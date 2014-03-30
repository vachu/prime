package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
	"os"
	"primelib/v3"
	"strings"
)

const helpMsg = `Available Commands:

LIST n: list out the first 'n' primes; n > 0 && < 2^32 - 1

LIST m n: list out all the primes between the numbers 'm' and 'n'
          where 'm' <= 'n' and ('m', 'n') > 0 && < 2^32 - 1

CLOSE | EXIT: quits program
`
const eotMsg = "==== EOT ===="
const warnArgCount = "WARNING: Invalid number of args\n"
const errorUnknownCmd = "ERROR: Unknown request / command"

func primeRequest(ws *websocket.Conn) {
	var req string
	var err error
	err = websocket.Message.Receive(ws, &req)
	for ; err == nil; err = websocket.Message.Receive(ws, &req) {
		fields := strings.Fields(req)
		fmt.Println("DEBUG:", fields)
		if len(fields) == 0 {
			fmt.Fprintln(os.Stderr, warnArgCount)
			err = websocket.Message.Send(ws, warnArgCount)
			err = websocket.Message.Send(ws, eotMsg)
			continue
		}

		method := strings.ToUpper(fields[0])
		switch method {
		case "?":
			err = websocket.Message.Send(ws, helpMsg)
			err = websocket.Message.Send(ws, eotMsg)
		case "LIST":
			out := make(chan uint32, 10000)
			var cnt, from, to uint32
			switch len(fields) {
			case 3:
				fmt.Sscanf(fields[1], "%d", &from)
				fmt.Sscanf(fields[2], "%d", &to)
				go primelib.ListPrimesBetween(out, from, to)
			case 2:
				fmt.Sscanf(fields[1], "%d", &cnt)
				go primelib.ListPrimes(out, cnt)
			default:
				fmt.Fprintln(os.Stderr, warnArgCount)
				err = websocket.Message.Send(ws, warnArgCount)
				err = websocket.Message.Send(ws, eotMsg)
				continue
			}
			for p := range out {
				err = websocket.Message.Send(ws, fmt.Sprintf("%d\n", p))
			}
			err = websocket.Message.Send(ws, eotMsg)
		case "EXIT":
			fallthrough
		case "QUIT":
			fallthrough
		case "CLOSE":
			fmt.Fprintln(os.Stderr, "INFO: Closing Websocket...")
			break
		default:
			fmt.Fprintln(os.Stderr, errorUnknownCmd, method)
			err = websocket.Message.Send(ws, errorUnknownCmd)
			err = websocket.Message.Send(ws, eotMsg)
		}
	}
	ws.Close()
}

func main() {
	http.Handle("/", websocket.Handler(primeRequest))
	fmt.Fprintln(os.Stderr, "Starting the Prime WebSocket Server...")
	err := http.ListenAndServe(":7573", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "FATAL ERROR:", err.Error())
	} else {
		fmt.Fprintln(os.Stderr, "Prime WebSocket Server shutdown normally")
	}
}
