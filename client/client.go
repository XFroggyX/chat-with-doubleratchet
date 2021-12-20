package main

import (
	"bufio"
	"fmt"
	"github.com/farazdagi/x3dh"
	"github.com/tiabc/doubleratchet"
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

type keysPair struct {
	curve      x3dh.Curve25519
	publicKey  x3dh.PublicKey
	privateKey x3dh.PrivateKey
	sharedKey  []byte
}

func (k *keysPair) generatingKeyPair() {
	k.curve = x3dh.NewCurve25519()
	k.privateKey, _ = k.curve.GenerateKey(nil)
	k.publicKey = k.curve.PublicKey(k.privateKey)
}

var sk = [32]byte{
	0xeb, 0x8, 0x10, 0x7c, 0x33, 0x54, 0x0, 0x20,
	0xe9, 0x4f, 0x6c, 0x84, 0xe4, 0x39, 0x50, 0x5a,
	0x2f, 0x60, 0xbe, 0x81, 0xa, 0x78, 0x8b, 0xeb,
	0x1e, 0x2c, 0x9, 0x8d, 0x4b, 0x4d, 0xc1, 0x40,
}

func main() {
	myKeyPair := &keysPair{}
	myKeyPair.generatingKeyPair()

	tcpAddr, err := net.ResolveTCPAddr(connectType, connectHost+":"+connectPort)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	/*
		//write
		_, err = conn.Write(myKeyPair.publicKey[:])
		if err != nil {
			return
		}

		var userPublicKey x3dh.PublicKey
		_, err = conn.Read(userPublicKey[:])
		if err != nil {
			return
		}

		myKeyPair.sharedKey = myKeyPair.curve.ComputeSecret(myKeyPair.privateKey, userPublicKey)
		fmt.Println("Key: ", myKeyPair.sharedKey)

		fmt.Println("Key: ", myKeyPair.publicKey)*/

	keyPair, err := doubleratchet.DefaultCrypto{}.GenerateDH()
	if err != nil {
		log.Fatal(err)
	}

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
