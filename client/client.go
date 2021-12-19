package main

import (
	"bufio"
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

func main() {
	tcpAddr, err := net.ResolveTCPAddr(connectType, connectHost+":"+connectPort)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go printOutput(conn)
	writeInput(conn)
}

func writeInput(conn *net.TCPConn) {
	fmt.Print("Enter username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	username = username[:len(username)-1]
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Enter text: ")
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		err = encodeCharset.WriteMsg(conn, username+": "+text)
		if err != nil {
			log.Println(err)
		}
	}
}

func printOutput(conn *net.TCPConn) {
	for {

		msg, err := encodeCharset.ReadMsg(conn)
		if err == io.EOF {
			conn.Close()
			fmt.Println("Connection Closed. Bye bye.")
			os.Exit(0)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(msg)
	}
}
