package main

import (
	"fmt"
	"github.com/farazdagi/x3dh"
	"log"
	"net"
)

func main() {
	curve := x3dh.NewCurve25519()
	a_PrivateKey, err := curve.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}

	a_PublicKey := curve.PublicKey(a_PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A private: ", a_PrivateKey)
	fmt.Println("A public: ", a_PublicKey)

	curve1 := x3dh.NewCurve25519()
	b_PrivateKey, err := curve1.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}

	b_PublicKey := curve1.PublicKey(b_PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Println("B private: ", b_PrivateKey)
	fmt.Println("B public: ", b_PublicKey)

	sec_a := curve.ComputeSecret(a_PrivateKey, b_PublicKey)

	fmt.Println()
	fmt.Println("Sec A: ", sec_a)

	sec_b := curve1.ComputeSecret(b_PrivateKey, a_PublicKey)
	fmt.Println("Sec B: ", sec_b, 2%2)

	type user struct {
		userConn net.Conn
		userName string
		user     x3dh.PublicKey
	}
	fmt.Println(user{}.userConn)
}
