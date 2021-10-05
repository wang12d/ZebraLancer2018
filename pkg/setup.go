package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
)

// Setup generate the public parameter of ZSK and public key and private key
// of digital signature scheme
func Setup() (*rsa.PublicKey, *rsa.PrivateKey, error) {
	msk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Generate digital signature key error: %v\n", err)
	}
	return &msk.PublicKey, msk, err
}
