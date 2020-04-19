package accountability

import (
	"github.com/mikanikos/Fork-Accountability/common"
	"sync"
)

// MAIN ALGORITHM MOVED TO ANOTHER FILE FOR BETTER ORGANIZATION

// Run starts the accountability algorithm to detect which processes caused the fork and finds all processes that had bad behavior
func (acc *Accountability) Run(firstDecisionRound, secondDecisionRound uint64) {

	// lock logs to prevent other additions during the execution
	acc.heightLogs.mutex.Lock()
	defer acc.heightLogs.mutex.Unlock()

	// clear faulty set
	acc.faultySet.Clear()

	// first, preprocess messages by scanning all the received vote sets and add missing messages in the processes which omitted to have sent some messages
	acc.preprocessPhase(firstDecisionRound, secondDecisionRound)

	// then, find faulty processes by analyzing their message logs
	acc.faultDetectionPhase(firstDecisionRound, secondDecisionRound)
}

// Preprocess messages by scanning all the received vote sets and add missing messages in the respective votes sets of processes which omitted to have sent some messages
func (acc *Accountability) preprocessPhase(firstDecisionRound, secondDecisionRound uint64) {
	for _, hvs := range acc.heightLogs.messageLogs {
		for round, vs := range hvs.VoteSetMap {
			if round >= firstDecisionRound && round <= secondDecisionRound {
				// Processing of the received prevote messages
				acc.addMissingVotes(vs.ReceivedPrevoteMessages)
				// Processing of the received precommit messages
				acc.addMissingVotes(vs.ReceivedPrecommitMessages)
			}
		}
	}
}

// Add missing votes to the other processes based on the messages received by the current process
func (acc *Accountability) addMissingVotes(receivedMessages []*common.Message) {
	// add all received messages
	for _, mes := range receivedMessages {
		senderHeightVoteSet, loaded := acc.heightLogs.messageLogs[mes.SenderID]
		if senderHeightVoteSet == nil || !loaded {
			senderHeightVoteSet = common.NewHeightVoteSet()
			acc.heightLogs.messageLogs[mes.SenderID] = senderHeightVoteSet
		}

		// add message if not already present in the sender vote set
		senderHeightVoteSet.AddMessage(mes)
	}
}

// Check for faultiness in each process by analyzing the history of messages and making sure it followed the consensus algorithm
func (acc *Accountability) faultDetectionPhase(firstDecisionRound, secondDecisionRound uint64) {
	wg := sync.WaitGroup{}

	// check for faultiness for each process by analyzing the history of messages and making sure it followed the consensus algorithm
	for processID := range acc.heightLogs.messageLogs {
		wg.Add(1)
		// optimize the execution by running the algorithm concurrently for each process
		go acc.isProcessFaulty(firstDecisionRound, secondDecisionRound, processID, &wg)
	}

	// wait for all the goroutines to complete
	wg.Wait()
}

// Check if a process is faulty in every round and detect all the faultiness reasons for it
func (acc *Accountability) isProcessFaulty(firstDecisionRound, secondDecisionRound uint64, processID string, wg *sync.WaitGroup) {
	var lockedValue *common.Value
	lockedRound := int64(-1)

	hvs := acc.heightLogs.messageLogs[processID]
	isHvsReceived := acc.heightLogs.receivedLogsMap[processID]

	// go from the first to the last round (the order is important)
	for round := firstDecisionRound; round <= secondDecisionRound; round++ {

		vs, vsLoad := hvs.VoteSetMap[round]
		// if process doesn't have a voteset, just go to the next round
		if vs == nil || !vsLoad {
			continue
		}

		// check if process equivocated in the round
		acc.checkForEquivocation(processID, round, vs)

		// if height vote set was not received, we can only check for equivocation
		if isHvsReceived {

			// if only one prevote message has been sent AND the process had previously sent precommit for some value
			if len(vs.SentPrevoteMessages) == 1 && lockedValue != nil {
				message := vs.SentPrevoteMessages[0]

				if message.Value != nil {

					// Only if two values are not the same, we should look for 2f + 1 prevote messages
					if !acc.checkQuorumPrevotesForPrevote(hvs, lockedValue, lockedRound, message) {
						acc.faultySet.AddFaultiness(processID, round, faultinessMissingQuorumForPrevote)
					}
				}
			}

			// if only one precommit message has been sent
			if len(vs.SentPrecommitMessages) == 1 {
				message := vs.SentPrecommitMessages[0]

				if message.Value != nil {

					// we should look for 2f + 1 prevote messages
					if !acc.checkQuorumPrevotesForPrecommit(vs, message) {
						acc.faultySet.AddFaultiness(processID, round, faultinessMissingQuorumForPrecommit)
					}

					// set lock value and lock round
					lockedValue = common.NewValue(message.Value.Data)
					lockedRound = int64(round)
				}
			}
		}
	}

	wg.Done()
}

// check if a process equivocated (sent two or more messages with the same type in the same round but with different values)
func (acc *Accountability) checkForEquivocation(processID string, round uint64, vs *common.VoteSet) {
	// check for duplicates prevotes
	if len(vs.SentPrevoteMessages) > 1 {
		acc.faultySet.AddFaultiness(processID, round, faultinessMultiplePrevotes)
	}

	// check for duplicates precommits
	if len(vs.SentPrecommitMessages) > 1 {
		acc.faultySet.AddFaultiness(processID, round, faultinessMultiplePrecommits)
	}
}

// check if there are enough prevotes to justify a precommit given a quorum
func (acc *Accountability) checkQuorumPrevotesForPrecommit(vs *common.VoteSet, precommit *common.Message) bool {
	numberOfAppropriateMessages := uint64(0)
	for _, receivedPrevoteMessage := range vs.ReceivedPrevoteMessages {
		if receivedPrevoteMessage.Value != nil && receivedPrevoteMessage.Value.Equal(precommit.Value) && receivedPrevoteMessage.Round == precommit.Round {
			numberOfAppropriateMessages++
		}
	}
	return numberOfAppropriateMessages >= acc.getQuorumThreshold()
}

// check if there are enough prevotes to justify another prevote given a quorum
func (acc *Accountability) checkQuorumPrevotesForPrevote(hvs *common.HeightVoteSet, lockedValue *common.Value, lockedRound int64, prevote *common.Message) bool {
	// if not enough justifications, the process is faulty
	if uint64(len(prevote.Justifications)) < acc.getQuorumThreshold() {
		return false
	}

	// go over all justifications provided and check that each one exists and is appropriate
	for _, justification := range prevote.Justifications {
		// if it's not between the lockedRound and the current round, it's not valid according to the consensus algorithm
		if !(prevote.Round >= 0 && justification.Round < prevote.Round &&
			(int64(justification.Round) >= lockedRound || lockedValue.Equal(prevote.Value))) {
			return false
		}

		// load vote set
		vs, vsLoaded := hvs.VoteSetMap[justification.Round]

		// if vote set not present, the justification is not real and, thus, not valid
		if vs == nil || !vsLoaded {
			return false
		}

		foundJustification := false
		for _, receivedPrevoteMessage := range vs.ReceivedPrevoteMessages {
			// find the justification and check that is equal to the one contained in the prevote message and corresponds to the same value
			if receivedPrevoteMessage.Value != nil && receivedPrevoteMessage.Value.Equal(prevote.Value) && justification.Equal(receivedPrevoteMessage) {
				foundJustification = true
				break
			}
		}

		// if not found, justification is fake
		if !foundJustification {
			return false
		}
	}

	return true
}
