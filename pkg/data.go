package pkg

import "github.com/wang12d/GoMarlin/marlin"

type Auxiliary struct {
	Prefix []byte
	Msg    []byte
}
type Proof struct {
	TagPrefix    []byte
	TagPrefixMsg []byte
	ZSKProof     marlin.Proof
}
