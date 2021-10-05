package requester

import (
	"log"

	crequester "github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/requester"
	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/task"
	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/utils/reward"
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
)

type R struct {
	sk   []byte
	pk   []byte
	cert ra.Certificate
	cr   *crequester.Requester
}

func NewR(byteSize int) *R {
	pk, sk, err := pkg.KeyGeneration()
	if err != nil {
		log.Fatalf("Key generation error: %v\n", err)
	}
	return &R{
		sk:   sk,
		pk:   pk,
		cert: nil,
		cr:   crequester.NewRequester(byteSize),
	}
}

// Register to binding a certificate to its public key
func (r *R) Register(RA ra.RegisterAuthority) ra.Certificate {
	var err error
	r.cert, err = RA.CertGen(r.pk)
	if err != nil {
		log.Fatalf("Register obtain certificate error: %v\n", err)
	}
	return r.cert
}

// TaskPublish publish a crowdsourcing task to public
func (r *R) TaskPublish(workerRequired int, reward int64, description string) (pkg.Auxiliary, pkg.Proof, marlin.VerifyKey) {
	r.cr.Register()
	r.cr.PostTask(workerRequired, reward, description) // Post the task to blockchain network
	taskAddress := r.cr.Task().Address().Bytes()
	reqeusterAddress := r.cr.Address().Bytes()
	proof, vk := ra.Auth(taskAddress, reqeusterAddress, r.sk, r.pk, r.cert, ra.RA.Mpk())
	return pkg.Auxiliary{Prefix: taskAddress, Msg: reqeusterAddress}, proof, vk
}

// Reward awarding all of the workers after submitted their task
func (r *R) Reward(rewardingPolicy reward.Policy) (marlin.Proof, marlin.VerifyKey) {
	r.cr.Rewarding(rewardingPolicy)
	return r.cr.GeneateZKProof()
}

// Task return the task published by the requester
func (r *R) Task() (*task.Task, error) {
	if r.cr.Task() == nil {
		return nil, TaskNotPushlied
	}
	return r.cr.Task(), nil
}
