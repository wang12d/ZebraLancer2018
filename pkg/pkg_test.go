package pkg_test

import (
	"fmt"
	"log"
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
	cert, err := ra.RA.CertGen(pk)
	fmt.Printf("%x\n", cert)
	if err != nil {
		log.Fatalf("Worker obtain signature error: %v\n", err)
	}
	fmt.Printf("%s\n", ra.RA.Mpk())
	if !ra.CertVrfy(cert, pk, ra.RA.PublickKey()) {
		t.Logf("Certificate verification error")
	}
	p, v := ra.Auth(prefix, msg, sk, pk, cert, ra.RA.Mpk())
	if !ra.Verify(prefix, msg, ra.RA.Mpk(), p, v) {
		t.Logf("Verification failed")
	}
}
