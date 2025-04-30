package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				var cmd Command
				fmt.Printf("this is the value we are getting %s", value.String())
				switch value.String() {
				case CommandSet:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of arguments in set command")
					}
					cmd = SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}

				case CommandGet:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of arguments in GET command")
					}
					cmd = GetCommand{
						key: v.Array()[1].Bytes(),
					}
				case CommandHello:
					cmd = HelloCommand{
						value: v.Array()[1].String(),
					}
				default:
					fmt.Printf("got unknoen command => %+v\n", v.Array())
				}
				p.msgCh <- Message{
					cmd:  cmd,
					peer: p,
				}
			}
		}
		// return fmt.Errorf("invalid or unknown command recived: ", raw)
	}
	return nil
}
