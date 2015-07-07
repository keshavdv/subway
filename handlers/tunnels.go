package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/context"
	"github.com/hashicorp/yamux"
	"github.com/unrolled/render"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

func GetTunnels(w http.ResponseWriter, req *http.Request) {
	// TODO
}

func CreateTunnel(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)

	port, err := strconv.Atoi(req.URL.Query().Get("port"))
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "failed", "msg": "Port is invalid."})
		return
	}

	session := context.Get(req, "Mux").(*yamux.Session)
	url, err := handleCreateTunnel(session, port)
	if err != nil {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed", "msg": err.Error()})
	} else {
		r.JSON(w, http.StatusOK, map[string]string{"status": "started", "url": url})
	}
}

func DeleteTunnel(w http.ResponseWriter, req *http.Request) {
	// TODO
}

func handleCreateTunnel(session *yamux.Session, port int) (url string, err error) {
	// TODO: Attempt to bind to local port
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	local, err := net.Dial("tcp", addr)
	if err != nil {
		return "", errors.New("Could not establish connection to host port.")
	}

	// Establish stream
	stream, err := session.Open()
	if err != nil {
		return "", errors.New("Could not create stream.")
	}

	// TODO: send control message

	// TODO: wait for ack w/ url
	url = "tcp://somehost:8000"

	// TODO: Update global state for later control messages

	// TODO: launch goroutine to proxy data
	go proxy(stream, local)

	return url, nil
}

func proxy(remote net.Conn, local net.Conn) {
	defer remote.Close()
	defer local.Close()

	log.Println("test")
	_, err := remote.Write([]byte("hello "))
	if err != nil {
		return
	}
	// TODO: proxy data in both directions
	var waitFor chan struct{}

	go io.Copy(remote, local)
	go io.Copy(local, remote)
	<-waitFor
	log.Println("Done.")
}
