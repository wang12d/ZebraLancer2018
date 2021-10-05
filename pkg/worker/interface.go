package worker

import (
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
)

type worker interface {
	Register() ra.Certificate
	AnswerCollection(data []byte) (pkg.Proof, marlin.VerifyKey)
}
