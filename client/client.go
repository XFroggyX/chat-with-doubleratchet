package main

import (
	"bufio"
	"fmt"
	"github.com/XFroggyX/chat-with-doubleratchet/encodeCharset"
	"github.com/farazdagi/x3dh"
	"github.com/tiabc/doubleratchet"
	"log"
	"net"
	"os"
)

const (
	connectType = "tcp"
	connectPort = "3333"
	connectHost = "0.0.0.0"
)

type KeysPair struct {
	curve      x3dh.Curve25519
	publicKey  x3dh.PublicKey
	privateKey x3dh.PrivateKey
	sharedKey  []byte
}

func (k *KeysPair) GeneratingKeyPair() {
	k.curve = x3dh.NewCurve25519()
	k.privateKey, _ = k.curve.GenerateKey(nil)
	k.publicKey = k.curve.PublicKey(k.privateKey)
}

func generatingKeys(listKeys []*KeysPair) {
	for i := 0; i < 3; i++ {
		listKeys[i] = &KeysPair{}
		listKeys[i].GeneratingKeyPair()
	}

}

func main() {
	tcpAddr, err := net.ResolveTCPAddr(connectType, connectHost+":"+connectPort)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			log.Panicln("Error close connect: ", err.Error())
		}
	}(conn)

	listKeys := make([]*KeysPair, 3)
	generatingKeys(listKeys)

	for _, myKeyPair := range listKeys {
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
	}

	keyPair, err := doubleratchet.DefaultCrypto{}.GenerateDH()
	if err != nil {
		log.Panic(err)
	}

	var sharedKey [32]byte
	copy(sharedKey[:], listKeys[0].sharedKey)
	var sharedHka [32]byte
	copy(sharedHka[:], listKeys[1].sharedKey)
	var sharedNhkb [32]byte
	copy(sharedNhkb[:], listKeys[2].sharedKey)

	var bobPublic [32]byte
	_, err = conn.Read(bobPublic[:])
	if err != nil {
		return
	}

	alicePublic := keyPair.PublicKey()
	_, err = conn.Write(alicePublic[:])
	if err != nil {
		return
	}

	countMsg := 0
	for countMsg < 2000 {
		alice, err := doubleratchet.NewHEWithRemoteKey(sharedKey, sharedHka, sharedNhkb, bobPublic)
		if err != nil {
			log.Panic(err)
		}

		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		m := alice.RatchetEncrypt([]byte(msg), nil)

		err = encodeCharset.WriteMsg(conn, string(m.Header))
		if err != nil {
			log.Println(err)
		}

		err = encodeCharset.WriteMsg(conn, string(m.Ciphertext))
		if err != nil {
			log.Println(err)
		}

		// send

		aliceR, err := doubleratchet.NewHE(sharedKey, sharedHka, sharedNhkb, keyPair)
		if err != nil {
			log.Fatal(err)
		}

		ms := doubleratchet.MessageHE{}

		head, err := encodeCharset.ReadMsg(conn)

		text, err := encodeCharset.ReadMsg(conn)

		ms.Header = []byte(head)
		ms.Ciphertext = []byte(text)

		plaintext, err := aliceR.RatchetDecrypt(ms, nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(plaintext))

		countMsg++
	}

}
