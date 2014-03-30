package main

import (
	"bufio"
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

	for msg := "?"; msg != "EXIT" && msg != "QUIT" && msg != "CLOSE"; msg = getLine() {
		//fmt.Println("DEBUG: msg =", msg)
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
