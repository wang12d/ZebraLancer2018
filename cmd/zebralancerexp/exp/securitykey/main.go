package main

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/task"
	"github.com/wang12d/Go-Crowdsourcing-DApp/pkg/crowdsourcing/utils/encoder"
	"github.com/wang12d/GoMarlin/marlin"
	"github.com/wang12d/ZebraLancer2018/pkg/ra"
	"github.com/wang12d/ZebraLancer2018/pkg/requester"
	"github.com/wang12d/ZebraLancer2018/pkg/worker"
)

const (
	reward = 5000
	mu     = 1000
	sigma  = 250
)

var ()

func main() {
	numberOfIteration := 5

	numberOfWorkers := make([]int, 0)
	byteSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Convert to number err: %v\n", err)
	}
	for _, str := range os.Args[2:] {
		nw, err := strconv.Atoi(str)
		if err != nil {
			log.Fatalf("Conver to number err: %v\n", err)
		}
		numberOfWorkers = append(numberOfWorkers, nw)
	}
	length := len(numberOfWorkers)
	registerationTimeCost, taskPublicationTimeCost := make([]time.Duration, length), make([]time.Duration, length)
	requesterRegistrationTimeCost, workerRegistrationTimeCost := make([]time.Duration, length), make([]time.Duration, length)
	taskParticipationTimeCost, dataCollectionTimeCost := make([]time.Duration, length), make([]time.Duration, length)
	rewardingTimeCost := make([]time.Duration, length)
	onChainStorage := make([]int, length)
	communication := make([]int, length)

	var registerationCost, taskPublicationCost, taskParticipationCost, dataCollectionCost, rewardingCost time.Duration
	var requesterRegistrationCost, workerRegistrationCost time.Duration
	var onChainBytes, communicationCost int

	var timeStart time.Time
	for i := 0; i < int(numberOfIteration); i++ {
		for n, workerRequired := range numberOfWorkers {
			onChainBytes, registerationCost, taskPublicationCost, taskParticipationCost, dataCollectionCost, rewardingCost = 0, 0, 0, 0, 0, 0
			requesterRegistrationCost, workerRegistrationCost = 0, 0
			communicationCost = 0
			r := requester.NewR(byteSize)
			timeStart = time.Now()
			r.Register(ra.RA)
			requesterRegistrationCost += time.Since(timeStart)
			registerationCost += requesterRegistrationCost
			requesterRegistrationTimeCost[n] += requesterRegistrationCost
			taskDescription := "Collecting the time of daliy smartphone usage"
			timeStart = time.Now()
			aulxiliary, proof, vk := r.TaskPublish(workerRequired, reward, taskDescription)
			taskPublicationCost += time.Since(timeStart)
			onChainBytes += len(proof.TagPrefix) + len(proof.TagPrefixMsg) + len(proof.ZSKProof) + len(vk)
			communicationCost += len(proof.TagPrefix) + len(proof.TagPrefixMsg) + len(proof.ZSKProof) + len(vk)

			timeStart = time.Now()
			verifyResult := ra.Verify(aulxiliary.Prefix, aulxiliary.Msg, ra.RA.Mpk(), proof, vk)
			taskPublicationCost += time.Since(timeStart)
			if !verifyResult {
				log.Fatalf("Task publish verficiation failed\n")
			}
			taskPublicationTimeCost[n] += taskPublicationCost

			workers := make([]*worker.W, workerRequired)
			for ii := 0; ii < workerRequired; ii++ {
				workers[ii] = worker.NewW()
				timeStart = time.Now()
				workers[ii].Register(ra.RA)
				workerRegistrationCost += time.Since(timeStart)
			}
			registerationCost += workerRegistrationCost
			workerRegistrationTimeCost[n] += workerRegistrationCost / time.Duration(workerRequired)
			registerationTimeCost[n] += registerationCost / time.Duration(workerRequired+1)

			t, err := r.Task()
			if err != nil {
				log.Fatalf("Obtain task of the requester error: %v\n", err)
			}
			evalResults := make([]marlin.EvaluationResults, workerRequired)
			encryptedData := make([][]byte, workerRequired)
			for ii, w := range workers {
				daliyTime := mu + 3*sigma + rand.Uint64()%5000
				fmt.Printf("Go data: %v\n", daliyTime)
				data := encoder.Uint64ToBytes(daliyTime)

				timeStart = time.Now()
				w.ParticipantTask(t)
				taskParticipationCost += time.Since(timeStart)

				timeStart = time.Now()
				aulxiliary, proof, vk = w.AnswerCollection(t, data)
				dataCollectionCost += time.Since(timeStart)
				if ii == 0 {
					communicationCost += len(proof.TagPrefix) + len(proof.TagPrefixMsg) + len(proof.ZSKProof)
				}
				onChainBytes += len(proof.TagPrefix) + len(proof.TagPrefixMsg) + len(proof.ZSKProof)
				// onChainBytes += len(vk)	// Communication Cost
				timeStart = time.Now()
				verifyResult = ra.Verify(aulxiliary.Prefix, aulxiliary.Msg, ra.RA.Mpk(), proof, vk)
				dataCollectionCost += time.Since(timeStart)
				evalResults[workerRequired-ii-1] = marlin.EvaluationResults{
					uint64(daliyTime - mu + 3*sigma),
					uint64(daliyTime - mu - 3*sigma),
				}
				encryptedData[workerRequired-ii-1] = w.EncryptedData()
				onChainBytes += len(encryptedData[workerRequired-ii-1])
				if ii == 0 {
					communicationCost += len(encryptedData[workerRequired-ii-1])
				}
				if !verifyResult {
					log.Fatalf("Worker answer collection verification failed\n")
				}
			}

			taskParticipationTimeCost[n] += taskParticipationCost / time.Duration(workerRequired)
			dataCollectionTimeCost[n] += dataCollectionCost / time.Duration(workerRequired)

			marlinProof, marlinVk, timeCost := r.Reward(&PC{})
			rewardingCost = timeCost
			rewardingTimeCost[n] += rewardingCost
			onChainBytes += len(marlinProof) + len(marlinVk)
			communicationCost += len(marlinProof) + len(marlinVk)

			onChainStorage[n] += onChainBytes
			communication[n] += communicationCost
			fmt.Printf("Go evalReuslt: %v\n Go encrypted: %v\n", evalResults, encryptedData)
			rewardResult := marlin.VerifyEncryptionZKProof(evalResults, encryptedData, marlinProof, marlinVk)
			if !rewardResult {
				log.Fatalln("Reward verification failed")
			}
		}
	}

	fmt.Printf("ZebraLancer: %v\n", byteSize)
	for i := 0; i < length; i++ {
		fmt.Printf("%v(%v, %v),%v,%v,%v,%v,%v,%v\n", registerationTimeCost[i]/time.Duration(numberOfIteration),
			requesterRegistrationTimeCost[i]/time.Duration(numberOfIteration), workerRegistrationTimeCost[i]/time.Duration(numberOfIteration),
			taskPublicationTimeCost[i]/time.Duration(numberOfIteration), taskParticipationTimeCost[i]/time.Duration(numberOfIteration),
			dataCollectionTimeCost[i]/time.Duration(numberOfIteration), rewardingTimeCost[i]/time.Duration(numberOfIteration),
			float64(onChainStorage[i])/1024.0/float64(numberOfIteration), float64(communication[i])/1024.0/float64(numberOfIteration),
		)
	}
}

type PC struct {
}

func (cp *PC) CalculateRewards(t *task.Task, reward *big.Int, workerID int) *big.Int {
	return reward
}
