package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/context"
	"github.com/hashicorp/yamux"
	"github.com/keshavdv/subway/msg"
	"github.com/unrolled/render"
	"io"
	"github.com/Sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
)

var log = logrus.New()

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
	log.WithFields(logrus.Fields{
		"server": addr,
	}).Info("attempting to connect to server")
	local, err := net.Dial("tcp", addr)
	if err != nil {
		return "", errors.New("Could not establish connection to host port.")
	}
	log.WithFields(logrus.Fields{
		"server": addr,
	}).Info("connected to server")

	// Establish stream
	stream, err := session.Open()
	if err != nil {
		return "", errors.New("Could not create stream.")
	}

	// TODO: send control message
	log.WithFields(logrus.Fields{
		"server": addr,
	}).Info("sending TunnelRequest")
	req := &msg.TunnelRequest{
		Port:  port,
	}
	if err := msg.WriteMsg(stream, req); err != nil {
		return "", err
	}

	// TODO: wait for ack w/ url
	log.WithFields(logrus.Fields{
		"server": addr,
	}).Info("waiting for TunnelReply")

	var res msg.TunnelReply
	if err := msg.ReadMsgInto(stream, &res); err != nil {
		return "", err
	}
	log.WithFields(logrus.Fields{
		"server": addr,
		"port": port,
		"uri": res.URI,
	}).Info("received TunnelReply")
	url = res.URI

	// TODO: Update global state for later control messages

	
	// TODO: launch goroutine to proxy data
	go proxy(stream, local)

	return url, nil
}

func proxy(remote net.Conn, local net.Conn) {
	defer remote.Close()
	defer local.Close()

	// TODO: proxy data in both directions
	var waitFor chan struct{}

	go io.Copy(remote, local)
	go io.Copy(local, remote)
	<-waitFor
	log.Println("Done.")
}
