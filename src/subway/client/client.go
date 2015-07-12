package client

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/hashicorp/yamux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/unrolled/render"
	"gopkg.in/alecthomas/kingpin.v2"
	"subway/context"
	"subway/handlers"
)

var log = logrus.New()

var (
	verbose = kingpin.Flag("verbose", "show debug information").Short('v').Bool()
	port    = kingpin.Flag("port", "port to listen on for REST api").Short('p').Default("3001").Int()
	host    = kingpin.Flag("host", "subway server to connect to (in the form host:port)").Short('h').Default("localhost:3000").String()
)

func Main() {
	kingpin.Parse()

	// Create websocket session for tunnels
	origin := fmt.Sprintf("http://%s/", *host)
	url := fmt.Sprintf("ws://%s/", *host)
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not connect to subway server (%s)!", url))
	}

	session, err := yamux.Client(conn, nil)
	if err != nil {
		panic(err)
	}

	// Start REST api
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.CreateTunnel)

	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	render := render.New(render.Options{})
	subway := context.CreateSubway(session, render)
	n.Use(subway)
	n.UseHandler(router)
	n.Run(fmt.Sprintf(":%v", *port))

}
