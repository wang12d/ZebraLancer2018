package requester

import (
	"errors"

	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/task"
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
)

var (
	TaskNotPushlied = errors.New("The task has not been published yet")
)

type requester interface {
	// Register to obtain a certificate binding to its public key
	Register(msk []byte) ra.Certificate
	TaskPublish() (pkg.Proof, marlin.VerifyKey)
	PublishedTask() (*task.Task, error)
	// Reward all workers after all data has been collected
	Reward()
}
