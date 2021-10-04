package pkg

import (
	"crypto/ed25519"
	"crypto/rand"
	"log"
)

// Setup generate the public parameter of ZSK and public key and private key
// of digital signature scheme
func Setup() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	mpk, msk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Generate digital signature key error: %v\n", err)
	}
	return mpk, msk, err
}
