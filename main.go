package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/gorilla/websocket"
	"github.com/matti/betterio"
)

func main() {
	listenAddress := os.Args[1]
	upstreamURL := os.Args[2]
	fmt.Println("listen", listenAddress)
	ln, err := net.Listen("tcp", listenAddress)

	if err != nil {
		panic(err)
	}

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			log.Printf("accept err: %v", e)
			return
		}

		go func() {
			log.Println(conn.RemoteAddr().String(), "handling")
			handle(conn, upstreamURL)
			log.Println(conn.RemoteAddr().String(), "handled")
		}()
	}
}

// from https://github.com/gorilla/websocket/issues/282
type rwc struct {
	r io.Reader
	c *websocket.Conn
}

func (c *rwc) Write(p []byte) (int, error) {
	err := c.c.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *rwc) Read(p []byte) (int, error) {
	for {
		if c.r == nil {
			// Advance to next message.
			var err error
			_, c.r, err = c.c.NextReader()
			if err != nil {
				return 0, err
			}
		}
		n, err := c.r.Read(p)
		if err == io.EOF {
			// At end of message.
			c.r = nil
			if n > 0 {
				return n, nil
			} else {
				// No data read, continue to next message.
				continue
			}
		}
		return n, err
	}
}

func (c *rwc) Close() error {
	return c.c.Close()
}

func handle(conn net.Conn, upstreamURL string) {
	defer conn.Close()

	ws, _, err := websocket.DefaultDialer.Dial(upstreamURL, nil)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	upstream := &rwc{c: ws}

	log.Println(conn.RemoteAddr().String(), "copy", ws.RemoteAddr().String())

	betterio.CopyBidirUntilCloseAndReturnBytesWritten(conn, upstream)
}
