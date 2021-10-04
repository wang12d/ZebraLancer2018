package pkg_test

import (
	"fmt"
	"testing"

	"github.com/wang12d/ZebraLancer2018/pkg"
)

func TestAuth(t *testing.T) {
	mpk, msk, err := pkg.Setup()
	if err != nil {
		t.Logf("Key generation error: %v", err)
	}
	pk, sk, err := pkg.KeyGeneration()
	if err != nil {
		t.Logf("User key generation error: %v", err)
	}
	fmt.Printf("pk: %x\nsk: %x\n", pk, sk)
	fmt.Printf("mpk: %x\nmsk: %x\n", mpk, msk)
	prefix, msg := []byte("hello"), []byte(" world.")
	cert := pkg.CertGen(msk, pk)
	fmt.Printf("prefix: %x\nmsg: %x\ncert: %x\n", prefix, msg, cert)
	if !pkg.CertVrfy(cert, pk, mpk) {
		t.Logf("Certificate verification error")
	}
}
