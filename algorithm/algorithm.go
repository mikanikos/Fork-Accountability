package algorithm

import (
	"github.com/mikanikos/Fork-Accountability/common"
)

var faultyProcesses map[uint64][]*FaultinessReason

// IdentifyFaultyProcesses detects which processes caused the fork and finds all processes that have bad behavior
func IdentifyFaultyProcesses(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsMap map[uint64]*common.HeightVoteSet) map[uint64][]*FaultinessReason {
	faultyProcesses = make(map[uint64][]*FaultinessReason)

	// first, preprocess messages by scanning all the received vote sets and add missing messages in the processes which omitted sent messages
	preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound, hvsMap)

	// check for faultiness for each process by analyzing the history of messages and making sure it followed the consensus algorithm
	for i := uint64(1); i <= numProcesses; i++ {
		checkForFaultiness(numProcesses, firstDecisionRound, secondDecisionRound, hvsMap[i])
	}

	return faultyProcesses
}

/**
	Preprocess messages by scanning all the received vote sets and add missing messages in the respective votes sets of processes which omitted sent messages
*/
func preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsMap map[uint64]*common.HeightVoteSet) {
	for round := firstDecisionRound; round <= secondDecisionRound; round++ {
		for processIndex := uint64(1); processIndex <= numProcesses; processIndex++ {

			hvs := hvsMap[processIndex]
			// if process didn't send the hvs, it's faulty
			if hvs == nil {
				addFaultinessReason(NewFaultinessReason(processIndex, 0, errHVSNotSent))
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

func addMissingVotes(hvsMap map[uint64]*common.HeightVoteSet, receivedMessages []*common.Message) {
	for _, mes := range receivedMessages {
		senderHeightVoteSet := hvsMap[mes.SenderID]

		// sender didn't send hvs, it's faulty
		if senderHeightVoteSet == nil {
			addFaultinessReason(NewFaultinessReason(mes.SenderID, 0, errHVSNotSent))
			continue
		}

		 // add message if not already present in the sender vote set
		senderHeightVoteSet.AddMessage(mes)
		voteSet := senderHeightVoteSet.VoteSetMap[mes.Round]

		// if the process sent more than 1 prevote or precommit, it's faulty
		switch mes.Type {
		case common.Prevote:
			if len(voteSet.SentPrevoteMessages) > 1 {
				addFaultinessReason(NewFaultinessReason(mes.SenderID, mes.Round, errMultiplePrevote))
			}

		case common.Precommit:
			if len(voteSet.SentPrecommitMessages) > 1 {
				addFaultinessReason(NewFaultinessReason(mes.SenderID, mes.Round, errMultiplePrecommit))
			}
		}
	}
}

/**
	Add Faultiness reason in the faulty processes map if not already present
 */
func addFaultinessReason(fr *FaultinessReason) {
	reasons, reasonsLoad := faultyProcesses[fr.processID]
	// create list of reasons for the process if not present
	if reasons == nil || !reasonsLoad {
		faultyProcesses[fr.processID] = make([]*FaultinessReason, 0)
		reasons, _ = faultyProcesses[fr.processID]
	}

	// check if already present
	contains := false
	for _, f := range reasons {
		if f.equals(fr) {
			contains = true
			break
		}
	}
	if !contains {
		reasons = append(reasons, fr)
		faultyProcesses[fr.processID] = reasons
	}
}


func checkForFaultiness(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvs *common.HeightVoteSet) {

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

		if len(vs.SentPrevoteMessages) == 1 {
			// Only one prevote message has been sent
			// If the process had previously sent precommit for some value, it can only send prevote message for different value if it has received 2f + 1 (quorum) prevote messages for that value
			if precommitSent {
				message := vs.SentPrevoteMessages[0]
				// Only if two values are not the same, we should look for 2f + 1 prevote messages
				if message.Value != precommitValue {
					if !hvs.ThereAreQuorumPrevoteMessagesForPrevote(precommitRound, round, quorum, message) {
						addFaultinessReason(NewFaultinessReason(hvs.OwnerID, round, errNotEnoughPrevotesForPrevote))
					}
				}
			}
		}

		if len(vs.SentPrecommitMessages) == 1 {
			message := vs.SentPrecommitMessages[0]
			if message.Value != -1 && !vs.ThereAreQuorumPrevoteMessagesForPrecommit(round, quorum, message) {
				addFaultinessReason(NewFaultinessReason(hvs.OwnerID, round, errNotEnoughPrevotesForPrecommit))
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
