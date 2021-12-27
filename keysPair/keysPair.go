package keyspair

import "github.com/farazdagi/x3dh"

type KeysPair struct {
	curve      x3dh.Curve25519
	publicKey  x3dh.PublicKey
	privateKey x3dh.PrivateKey
	sharedKey  []byte
}

func (k *KeysPair) generatingKeyPair() {
	k.curve = x3dh.NewCurve25519()
	k.privateKey, _ = k.curve.GenerateKey(nil)
	k.publicKey = k.curve.PublicKey(k.privateKey)
}
