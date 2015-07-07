package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/yamux"
	"io"
	"net"
	"net/http"
)

var log = logrus.New()

func tunnel(sconn net.Conn) {
	buff := make([]byte, 0xff)

	// TOOD: read control message to figure out what port is being forwarded

	// TODO: bind to available local port to proxy
	local, _ := net.Listen("tcp", ":0")
	defer local.Close()
	log.Info(local.Addr())

	// TODO: send url back in ack message

	// TODO: start proxying

	for {
		_, err := sconn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Stream read error: %s", err)
			break
		}
		// log.Printf("stream sent %d bytes: %s", n, buff[:n])
		sconn.Write([]byte("GET /\n\n"))
	}
}

func handler(ws *websocket.Conn) {
	// Setup server side of yamux
	session, err := yamux.Server(ws, nil)
	if err != nil {
		panic(err)
	}

	// Handle new streams
	for {
		stream, err := session.Accept()
		if err != nil {
			if session.IsClosed() {
				// TODO: tunnel is no longer needed, close locally bound ports for this session
				log.Printf("TCP closed")
				break
			}
			// Print erros
			log.Printf("Yamux accept: %s", err)
			continue
		}
		go tunnel(stream)
	}
}

func main() {
	log.Println("Starting subway server on port 3000...")
	http.Handle("/", websocket.Handler(handler))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
