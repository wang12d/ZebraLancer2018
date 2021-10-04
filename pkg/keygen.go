package pkg

import (
	"crypto/ed25519"
	"crypto/rand"
)

// KeyGeneration generates a pair of public key and private key for
// identity identification
func KeyGeneration() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return pk, sk, err
}
