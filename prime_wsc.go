package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"os"
	"strings"
)

var wsUrl = flag.String("wsurl", "ws://localhost:7573", "the websocket Url")

func main() {
	flag.Parse()
	fmt.Fprintln(os.Stderr, "INFO: Connecting to WebSocket Server @", *wsUrl, "...")
	ws, err := websocket.Dial(*wsUrl, "", "http://localhost")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: connecting to %s - %s\n", wsUrl, err)
		return
	}
	fmt.Fprintln(os.Stderr, "INFO: Connected\n")

	for msg := "?"; msg != "EXIT" && msg != "QUIT" && msg != "CLOSE"; msg = getLine() {
		err = websocket.Message.Send(ws, msg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: could not send Websocket request - ", err)
			continue
		}

		err = websocket.Message.Receive(ws, &msg)
		msg = strings.Trim(msg, " \t\v\f\r\n")
		for err == nil && msg != "==== EOT ====" {
			fmt.Println(msg)

			err = websocket.Message.Receive(ws, &msg)
			msg = strings.Trim(msg, " \t\v\f\r\n")
		}
		fmt.Print("Enter Request / Command: ")
	}
	websocket.Message.Send(ws, "CLOSE")
	ws.Close()
}

func getLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.Trim(line, " \t\v\f\r\n")
	line = strings.ToUpper(line)
	return line
}
