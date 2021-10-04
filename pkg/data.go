package pkg

import "github.com/wang12d/GoMarlin/marlin"

type Proof struct {
	TagPrefix    []byte
	TagPrefixMsg []byte
	ZSKProof     marlin.Proof
}
