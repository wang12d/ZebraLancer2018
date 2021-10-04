package ra

import (
	"crypto/ed25519"
	"log"
	"sync"

	"github.com/wang12d/ZebraLancer2018/pkg"
)

type Certificate []byte

type RegisterAuthority interface {
	CertGen(msg []byte) Certificate
}
type ra struct {
	msk []byte
	mpk []byte
}

var (
	RA     *ra
	RAInit sync.Once
)

func init() {
	RAInit.Do(
		func() {
			pk, sk, err := pkg.Setup()
			if err != nil {
				log.Fatalf("RA key generation error: %v\n", err)
			}
			RA.mpk, RA.msk = pk, sk
		})
}

// CertGen certificate the public key using master secret key by the CA
func (a *ra) CertGen(msg []byte) Certificate {
	return ed25519.Sign(a.msk, msg)
}

// Mpk return the public key of register authority
func (a *ra) Mpk() []byte {
	return a.mpk
}

// CertVrfy verifies the certification of public key generated by the CA
func CertVrfy(cert Certificate, pk, mpk []byte) bool {
	return ed25519.Verify(mpk, pk, cert)
}

// Pair verifies that the public key is corresponding to the private key
func Pair(pk ed25519.PublicKey, sk ed25519.PrivateKey) bool {
	return pk.Equal(sk.Public())
}
