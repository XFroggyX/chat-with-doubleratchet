package main

import (
	"bufio"
	"fmt"
	"github.com/XFroggyX/chat-with-doubleratchet/encodeCharset"
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

var (
	sharedHka = [32]byte{
		0xbd, 0x29, 0x18, 0xcb, 0x18, 0x6c, 0x26, 0x32,
		0xd5, 0x82, 0x41, 0x2d, 0x11, 0xa4, 0x55, 0x87,
		0x1e, 0x5b, 0xa3, 0xb5, 0x5a, 0x6d, 0xe1, 0x97,
		0xde, 0xf7, 0x5e, 0xc3, 0xf2, 0xec, 0x1d, 0xd,
	}
	sharedNhkb = [32]byte{
		0x32, 0x89, 0x3a, 0xed, 0x4b, 0xf0, 0xbf, 0xc1,
		0xa5, 0xa9, 0x53, 0x73, 0x5b, 0xf9, 0x76, 0xce,
		0x70, 0x8e, 0xe1, 0xa, 0xed, 0x98, 0x1d, 0xe3,
		0xb4, 0xe9, 0xa9, 0x88, 0x54, 0x94, 0xaf, 0x23,
	}
)

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

	//key

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

	//doubleratchet
	keyPair, err := doubleratchet.DefaultCrypto{}.GenerateDH()
	if err != nil {
		log.Fatal(err)
	}

	var sharedKey [32]byte
	copy(sharedKey[:], myKeyPair.sharedKey)
	bob, err := doubleratchet.NewHE(sharedKey, sharedHka, sharedNhkb, keyPair)
	if err != nil {
		log.Fatal(err)
	}

	bobPublic := keyPair.PublicKey()
	_, err = conn.Write(bobPublic[:])
	if err != nil {
		return
	}

	ms := doubleratchet.MessageHE{}

	head, err := encodeCharset.ReadMsg(conn)

	text, err := encodeCharset.ReadMsg(conn)

	ms.Header = []byte(head)
	ms.Ciphertext = []byte(text)

	plaintext, err := bob.RatchetDecrypt(ms, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(plaintext))

	//chat
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		err = encodeCharset.WriteMsg(conn, text)
		if err != nil {
			log.Println(err)
		}

		msg, err := encodeCharset.ReadMsg(conn)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				return
			}
			log.Println(err)
			return
		}

		fmt.Println("Bob: ", msg)
	}
}
