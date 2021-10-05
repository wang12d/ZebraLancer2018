package main

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"

	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/task"
	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/utils/encoder"
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
	"github.com/wang12d/ZebraLancer2018/pkg/requester"
	"github.com/wang12d/ZebraLancer2018/pkg/worker"
)

const (
	workerRequired = 10
	reward         = 5000
	mu             = 0
	sigma          = 250
)

func main() {
	byteSize := 2048
	r := requester.NewR(byteSize)
	r.Register(ra.RA)
	taskDescription := "Collecting the time of daliy smartphone usage"
	aulxiliary, proof, vk := r.TaskPublish(workerRequired, reward, taskDescription)
	verifyResult := ra.Verify(aulxiliary.Prefix, aulxiliary.Msg, ra.RA.Mpk(), proof, vk)
	if !verifyResult {
		log.Fatalf("Task publish verficiation failed\n")
	}

	workers := make([]*worker.W, workerRequired)
	for i := 0; i < workerRequired; i++ {
		workers[i] = worker.NewW()
		workers[i].Register()
	}

	t, err := r.Task()
	if err != nil {
		log.Fatalf("Obtain task of the requester error: %v\n", err)
	}
	evalResults := make([]marlin.EvaluationResults, workerRequired)
	encryptedData := make([][]byte, workerRequired)
	for i, w := range workers {
		daliyTime := 3*sigma + rand.Uint64()%5000
		data := encoder.Uint64ToBytes(daliyTime)
		aulxiliary, proof, vk = w.AnswerCollection(t, data)
		verifyResult = ra.Verify(aulxiliary.Prefix, aulxiliary.Msg, ra.RA.Mpk(), proof, vk)
		evalResults[workerRequired-i-1] = marlin.EvaluationResults{
			uint64(daliyTime - mu + 3*sigma),
			uint64(daliyTime - mu - 3*sigma),
		}
		encryptedData[workerRequired-i-1] = w.EncryptedData()
		if !verifyResult {
			log.Fatalf("Worker answer collection verification failed\n")
		}
	}
	marlinProof, marlinVk := r.Reward(&PC{})
	rewardResult := marlin.VerifyEncryptionZKProof(evalResults, encryptedData, marlinProof, marlinVk)
	if !rewardResult {
		fmt.Println("Reward verification failed")
	}
}

type PC struct {
}

func (cp *PC) CalculateRewards(t *task.Task, reward *big.Int, workerID int) *big.Int {
	return reward
}
