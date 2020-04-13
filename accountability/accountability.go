package accountability

import (
	"github.com/mikanikos/Fork-Accountability/common"
	"strings"
	"sync"
)

// Accountability stores all the validators that are faulty and the corresponding faultiness proofs
type Accountability struct {
	MessageLogsReceived uint64
	HeightLogs          *HeightLogs
	faultySet           *FaultySet
}

// NewAccountability creates a new Accountability structure
func NewAccountability() *Accountability {
	return &Accountability{
		HeightLogs: NewHeightLogs(),
		faultySet:  NewFaultySet(),
	}
}

// String returns a string representation (result) of the accountability algorithm
func (acc *Accountability) String() string {
	var sb strings.Builder
	sb.WriteString("Accountability algorithm report\n\n")
	sb.WriteString(acc.HeightLogs.String())
	sb.WriteString(acc.faultySet.String())
	return sb.String()
}

// IsCompleted returns true if the algorithm has completed, false otherwise
func (acc *Accountability) IsCompleted(threshold int) bool {
	// if we have at least f + 1 faulty processes, the algorithm completed
	return acc.faultySet.Length() >= threshold
}

// CanRun returns true if the algorithm can run after the update, false otherwise
func (acc *Accountability) CanRun(threshold int) bool {
	// if we have delivered at least f + 1 message logs, run the monitor algorithm
	return int(acc.MessageLogsReceived) >= threshold
}

// Run starts the accountability algorithm to detect which processes caused the fork and finds all processes that have bad behavior
func (acc *Accountability) Run(numProcesses, firstDecisionRound, secondDecisionRound uint64) {

	// lock logs to prevent other additions during the execution
	acc.HeightLogs.mutex.Lock()
	defer acc.HeightLogs.mutex.Unlock()

	// clear faulty set
	acc.faultySet.Clear()

	// first, preprocess messages by scanning all the received vote sets and add missing messages in the processes which omitted to have sent some messages
	acc.preprocessMessages(firstDecisionRound, secondDecisionRound)

	// then, find faulty processes by analyzing their message logs
	acc.findFaultyProcesses(numProcesses, firstDecisionRound, secondDecisionRound)
}

// Preprocess messages by scanning all the received vote sets and add missing messages in the respective votes sets of processes which omitted to have sent some messages
func (acc *Accountability) preprocessMessages(firstDecisionRound, secondDecisionRound uint64) {
	for _, hvs := range acc.HeightLogs.logs {
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
		senderHeightVoteSet, loaded := acc.HeightLogs.logs[mes.SenderID]
		if senderHeightVoteSet == nil || !loaded {
			senderHeightVoteSet = common.NewHeightVoteSet()
			acc.HeightLogs.logs[mes.SenderID] = senderHeightVoteSet
		}

		// add message if not already present in the sender vote set
		senderHeightVoteSet.AddMessage(mes)
	}
}

// Check for faultiness in each process by analyzing the history of messages and making sure it followed the consensus algorithm
func (acc *Accountability) findFaultyProcesses(numProcesses, firstDecisionRound, secondDecisionRound uint64) {
	wg := sync.WaitGroup{}
	quorum := numProcesses - (numProcesses-1)/3 // quorum = 2f + 1

	// check for faultiness for each process by analyzing the history of messages and making sure it followed the consensus algorithm
	for processID, hvs := range acc.HeightLogs.logs {
		wg.Add(1)
		// optimize the execution by running the algorithm concurrently for each process
		go acc.isProcessFaulty(quorum, firstDecisionRound, secondDecisionRound, processID, hvs, &wg)
	}

	// wait for all the goroutines to complete
	wg.Wait()
}

// Check if a process is faulty in every round and detect all the faultiness reasons for it
func (acc *Accountability) isProcessFaulty(quorum, firstDecisionRound, secondDecisionRound uint64, processID string, hvs *common.HeightVoteSet, wg *sync.WaitGroup) {
	lockValue := -1        // null value, we assume positive values for messages
	lockRound := uint64(0) // null round, rounds start from 1

	// go from the first to the last round (the order is important)
	for round := firstDecisionRound; round <= secondDecisionRound; round++ {

		vs, vsLoad := hvs.VoteSetMap[round]
		// if process doesn't have a voteset, just go to the next round
		if vs == nil || !vsLoad {
			continue
		}

		// check for duplicates prevotes
		sentPrevotesLength := len(vs.SentPrevoteMessages)
		if sentPrevotesLength > 1 {
			acc.faultySet.AddFaultinessReason(NewFaultiness(processID, round, faultinessMultiplePrevotes))
		} else {

			// if only one prevote message has been sent AND
			// if the process had previously sent precommit for some value, it can only send prevote message for different value if it has received 2f + 1 (quorum) prevote messages for that value
			if sentPrevotesLength == 1 && lockValue != -1 {
				message := vs.SentPrevoteMessages[0]

				// Only if two values are not the same, we should look for 2f + 1 prevote messages
				if int(message.Value) != lockValue && !checkQuorumPrevotesForPrevote(hvs, lockRound, round, quorum, message) {
					acc.faultySet.AddFaultinessReason(NewFaultiness(processID, round, faultinessMissingQuorumForPrevote))
				}
			}
		}

		sentPrecommitsLength := len(vs.SentPrecommitMessages)
		// check for duplicates precommits
		if sentPrecommitsLength > 1 {
			acc.faultySet.AddFaultinessReason(NewFaultiness(processID, round, faultinessMultiplePrecommits))
		} else {

			// if only one precommit message has been sent
			if sentPrecommitsLength == 1 {
				message := vs.SentPrecommitMessages[0]

				// we should look for 2f + 1 precommit messages
				if !checkQuorumPrevotesForPrecommit(vs, round, quorum, message) {
					acc.faultySet.AddFaultinessReason(NewFaultiness(processID, round, faultinessMissingQuorumForPrecommit))
				}

				// set lock value and lock round
				lockValue = int(message.Value)
				lockRound = round
			}
		}
	}

	wg.Done()
}

// check if there are enough prevotes to justify a precommit given a quorum
func checkQuorumPrevotesForPrecommit(vs *common.VoteSet, round uint64, quorum uint64, precommit *common.Message) bool {
	numberOfAppropriateMessages := uint64(0)
	for _, receivedPrevoteMessage := range vs.ReceivedPrevoteMessages {
		if receivedPrevoteMessage.Value == precommit.Value && receivedPrevoteMessage.Round == round {
			numberOfAppropriateMessages++
		}
	}
	return numberOfAppropriateMessages >= quorum
}

// check if there are enough prevotes to justify another prevote given a quorum
func checkQuorumPrevotesForPrevote(hvs *common.HeightVoteSet, lockedRound, currentRound, quorum uint64, prevote *common.Message) bool {
	// if not enough justifications, the process is faulty
	if uint64(len(prevote.Justifications)) < quorum {
		return false
	}

	// go over all justifications provided and check that each one exists and is appropriate
	for _, justification := range prevote.Justifications {
		// if it's not between the lockedRound and the current round, it's not valid according to the consensus algorithm
		if justification.Round < lockedRound || justification.Round >= currentRound {
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
			if receivedPrevoteMessage.Value == prevote.Value && justification.Equal(receivedPrevoteMessage) {
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
