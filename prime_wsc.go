package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"os"
	"strings"
)

func main() {
	wssUrl := "ws://localhost:7573"
	ws, err := websocket.Dial(wssUrl, "", "http://localhost")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: connecting to %s - %s\n", wssUrl, err)
		return
	}

	msg := "ListPrimes 100"
	err = websocket.Message.Send(ws, msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: could nto send WS request")
		return
	}

	err = websocket.Message.Receive(ws, &msg)
	msg = strings.Trim(msg, " \t\v\f\r\n")
	for err == nil && len(msg) > 0 && msg != "==== EOT ====" {
		fmt.Println(msg)
		err = websocket.Message.Receive(ws, &msg)
		msg = strings.Trim(msg, " \t\v\f\r\n")
	}
}
