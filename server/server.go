package main

import (
	"fmt"
	"github.com/XFroggyX/chat-with-doubleratchet/encodeCharset"
	"io"
	"log"
	"net"
	"os"
)

const (
	connectType = "tcp"
	connectPort = "3333"
	connectHost = "0.0.0.0"
)

var connections []net.Conn

func main() {
	l, err := net.Listen(connectType, connectHost+":"+connectPort)
	if err != nil {
		log.Panicln("Error listening: ", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	log.Println("Listening on " + connectHost + ":" + connectPort)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln("Error accepting: ", err.Error())
			os.Exit(1)
		}
		connections = append(connections, conn)
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	for {
		msg, err := encodeCharset.ReadMsg(conn)
		if err != nil {
			if err == io.EOF {
				removeConn(conn)
				conn.Close()
				return
			}
			log.Println(err)
			return
		}
		fmt.Printf("Message Received: %s\n", msg)
		broadcast(conn, msg)
	}
}

func removeConn(conn net.Conn) {
	var i int
	for i := range connections {
		if connections[i] == conn {
			break
		}
	}
	connections = append(connections[:i], connections[i+1:]...)
}

func broadcast(conn net.Conn, msg string) {
	for i := range connections {
		if connections[i] != conn {
			err := encodeCharset.WriteMsg(connections[i], msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
