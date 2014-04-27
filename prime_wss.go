package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
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
		log.Fatalln("FATAL: Illegal port number -", *portNum)
	}

	log.Println("INFO: Starting the Prime WebSocket Server listening at port", *portNum, "...")
	http.HandleFunc("/", primeRequest)
	if err := http.ListenAndServe(url, nil); err != nil {
		log.Fatalln("FATAL:", err.Error())
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
const MB = 1024 * 1024

func primeRequest(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1*MB, 1*MB)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
	} else if err != nil {
		log.Println("ERROR: Websocket Handshake -", err)
		return
	}
	primeRequestMainLoop(ws)
	ws.Close()
}

func primeRequestMainLoop(ws *websocket.Conn) {
	msgType, msg, err := ws.ReadMessage()
	for ; err == nil && msgType == websocket.TextMessage; msgType, msg, err = ws.ReadMessage() {
		req := string(msg)
		fields := strings.Fields(req)
		if len(fields) == 0 {
			err = ws.WriteMessage(websocket.TextMessage, []byte(warnArgCount))
		} else {
			method := strings.ToUpper(fields[0])
			switch method {
			case "?":
				err = ws.WriteMessage(msgType, []byte(helpMsg))
			case "TEST":
				if len(fields) > 1 {
					err = OnTest(ws, fields[1:])
				} else {
					err = ws.WriteMessage(msgType, []byte("ERROR: nothing to test"))
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
				err = ws.WriteMessage(msgType, []byte(errorUnknownCmd))
			}
		}
		err = ws.WriteMessage(websocket.TextMessage, []byte(eotMsg))
	}
}

func OnTest(ws *websocket.Conn, fields []string) (err error) {
	bws := &BufferedWebSocket{bufSize: 1024 * 1024, ws: ws} // 1MB buffer
	for _, numStr := range fields {
		var num uint64
		var result string
		_, errFmt := fmt.Sscanf(numStr, "%d", &num)
		if errFmt == nil {
			result = fmt.Sprintf("%d: ", num)
			primeFactor := primelib.GetFirstPrimeFactor(num)
			if primeFactor == num {
				result += "PRIME\n"
			} else {
				result += fmt.Sprintf("divisible by %d\n", primeFactor)
			}
		} else {
			result = fmt.Sprintf("ERROR: %s - not a number\n", numStr)
		}

		if err = bws.Send(result); err != nil {
			break
		}
	}
	bws.Flush()
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
			err = ws.WriteMessage(websocket.TextMessage, []byte("ERROR: parsing arg(s) as number(s)"))
			return
		}
		go primelib.ListPrimesBetween(out, from, to)
	case 1:
		_, errFmt := fmt.Sscanf(fields[0], "%d", &cnt)
		if errFmt != nil {
			err = ws.WriteMessage(websocket.TextMessage, []byte("ERROR: parsing arg as number"))
			return
		}
		go primelib.ListPrimes(out, cnt)
	default:
		err = ws.WriteMessage(websocket.TextMessage, []byte(warnArgCount))
		return
	}

	bws := &BufferedWebSocket{bufSize: 1024 * 1024, ws: ws} // 1MB buffer
	for p := range out {
		if err == nil {
			err = bws.Send(fmt.Sprintf("%d\n", p))
		}
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

	if err == nil {
		bws.buf = append(bws.buf, msgBytes...)
	}
	return err
}

func (bws *BufferedWebSocket) Flush() error {
	err := bws.ws.WriteMessage(websocket.TextMessage, bws.buf)
	bws.buf = make([]byte, 0, bws.bufSize)
	return err
}
