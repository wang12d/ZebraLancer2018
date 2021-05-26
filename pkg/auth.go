package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
)

// Auth generates the zero-knowledge proof of the requirement
func Auth(prefix, msg, sk, pk, cert []byte, pp ZskPP) ([]byte, []byte) {
	t1 := pairHash(prefix, sk)
	prefixWithMsg := make([]byte, len(prefix)+len(msg))
	copy(prefixWithMsg[:len(prefix)], prefix)
	copy(prefixWithMsg[len(prefix):], msg)
	t2 := pairHash(prefixWithMsg, sk)
	return t1, t2
}

func pairHash(key, val []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(val)
	return h.Sum(nil)
}