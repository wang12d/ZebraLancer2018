package pkg_test

import (
	"fmt"
	"testing"

	"github.com/wang12d/ZebraLancer2018/pkg"
)

func TestAuth(t *testing.T) {
	pk, sk, err := pkg.Setup()
	if err != nil {
		t.Log("Key generation error")
	}
	mpk, msk, err := pkg.Setup()
	if err != nil {
		t.Log("Key generation error")
	}
	fmt.Printf("pk: %x\nsk: %x\n", pk, sk)
	fmt.Printf("mpk: %x\nmsk: %x\n", mpk, msk)
	prefix, msg := []byte("hello"), []byte(" world.")
	cert := pkg.CertGen(msk, pk)
	fmt.Printf("prefix: %x\nmsg: %x\ncert: %x\n", prefix, msg, cert)
}
