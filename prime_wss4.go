package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"os"
	"primelib/v4"
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

const beginHelpMsg = `
===============================
Primelib v4 Web(socket) Service
-------------------------------
Available Commands (case-sensitive):
`
const endHelpMsg = `
nodiag:
    turns off diagnostic messages

wantdiag:
	turns on diagnostic messages

exit | quit | close:
    closes the websocket client connection @server-end
===============================
`

const warnArgCount = "WARNING: Invalid number of args"
const errorUnknownCmd = "ERROR: Unknown request / command"

const MB = 1024 * 1024

func primeRequest(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1*MB, 1*MB)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println("ERROR: Websocket Handshake -", err)
		return
	}
	fmt.Fprintln(os.Stderr, "INFO: Websocket opened")
	primeRequestMainLoop(ws)
	ws.Close()
	fmt.Fprintln(os.Stderr, "INFO: Closing Websocket...")
}

func primeRequestMainLoop(ws *websocket.Conn) {
	wantDiagOutput := true
	msgType, msg, err := ws.ReadMessage()
MainLoop:
	for ; err == nil && msgType == websocket.TextMessage; msgType, msg, err = ws.ReadMessage() {
		req := strings.ToLower(strings.Trim(string(msg), " \r\n\t\v\f"))
		switch req {
		case "nodiag":
			wantDiagOutput = false
			ws.WriteMessage(websocket.TextMessage, []byte("Diag OFF\n"))
		case "wantdiag":
			wantDiagOutput = true
			ws.WriteMessage(websocket.TextMessage, []byte("Diag ON\n"))
		case "exit":
			fallthrough
		case "quit":
			fallthrough
		case "close":
			break MainLoop
		default:
			if w, e := ws.NextWriter(websocket.TextMessage); e == nil {
				if req == "?" {
					w.Write([]byte(beginHelpMsg))
				}
				_, out, diag := primelib.DoCmd(req)
				if out != nil {
					for output := range out {
						outputStr := fmt.Sprintf("%v", output)
						_, err = w.Write([]byte(outputStr))
						_, err = w.Write([]byte("\n"))
					}
				}
				if wantDiagOutput {
					for diagOutput := range diag {
						diagOutputStr := fmt.Sprintf("%v", diagOutput)
						_, err = w.Write([]byte(diagOutputStr))
						_, err = w.Write([]byte("\n"))
					}
				}
				if req == "?" {
					w.Write([]byte(endHelpMsg))
				}
				w.Close()
			} else {
				log.Printf("ERROR: %v\n", e.Error())
			}
		}
	}
}
