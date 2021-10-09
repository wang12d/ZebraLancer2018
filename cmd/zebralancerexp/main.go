package main

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/task"
	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/utils/encoder"
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
	"github.com/wang12d/ZebraLancer2018/pkg/requester"
	"github.com/wang12d/ZebraLancer2018/pkg/worker"
)

const (
	numberOfIteration = uint64(10)
	reward            = 5000
	mu                = 0
	sigma             = 250
)

var (
	numberOfWorkers = []int{
		3, 5, 7, 9, 11,
	}
)

func main() {
	byteSize := 2048

	length := len(numberOfWorkers)
	requesterProofSize, requesterVerifyKeySize := make([]uint64, length), make([]uint64, length)
	requesterProofGenTime, requesterVerifyTime := make([]time.Duration, length), make([]time.Duration, length)

	workerProofSize, workerVerifyKeySize := make([]uint64, length), make([]uint64, length)
	workerProofGenTime, workerVerifyTime := make([]time.Duration, length), make([]time.Duration, length)

	rewardProofSize, rewardVerifyKeySize := make([]uint64, length), make([]uint64, length)
	rewardProofGenTime, rewardVerifyTime := make([]time.Duration, length), make([]time.Duration, length)

	onChainStorage := make([]uint64, length)

	var timeStart time.Time
	for i := 0; i < int(numberOfIteration); i++ {
		for n, workerRequired := range numberOfWorkers {
			r := requester.NewR(byteSize)
			r.Register(ra.RA)
			taskDescription := "Collecting the time of daliy smartphone usage"
			timeStart = time.Now()
			aulxiliary, proof, vk := r.TaskPublish(workerRequired, reward, taskDescription)
			requesterProofGenTime[n] += time.Since(timeStart)
			requesterProofSize[n] += uint64(len(proof.TagPrefix) + len(proof.TagPrefixMsg) + len(proof.ZSKProof))
			requesterVerifyKeySize[n] += uint64(len(vk))
			timeStart = time.Now()
			verifyResult := ra.Verify(aulxiliary.Prefix, aulxiliary.Msg, ra.RA.Mpk(), proof, vk)
			requesterVerifyTime[n] += time.Since(timeStart)
			if !verifyResult {
				log.Fatalf("Task publish verficiation failed\n")
			}

			workers := make([]*worker.W, workerRequired)
			for i := 0; i < workerRequired; i++ {
				workers[i] = worker.NewW()
				workers[i].Register(ra.RA)
			}

			t, err := r.Task()
			if err != nil {
				log.Fatalf("Obtain task of the requester error: %v\n", err)
			}
			evalResults := make([]marlin.EvaluationResults, workerRequired)
			encryptedData := make([][]byte, workerRequired)
			for ii, w := range workers {
				daliyTime := 3*sigma + rand.Uint64()%5000
				data := encoder.Uint64ToBytes(daliyTime)
				timeStart = time.Now()
				aulxiliary, proof, vk = w.AnswerCollection(t, data)
				workerProofGenTime[n] += time.Since(timeStart)
				workerProofSize[n] += uint64(len(proof.TagPrefix) + len(proof.TagPrefixMsg) + len(proof.ZSKProof))
				workerVerifyKeySize[n] += uint64(len(vk))
				timeStart = time.Now()
				verifyResult = ra.Verify(aulxiliary.Prefix, aulxiliary.Msg, ra.RA.Mpk(), proof, vk)
				workerVerifyTime[n] += time.Since(timeStart)
				evalResults[workerRequired-ii-1] = marlin.EvaluationResults{
					uint64(daliyTime - mu + 3*sigma),
					uint64(daliyTime - mu - 3*sigma),
				}
				encryptedData[workerRequired-ii-1] = w.EncryptedData()
				onChainStorage[n] += uint64(len(encryptedData[workerRequired-ii-1]))
				if !verifyResult {
					log.Fatalf("Worker answer collection verification failed\n")
				}
			}
			onChainStorage[n] += uint64(len(evalResults) * 8)

			workerProofGenTime[n] = workerProofGenTime[n] / time.Duration(workerRequired)
			workerVerifyTime[n] = workerVerifyTime[n] / time.Duration(workerRequired)
			workerProofSize[n] = workerProofSize[n] / uint64(workerRequired)
			workerVerifyKeySize[n] = workerVerifyKeySize[n] / uint64(workerRequired)

			marlinProof, marlinVk, timeCost := r.Reward(&PC{})
			rewardProofGenTime[n] += timeCost
			rewardProofSize[n] += uint64(len(marlinProof))
			rewardVerifyKeySize[n] += uint64(len(marlinVk))
			timeStart = time.Now()
			rewardResult := marlin.VerifyEncryptionZKProof(evalResults, encryptedData, marlinProof, marlinVk)
			rewardVerifyTime[n] += time.Since(timeStart)
			if !rewardResult {
				fmt.Println("Reward verification failed")
			}
		}
	}

	for i := 0; i < length; i++ {
		fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", requesterProofSize[i]/numberOfIteration,
			requesterVerifyKeySize[i]/numberOfIteration,
			requesterProofGenTime[i]/time.Duration(numberOfIteration), requesterVerifyTime[i]/time.Duration(numberOfIteration),
			workerProofSize[i]/numberOfIteration, workerVerifyKeySize[i]/numberOfIteration, workerProofGenTime[i]/time.Duration(numberOfIteration),
			workerVerifyTime[i]/time.Duration(numberOfIteration), rewardProofSize[i]/numberOfIteration, rewardVerifyKeySize[i]/numberOfIteration,
			rewardProofGenTime[i]/time.Duration(numberOfIteration), rewardVerifyTime[i]/time.Duration(numberOfIteration), onChainStorage[i]/numberOfIteration,
		)
	}
}

type PC struct {
}

func (cp *PC) CalculateRewards(t *task.Task, reward *big.Int, workerID int) *big.Int {
	return reward
}
