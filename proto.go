package main

import (
	"bytes"
	"fmt"

	"github.com/tidwall/resp"
)

const (
	CommandSet    = "set"
	CommandGet    = "get"
	CommandHello  = "hello"
	CommandClient = "client"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

type ClientCommand struct {
	value string
}

type HelloCommand struct {
	value string
}

type GetCommand struct {
	key []byte
}

func respWriteMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteString(k)
		rw.WriteString(":" + v)
		// buf.WriteString(fmt.Sprintf("+%s\r\n", k))
		// buf.WriteString(fmt.Sprintf(":%s\r\n", v))
	}
	return buf.Bytes()
}
