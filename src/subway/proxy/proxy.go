package proxy

import (
	"io"
	"log"
	"net"
)

func Proxy(left, right net.Conn) {
	// channels to wait on the close event for each connection
	leftClosed := make(chan struct{}, 1)
	rightClosed := make(chan struct{}, 1)

	go broker(left, right, leftClosed)
	go broker(left, right, rightClosed)

	// wait for one half of the proxy to exit, then trigger a shutdown of the
	// other half by calling CloseRead(). This will break the read loop in the
	// broker and allow us to fully close the connection cleanly without a
	// "use of closed network connection" error.
	var waitUntilDisconnect chan struct{}
	select {
	case <-leftClosed:
		// faster.
		right.Close()
		waitUntilDisconnect = rightClosed
	case <-rightClosed:
		left.Close()
		waitUntilDisconnect = leftClosed
	}

	<-waitUntilDisconnect
}

func broker(dst, src net.Conn, srcClosed chan struct{}) {
	_, err := io.Copy(dst, src)

	if err != nil {
		log.Printf("Copy error: %s", err)
	}
	if err := src.Close(); err != nil {
		log.Printf("Close error: %s", err)
	}
	srcClosed <- struct{}{}
}
