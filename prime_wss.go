package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
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
		
TEST n1 [n2 ...]: Test the given number(s) n1 (, n2 ...) for Primality

CLOSE | EXIT: quits program
`
const eotMsg = "==== EOT ===="
const warnArgCount = "WARNING: Invalid number of args"
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
		case "TEST":
			if len(fields) > 1 {
				for _, numStr := range fields[1:] {
					var num uint64
					_, errFmt := fmt.Sscanf(numStr, "%d", &num)
					if errFmt == nil {
						result := fmt.Sprintf("%d: ", num)
						primeFactor := primelib.GetFirstPrimeFactor(num)
						if primeFactor == num {
							result += "PRIME"
						} else {
							result += fmt.Sprintf("divisible by %d", primeFactor)
						}
						err = websocket.Message.Send(ws, result)
					} else {
						err = websocket.Message.Send(ws, fmt.Sprintf("ERROR: %s - not a number", numStr))
					}
				}
			} else {
				err = websocket.Message.Send(ws, "ERROR: nothing to test")
			}
			err = websocket.Message.Send(ws, eotMsg)
		case "LIST":
			out := make(chan uint32, 10000)
			var cnt, from, to uint32
			switch len(fields) {
			case 3:
				_, errFmt1 := fmt.Sscanf(fields[1], "%d", &from)
				_, errFmt2 := fmt.Sscanf(fields[2], "%d", &to)
				if errFmt1 == nil && errFmt2 == nil {
					go primelib.ListPrimesBetween(out, from, to)
				} else {
					err = websocket.Message.Send(ws, "ERROR: parsing arg(s) as number(s)")
					err = websocket.Message.Send(ws, eotMsg)
				}
			case 2:
				_, errFmt := fmt.Sscanf(fields[1], "%d", &cnt)
				if errFmt == nil {
					go primelib.ListPrimes(out, cnt)
				} else {
					err = websocket.Message.Send(ws, "ERROR: parsing arg as number")
					err = websocket.Message.Send(ws, eotMsg)
				}
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

var portNum = flag.Uint("port", 7573, "(Websocket) port number to listen to")

func main() {
	flag.Parse()
	url := fmt.Sprintf(":%d", *portNum)
	fmt.Fprintln(os.Stderr, "DEBUG: portNum =", *portNum, "; ", url)

	fmt.Fprintln(os.Stderr, "Starting the Prime WebSocket Server listening at port", *portNum, "...")
	http.Handle("/", websocket.Handler(primeRequest))
	err := http.ListenAndServe(url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "FATAL ERROR:", err.Error())
	} else {
		fmt.Fprintln(os.Stderr, "Prime WebSocket Server shutdown normally")
	}
}
