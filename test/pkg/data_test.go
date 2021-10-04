package pkg_test

import (
	"math/rand"
	"testing"
)

const (
	randomDecodeTimes = 1000000
)

func TestRandomDecode(t *testing.T) {
	for i := 0; i < randomDecodeTimes; i++ {
		pp := pkg.ZskPP{
			OutputOne: uint(rand.Uint32()),
			OutputTwo: uint(rand.Uint32()),
		}
		buf := pp.ByteEncode()
		ppDecoded, err := pkg.ByteDecode(buf)
		if err != nil {
			t.Log(err)
		}
		if ppDecoded.OutputOne != pp.OutputOne || ppDecoded.OutputTwo != pp.OutputTwo {
			t.Log("Decode error, not equal!")
		}
	}
}
