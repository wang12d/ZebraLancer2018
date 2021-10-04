package requester

import (
	"log"

	crequester "github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/requester"
	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/utils/reward"
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

type R struct {
	sk   []byte
	pk   []byte
	cert ra.Certificate
	cr   *crequester.Requester
}

func newR(byteSize int) *R {
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
	r.cert = RA.CertGen(r.pk)
	return r.cert
}

// TaskPublish publish a crowdsourcing task to public
func (r *R) TaskPublish(workerRequestered int, reward int64, description string) (pkg.Proof, marlin.VerifyKey) {
	r.cr.Register()
	r.cr.PostTask(workerRequestered, reward, description) // Post the task to blockchain network
	taskAddress := r.cr.Task().Address().Bytes()
	reqeusterAddress := r.cr.Address().Bytes()
	tagPrefix, tagPrefixMsg, zskProof, zskVk := ra.Auth(taskAddress, reqeusterAddress, r.sk, r.pk, r.cert, ra.RA.Mpk())
	return pkg.Proof{
		TagPrefix:    tagPrefix,
		TagPrefixMsg: tagPrefixMsg,
		ZSKProof:     zskProof,
	}, zskVk
}

// Reward awarding all of the workers after submitted their task
func (r *R) Reward(rewardingPolicy reward.Policy) (marlin.Proof, marlin.VerifyKey) {
	r.cr.Rewarding(rewardingPolicy)
	return r.cr.GeneateZKProof()
}
