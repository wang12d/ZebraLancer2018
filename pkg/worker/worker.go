package worker

import (
	"log"

	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/task"
	cworker "github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/worker"
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
)

type W struct {
	sk   []byte
	pk   []byte
	cert ra.Certificate
	cw   *cworker.Worker
}

func newW() *W {
	pk, sk, err := pkg.KeyGeneration()
	if err != nil {
		log.Fatalf("Key generation error: %v\n", err)
	}
	return &W{
		pk:   pk,
		sk:   sk,
		cert: nil,
		cw:   cworker.NewWorker(),
	}
}

// AnswerCollection collects and uploads data to crowdsourcing blockchain
func (w *W) AnswerCollection(t *task.Task, data []byte) (pkg.Proof, marlin.VerifyKey) {
	w.cw.ParticipantTask(t)
	w.cw.CollectData(0, data)
	workerAddress := w.cw.Address().Bytes()
	taskAddress := t.Address().Bytes()
	encryptedData, err := t.Encryptor().EncryptData(data)
	if err != nil {
		log.Fatalf("Worker encrypts data error: %v\n", err)
	}
	msg := make([]byte, len(workerAddress)+len(encryptedData))
	copy(msg[:len(workerAddress)], workerAddress)
	copy(msg[len(workerAddress):], encryptedData)
	return ra.Auth(taskAddress, msg, w.sk, w.pk, w.cert, ra.RA.Mpk())
}
