package main

import (
	"fmt"
	"github.com/XFroggyX/chat-with-doubleratchet/encodeCharset"
	"github.com/farazdagi/x3dh"
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

type user struct {
	userConn net.Conn
	userName string
	userKey  x3dh.PublicKey
}

var channel [2]user

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

		idFreeChannel, statusChannel := channelFreeSpace()

		if !statusChannel {
			err := encodeCharset.WriteMsg(conn, "The channel is used")
			if err != nil {
				log.Panicln(err)
			}
		}

		channel[idFreeChannel].userConn = conn
		/*
			_, err = conn.Read(channel[idFreeChannel].userKey[:])
			if err != nil {
				return
			}

			_, err = conn.Write(channel[idFreeChannel+1%2].userKey[:])
			if err != nil {
				return
			}

			_, err = conn.Read(channel[idFreeChannel].userKey[:])
			if err != nil {
				return
			}
		*/
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

		log.Printf("Message Received: %s\n", msg)
		fmt.Println("Name: ", msg)
		broadcast(conn, msg)
	}

}

func removeConn(conn net.Conn) {
	var i int
	for i := range channel {
		if channel[i].userConn == conn {
			break
		}
	}
	channel[i].userConn = nil
	channel[i].userName = ""
	channel[i].userKey = [32]byte{}
}

func idConn(conn net.Conn) (i int) {
	for i := range channel {
		if channel[i].userConn == conn {
			return i
		}
	}
	return -1
}

func broadcast(conn net.Conn, msg string) {
	for i := range channel {
		if channel[i].userConn != conn && channel[i].userConn != nil {
			err := encodeCharset.WriteMsg(channel[i].userConn, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func channelFreeSpace() (int, bool) {
	for id, user := range channel {
		if user.userConn == nil {
			return id, true
		}
	}
	return -1, false
}
