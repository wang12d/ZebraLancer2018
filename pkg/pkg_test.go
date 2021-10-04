package pkg_test

import (
	"fmt"
	"testing"

	"github.com/wang12d/ZebraLancer2018/pkg"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
)

func TestAuth(t *testing.T) {
	pk, sk, err := pkg.KeyGeneration()
	if err != nil {
		t.Logf("User key generation error: %v", err)
	}
	fmt.Printf("pk: %x\nsk: %x\n", pk, sk)
	prefix, msg := []byte("hello"), []byte(" world.")
	cert := ra.RA.CertGen(pk)
	fmt.Printf("prefix: %x\nmsg: %x\ncert: %x\n", prefix, msg, cert)
	if !ra.CertVrfy(cert, pk, ra.RA.Mpk()) {
		t.Logf("Certificate verification error")
	}
}
