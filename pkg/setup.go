package pkg

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"log"
)

// Setup generate the public parameter of ZSK and public key and private key
// of digital signature scheme
func Setup() (ed25519.PublicKey, ed25519.PrivateKey, *rsa.PrivateKey) {
	sk, pk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Generate digital signature key error: %v\n", err)
	}
	esk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Generate encryption key error: %v\n", err)
	}
	return sk, pk, esk
}
