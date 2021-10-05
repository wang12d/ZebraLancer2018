package worker

import (
	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/task"
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
)

type worker interface {
	Register() ra.Certificate
	AnswerCollection(t *task.Task, data []byte) (pkg.Proof, marlin.VerifyKey)
}
