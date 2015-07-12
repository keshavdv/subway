package ma2in

import (
	"github.com/keshavdv/subway/proxy"
	"log"
	"net"
)

func main() {
	a, _ := net.Listen("tcp", ":0")
	defer a.Close()
	log.Println(a.Addr())

	b, _ := net.Listen("tcp", ":0")
	defer b.Close()
	log.Println(b.Addr())

	conn1, err := a.Accept()
	if err != nil {
		// handle error
		log.Println("err")
	}

	conn2, err := b.Accept()
	if err != nil {
		// handle error
		log.Println("err")
	}

	log.Println("Starting proxy")
	proxy.Proxy(conn1, conn2)

}
