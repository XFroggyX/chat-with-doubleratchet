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
	listen, err := net.Listen(connectType, connectHost+":"+connectPort)
	if err != nil {
		log.Panicln("Error listening: ", err.Error())
		os.Exit(1)
	}

	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Panicln("Error close listening: ", err.Error())
		}
	}(listen)
	log.Println("Listening on " + connectHost + ":" + connectPort)

	conn, err := listen.Accept()
	if err != nil {
		log.Panicln("Error accepting: ", err.Error())
		os.Exit(1)
	}

	listKeys := make([]*KeysPair, 3)
	generatingKeys(listKeys)

	for _, myKeyPair := range listKeys {
		var userPublicKey x3dh.PublicKey
		_, err = conn.Read(userPublicKey[:])
		if err != nil {
			return
		}

		_, err = conn.Write(myKeyPair.publicKey[:])
		if err != nil {
			return
		}

		myKeyPair.sharedKey = myKeyPair.curve.ComputeSecret(myKeyPair.privateKey, userPublicKey)
		fmt.Println("Key: ", myKeyPair.sharedKey)
	}

	keyPair, err := doubleratchet.DefaultCrypto{}.GenerateDH()
	if err != nil {
		log.Fatal(err)
	}

	var sharedKey [32]byte
	copy(sharedKey[:], listKeys[0].sharedKey)
	var sharedHka [32]byte
	copy(sharedHka[:], listKeys[1].sharedKey)
	var sharedNhkb [32]byte
	copy(sharedNhkb[:], listKeys[2].sharedKey)

	bobPublic := keyPair.PublicKey()
	_, err = conn.Write(bobPublic[:])
	if err != nil {
		return
	}

	var alicePublic [32]byte
	_, err = conn.Read(alicePublic[:])
	if err != nil {
		return
	}

	countMsg := 0
	for countMsg < 2000 {
		bob, err := doubleratchet.NewHE(sharedKey, sharedHka, sharedNhkb, keyPair)
		if err != nil {
			log.Fatal(err)
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

		fmt.Print(string(plaintext))

		// send

		bobR, err := doubleratchet.NewHEWithRemoteKey(sharedKey, sharedHka, sharedNhkb, alicePublic)
		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		m := bobR.RatchetEncrypt([]byte(msg), nil)

		err = encodeCharset.WriteMsg(conn, string(m.Header))
		if err != nil {
			log.Println(err)
		}

		err = encodeCharset.WriteMsg(conn, string(m.Ciphertext))
		if err != nil {
			log.Println(err)
		}

		countMsg++
	}
	fmt.Println()
}
