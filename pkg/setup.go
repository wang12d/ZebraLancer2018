package pkg

import (
	"crypto/ed25519"
	"crypto/rand"
)

// Setup generate the public parameter of ZSK and public key and private key
// of digital signature scheme
func Setup() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(rand.Reader)
}