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
		var cmd Command
		if v.Type() == resp.Array {
			rawCmd := v.Array()[0]
			switch rawCmd.String() {
			case CommandGet:
				if len(v.Array()) != 2 {
					return fmt.Errorf("invalid number of arguments in GET command")
				}
				cmd = GetCommand{
					key: v.Array()[1].Bytes(),
				}
			case CommandSet:
				if len(v.Array()) != 3 {
					return fmt.Errorf("invalid number of arguments in set command")
				}
				cmd = SetCommand{
					key: v.Array()[1].Bytes(),
					val: v.Array()[2].Bytes(),
				}
			case CommandHello:
				cmd = HelloCommand{
					value: v.Array()[1].String(),
				}
			case CommandClient:
				cmd = ClientCommand{
					value: v.Array()[1].String(),
				}
			default:
				fmt.Println("got this unhandled commmand", rawCmd)
			}
			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,
			}

			fmt.Println("this should be the cmd", v.Array()[0])
		}
	}
	return nil
}
