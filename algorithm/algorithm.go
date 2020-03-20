package algorithm

import (
	"errors"
	"github.com/mikanikos/Fork-Accountability/common"
)

var (
	ErrHVSNotSent                    = errors.New("the process did not send its HeightVoteSet")
	ErrMultiplePrevote               = errors.New("the process sent more than one PREVOTE message in a round")
	ErrMultiplePrecommit             = errors.New("the process sent more than one PRECOMMIT message in a round")
	ErrNotEnoughPrevotesForPrecommit = errors.New("the process did not receive 2f + 1 PREVOTE messages for a sent PRECOMMIT message to be issued")
	ErrNotEnoughPrevotesForPrevote   = errors.New("the process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued")
)

// faulty set
var faultyProcesses *FaultySet

// IdentifyFaultyProcesses detects which processes caused the fork and finds all processes that have bad behavior
func IdentifyFaultyProcesses(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsMap *common.HeightLogs) *FaultySet {

	faultyProcesses = NewFaultySet()

	// first, preprocess messages by scanning all the received vote sets and add missing messages in the processes which omitted sent messages
	preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound, hvsMap)

	// check for faultiness for each process by analyzing the history of messages and making sure it followed the consensus algorithm
	for i := uint64(1); i <= numProcesses; i++ {
		findFaultiness(numProcesses, firstDecisionRound, secondDecisionRound, hvsMap.Logs[i])
	}

	return faultyProcesses
}

// Preprocess messages by scanning all the received vote sets and add missing messages in the respective votes sets of processes which omitted sent messages
func preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsMap *common.HeightLogs) {
	for round := firstDecisionRound; round <= secondDecisionRound; round++ {
		for processIndex := uint64(1); processIndex <= numProcesses; processIndex++ {

			hvs := hvsMap.Logs[processIndex]
			// if process didn't send the hvs, it's faulty
			if hvs == nil {
				faultyProcesses.AddFaultinessReason(NewFaultiness(processIndex, 0, ErrHVSNotSent))
				continue
			}

			vs, vsLoad := hvs.VoteSetMap[round]
			// if process doesn't have a voteset, just ignore
			if vs == nil || !vsLoad {
				continue
			}

			// Processing of the received prevote messages
			addMissingVotes(hvsMap, vs.ReceivedPrevoteMessages)
			// Processing of the received precommit messages
			addMissingVotes(hvsMap, vs.ReceivedPrecommitMessages)
		}
	}
}

// add missing votes to the other processes based on the messages received by the current process
func addMissingVotes(hvsMap *common.HeightLogs, receivedMessages []*common.Message) {
	for _, mes := range receivedMessages {
		senderHeightVoteSet := hvsMap.Logs[mes.SenderID]

		// sender didn't send hvs, it's faulty
		if senderHeightVoteSet == nil {
			faultyProcesses.AddFaultinessReason(NewFaultiness(mes.SenderID, 0, ErrHVSNotSent))
			continue
		}

		// add message if not already present in the sender vote set
		senderHeightVoteSet.AddMessage(mes)
	}
}

// tries to find all the faultiness reasons in all the rounds given the current hvs of a process
func findFaultiness(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvs *common.HeightVoteSet) {

	if hvs == nil {
		return
	}

	quorum := numProcesses - (numProcesses-1)/3 // quorum = 2f + 1
	precommitSent := false
	precommitValue := -1
	precommitRound := uint64(0)

	for round := firstDecisionRound; round <= secondDecisionRound; round++ {
		vs, vsLoad := hvs.VoteSetMap[round]
		// Maybe some process does not have a voteset for this round
		if vs == nil || !vsLoad {
			continue
		}

		// check for duplicates prevotes
		if len(vs.SentPrevoteMessages) > 1 {
			faultyProcesses.AddFaultinessReason(NewFaultiness(hvs.OwnerID, round, ErrMultiplePrevote))
		}

		if len(vs.SentPrevoteMessages) == 1 {
			// Only one prevote message has been sent
			// If the process had previously sent precommit for some value, it can only send prevote message for different value if it has received 2f + 1 (quorum) prevote messages for that value
			if precommitSent {
				message := vs.SentPrevoteMessages[0]
				// Only if two values are not the same, we should look for 2f + 1 prevote messages
				if message.Value != precommitValue {
					if !hvs.ThereAreQuorumPrevoteMessagesForPrevote(precommitRound, round, quorum, message) {
						faultyProcesses.AddFaultinessReason(NewFaultiness(hvs.OwnerID, round, ErrNotEnoughPrevotesForPrevote))
					}
				}
			}
		}

		// check for duplicates precommits
		if len(vs.SentPrecommitMessages) > 1 {
			faultyProcesses.AddFaultinessReason(NewFaultiness(hvs.OwnerID, round, ErrMultiplePrecommit))
		}

		if len(vs.SentPrecommitMessages) == 1 {
			message := vs.SentPrecommitMessages[0]
			if message.Value != -1 && !vs.ThereAreQuorumPrevoteMessagesForPrecommit(round, quorum, message) {
				faultyProcesses.AddFaultinessReason(NewFaultiness(hvs.OwnerID, round, ErrNotEnoughPrevotesForPrecommit))
			}

			// If not nil is precommited
			if message.Value != -1 {
				precommitSent = true
				precommitValue = message.Value
				precommitRound = round
			}
		}
	}
}
