package gork

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh/terminal"
)

type ZIODev interface {
	Print(...interface{})
	ReadLine() string
}

type ZTerminal struct{}

func (_ ZTerminal) Print(s ...interface{}) {
	for _, si := range s {
		fmt.Print(si)
	}
}

func (_ ZTerminal) ReadLine() string {
	r := bufio.NewReader(os.Stdin)

	s, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}

	return s
}

type ZSshTerminal struct {
	Term *terminal.Terminal
}

func (sshTerm ZSshTerminal) Print(s ...interface{}) {
	for _, si := range s {
		sis := fmt.Sprint(si)
		sis = strings.Replace(sis, "\n", "\r\n", -1)
		sshTerm.Term.Write([]byte(sis))
	}
}

func (sshTerm ZSshTerminal) ReadLine() string {
	l, err := sshTerm.Term.ReadLine()

	if err != nil {
		panic(err)
	}

	return l
}

type ZWSDev struct {
	Conn *websocket.Conn
}

func (ws *ZWSDev) Print(s ...interface{}) {
	for _, si := range s {
		sis := fmt.Sprint(si)
		ws.Conn.WriteMessage(websocket.TextMessage, []byte(sis))
	}
}

func (ws *ZWSDev) ReadLine() string {
	msg_type, l, err := ws.Conn.ReadMessage()

	if msg_type != websocket.TextMessage || err != nil {
		panic(err)
	}

	return string(l)
}
