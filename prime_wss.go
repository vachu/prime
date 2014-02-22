package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
	"os"
	"primelib"
	"time"
)

func primeRequest(ws *websocket.Conn) {
	var req string
	err := websocket.Message.Receive(ws, &req)
	if err != nil {
		fmt.Fprintln(os.Stderr, time.Now, ": ERROR -", err)
	} else {
		var method string
		var cnt uint32
		_, err2 := fmt.Sscanf(req, "%s%d", &method, &cnt)
		if err2 != nil {
			fmt.Println(os.Stderr, time.Now, ": ERROR - invalid request;", req)
		} else {
			out := make(chan uint32, 1000)
			go primelib.ListPrimes(out, cnt)
			for p := range out {
				websocket.Message.Send(ws, fmt.Sprintf("%d", p))
			}
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(primeRequest))
	fmt.Fprintln(os.Stderr, "Starting the Prime WebSocket Server...")
	err := http.ListenAndServe(":7573", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL ERROR:", err.Error())
	} else {
		fmt.Fprintln(os.Stderr, "Prime WebSocket Server shutdown normally")
	}
}
