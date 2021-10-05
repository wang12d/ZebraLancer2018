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
	sk            []byte
	pk            []byte
	encryptedData []byte
	cert          ra.Certificate
	cw            *cworker.Worker
}

func NewW() *W {
	pk, sk, err := pkg.KeyGeneration()
	if err != nil {
		log.Fatalf("Key generation error: %v\n", err)
	}
	return &W{
		pk:            pk,
		sk:            sk,
		cert:          nil,
		encryptedData: nil,
		cw:            cworker.NewWorker(),
	}
}

// Register binds the public key of a worker with a certificate from register authority
func (w *W) Register() ra.Certificate {
	w.cert = ra.RA.CertGen(w.pk)
	w.cw.Register(1)
	return w.cert
}

// AnswerCollection collects and uploads data to crowdsourcing blockchain
func (w *W) AnswerCollection(t *task.Task, data []byte) (pkg.Auxiliary, pkg.Proof, marlin.VerifyKey) {
	w.cw.ParticipantTask(t)
	w.cw.CollectData(0, data)
	w.cw.SubmitData(0)
	var err error
	workerAddress := w.cw.Address().Bytes()
	taskAddress := t.Address().Bytes()
	w.encryptedData, err = t.Encryptor().EncryptData(data)
	if err != nil {
		log.Fatalf("Worker encrypts data error: %v\n", err)
	}
	msg := make([]byte, len(workerAddress)+len(w.encryptedData))
	copy(msg[:len(workerAddress)], workerAddress)
	copy(msg[len(workerAddress):], w.encryptedData)
	proof, vk := ra.Auth(taskAddress, msg, w.sk, w.pk, w.cert, ra.RA.Mpk())
	return pkg.Auxiliary{Prefix: taskAddress, Msg: msg}, proof, vk
}

// EncryptedData returns encrypted data collected by the worker
func (w *W) EncryptedData() []byte {
	return w.encryptedData
}
