package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSet = "set"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSet:
					fmt.Printf("%+v\n", v.Array())
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("invalid number of arguments in set command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					return cmd, nil

				}
			}
		}
		return nil, fmt.Errorf("invalid or unknown command recived: ", raw)
	}

	return nil, fmt.Errorf("invalid or unknown command recived: ", raw)
}
