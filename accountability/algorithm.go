package accountability

import (
	"github.com/mikanikos/Fork-Accountability/common"
	"sync"
)

// IdentifyFaultyProcesses detects which processes caused the fork and finds all processes that have bad behavior
func (acc *Accountability) IdentifyFaultyProcesses(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsMap *common.HeightLogs) {
	// first, preprocess messages by scanning all the received vote sets and add missing messages in the processes which omitted sent messages
	acc.preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound, hvsMap)

	// then, find faulty processes by analyzing the message logs
	acc.findFaultyProcesses(numProcesses, firstDecisionRound, secondDecisionRound, hvsMap)
}

// Preprocess messages by scanning all the received vote sets and add missing messages in the respective votes sets of processes which omitted sent messages
func (acc *Accountability) preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsMap *common.HeightLogs) {
	for round := firstDecisionRound; round <= secondDecisionRound; round++ {
		for processIndex := uint64(1); processIndex <= numProcesses; processIndex++ {

			hvs, hvsLoad := hvsMap.Logs[processIndex]
			// if process didn't send the hvs, it's faulty
			if hvs == nil || !hvsLoad {
				acc.AddFaultinessReason(NewFaultiness(processIndex, 0, FaultinessHVSNotSent))
				continue
			}

			vs, vsLoad := hvs.VoteSetMap[round]
			// if process doesn't have a voteset, just ignore
			if vs == nil || !vsLoad {
				continue
			}

			// Processing of the received prevote messages
			acc.addMissingVotes(hvsMap, vs.ReceivedPrevoteMessages)
			// Processing of the received precommit messages
			acc.addMissingVotes(hvsMap, vs.ReceivedPrecommitMessages)
		}
	}
}

// add missing votes to the other processes based on the messages received by the current process
func (acc *Accountability) addMissingVotes(hvsMap *common.HeightLogs, receivedMessages []*common.Message) {
	for _, mes := range receivedMessages {
		senderHeightVoteSet := hvsMap.Logs[mes.SenderID]

		// sender didn't send hvs, it's faulty
		if senderHeightVoteSet == nil {
			acc.AddFaultinessReason(NewFaultiness(mes.SenderID, 0, FaultinessHVSNotSent))
			continue
		}

		// add message if not already present in the sender vote set
		senderHeightVoteSet.AddMessage(mes)
	}
}

// check for faultiness in each process by analyzing the history of messages and making sure it followed the consensus algorithm
func (acc *Accountability) findFaultyProcesses(numProcesses, firstRound, secondRound uint64, hvsMap *common.HeightLogs) {
	wg := sync.WaitGroup{}
	quorum := numProcesses - (numProcesses-1)/3 // quorum = 2f + 1

	// check for faultiness for each process by analyzing the history of messages and making sure it followed the consensus algorithm
	for processId := uint64(1); processId <= numProcesses; processId++ {

		hvs, hvsLoad := hvsMap.Logs[processId]
		// if process didn't send the hvs, ignore because pre-processing already caught that
		if hvs == nil || !hvsLoad {
			continue
		}

		wg.Add(1)
		go acc.isProcessFaulty(quorum, firstRound, secondRound, hvsMap.Logs[processId], &wg)
	}

	wg.Wait()
}

func (acc *Accountability) checkForDuplicateMessages(processId, round uint64, vs *common.VoteSet) {
	// check for duplicates prevotes
	if len(vs.SentPrevoteMessages) > 1 {
		acc.AddFaultinessReason(NewFaultiness(processId, round, FaultinessMultiplePrevotes))
	}

	// check for duplicates precommits
	if len(vs.SentPrecommitMessages) > 1 {
		acc.AddFaultinessReason(NewFaultiness(processId, round, FaultinessMultiplePrecommits))
	}
}

func (acc *Accountability) isProcessFaulty(quorum, firstRound, secondRound uint64, hvs *common.HeightVoteSet, wg *sync.WaitGroup) {
	lockValue := -1
	lockRound := uint64(0)

	for round := firstRound; round <= secondRound; round++ {
		vs, vsLoad := hvs.VoteSetMap[round]
		// if process doesn't have a voteset, just ignore
		if vs == nil || !vsLoad {
			continue
		}

		// check if multiple prevotes/precommits have been sent in the same round
		acc.checkForDuplicateMessages(hvs.OwnerID, round, vs)

		if len(vs.SentPrevoteMessages) == 1 {
			// Only one prevote message has been sent
			// If the process had previously sent precommit for some value, it can only send prevote message for different value if it has received 2f + 1 (quorum) prevote messages for that value
			if lockValue != -1 {
				message := vs.SentPrevoteMessages[0]
				// Only if two values are not the same, we should look for 2f + 1 prevote messages
				if message.Value != lockValue {
					if !hvs.ThereAreQuorumPrevoteMessagesForPrevote(lockRound, round, quorum, message) {
						acc.AddFaultinessReason(NewFaultiness(hvs.OwnerID, round, FaultinessNotEnoughPrevotesForPrevote))
					}
				}
			}
		}

		if len(vs.SentPrecommitMessages) == 1 {
			message := vs.SentPrecommitMessages[0]
			if message.Value != -1 && !vs.ThereAreQuorumPrevoteMessagesForPrecommit(round, quorum, message) {
				acc.AddFaultinessReason(NewFaultiness(hvs.OwnerID, round, FaultinessNotEnoughPrevotesForPrecommit))
			}

			// If not -1 is precommited
			if message.Value != -1 {
				lockValue = message.Value
				lockRound = round
			}
		}
	}

	wg.Done()
}
