package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/yamux"
	"net"
	"net/http"
	"os"
	"subway/msg"
	"subway/proxy"
)

var log = logrus.New()

// var sessions = map[int]net.Conn

func init() {
	// sessions = make(map[string]int)
}

func tunnel(sconn net.Conn) {
	// buff := make([]byte, 0xff)

	// read control message to figure out what port is being forwarded
	var req msg.TunnelRequest
	if err := msg.ReadMsgInto(sconn, &req); err != nil {
		log.Error(err.Error())
	}

	log.WithFields(logrus.Fields{
		"client":     sconn.RemoteAddr(),
		"proxy_port": req.Port,
	}).Info("received TunnelRequest")

	// bind to available local port to proxy
	local, _ := net.Listen("tcp", ":0")
	defer local.Close()
	log.Info(local.Addr())

	// send url back in ack message
	hostname, _ := os.Hostname()
	port := local.Addr().(*net.TCPAddr).Port
	res := &msg.TunnelReply{
		URI: fmt.Sprintf("tcp://%s:%d", hostname, port),
	}
	if err := msg.WriteMsg(sconn, res); err != nil {
		log.Error(err.Error())
		return
	}

	// TODO: start proxying
	for {
		conn, err := local.Accept()
		if err != nil {
			// handle error
			log.Error("err")
		}
		log.Info("new conn")
		go proxy.Proxy(conn, sconn)
		// _, err := sconn.Read(buff)
		// if err != nil {
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	log.Printf("Stream read error: %s", err)
		// 	break
		// }
		// // log.Printf("stream sent %d bytes: %s", n, buff[:n])
		// sconn.Write([]byte("GET /\n\n"))
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
				log.Info("session closed")
				break
			}
			// Print errors
			log.Error("yamux error: %s", err)
			continue
		}
		go tunnel(stream)
	}
}

func Main() {
	log.Println("Starting subway server on port 3000...")
	http.Handle("/", websocket.Handler(handler))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
