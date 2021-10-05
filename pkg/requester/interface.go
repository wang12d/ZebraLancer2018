package requester

import (
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
)

type requester interface {
	// Register to obtain a certificate binding to its public key
	Register(msk []byte) ra.Certificate
	TaskPublish() (pkg.Proof, marlin.VerifyKey)
	// Reward all workers after all data has been collected
	Reward()
}
