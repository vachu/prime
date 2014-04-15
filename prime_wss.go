package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"primelib/v3"
	"strings"
)

var portNum = flag.Uint("port", 7573, "(Websocket) port number to listen to")

func main() {
	flag.Parse()
	url := fmt.Sprintf(":%d", *portNum)
	if *portNum > math.MaxUint16 {
		fmt.Fprintln(os.Stderr, "ERROR: Illegal port number -", *portNum)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "INFO: Starting the Prime WebSocket Server listening at port", *portNum, "...")
	http.Handle("/", websocket.Handler(primeRequest))
	if err := http.ListenAndServe(url, nil); err != nil {
		fmt.Fprintln(os.Stderr, "FATAL ERROR:", err.Error())
	}
}

const helpMsg = `Available Commands (case-insensitive):

LIST n: list out the first 'n' primes; n > 0 && < 2^32 - 1

LIST m n: list out all the primes between the numbers 'm' and 'n'
          where 'm' <= 'n' and ('m', 'n') > 0 && < 2^32 - 1
		
TEST n1 [n2 ...]: Test the given number(s) n1 (, n2 ...) for Primality

CLOSE | EXIT | QUIT: quits program
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
		if len(fields) == 0 {
			err = websocket.Message.Send(ws, warnArgCount)
		} else {
			method := strings.ToUpper(fields[0])
			switch method {
			case "?":
				err = websocket.Message.Send(ws, helpMsg)
			case "TEST":
				if len(fields) > 1 {
					err = OnTest(ws, fields[1:])
				} else {
					err = websocket.Message.Send(ws, "ERROR: nothing to test")
				}
			case "LIST":
				err = OnList(ws, fields[1:])
			case "EXIT":
				fallthrough
			case "QUIT":
				fallthrough
			case "CLOSE":
				fmt.Fprintln(os.Stderr, "INFO: Closing Websocket...")
				break
			default:
				err = websocket.Message.Send(ws, errorUnknownCmd)
			}
		}
		err = websocket.Message.Send(ws, eotMsg)
	}
	ws.Close()
}

func OnTest(ws *websocket.Conn, fields []string) (err error) {
	for _, numStr := range fields {
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

		if err != nil {
			break
		}
	}
	return
}

func OnList(ws *websocket.Conn, fields []string) (err error) {
	out := make(chan uint32, 10000)
	var cnt, from, to uint32
	switch len(fields) {
	case 2:
		_, errFmt1 := fmt.Sscanf(fields[0], "%d", &from)
		_, errFmt2 := fmt.Sscanf(fields[1], "%d", &to)
		if errFmt1 != nil || errFmt2 != nil {
			err = websocket.Message.Send(ws, "ERROR: parsing arg(s) as number(s)")
			return
		}
		go primelib.ListPrimesBetween(out, from, to)
	case 1:
		_, errFmt := fmt.Sscanf(fields[0], "%d", &cnt)
		if errFmt != nil {
			err = websocket.Message.Send(ws, "ERROR: parsing arg as number")
			return
		}
		go primelib.ListPrimes(out, cnt)
	default:
		err = websocket.Message.Send(ws, warnArgCount)
		return
	}

	bws := &BufferedWebSocket{bufSize: 1024 * 1024, ws: ws} // 1MB buffer
	for p := range out {
		err = bws.Send(fmt.Sprintf("%d\n", p))
	}
	err = bws.Flush()
	return
}

////////////////////////////////////////////////////////////////////////////////
type BufferedWebSocket struct {
	bufSize uint32
	ws      *websocket.Conn
	buf     []byte
}

func (bws *BufferedWebSocket) Send(msg string) error {
	if bws.buf == nil {
		bws.buf = make([]byte, 0, bws.bufSize)
	}

	var err error
	msgBytes := []byte(msg)
	if uint32(len(bws.buf)+len(msgBytes)) > bws.bufSize {
		err = bws.Flush()
	}

	bws.buf = append(bws.buf, msgBytes...)
	return err
}

func (bws *BufferedWebSocket) Flush() error {
	err := websocket.Message.Send(bws.ws, string(bws.buf))
	bws.buf = make([]byte, 0, bws.bufSize)
	return err
}
