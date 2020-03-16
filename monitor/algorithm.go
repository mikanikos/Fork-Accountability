package monitor

import (
	"github.com/mikanikos/Fork-Accountability/common"
)

var faultyProcesses map[uint64][]*FaultinessReason

// IdentifyFaultyProcesses detects which processes caused the fork and finds all processes that have bad behavior
func IdentifyFaultyProcesses(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsMap map[uint64]*common.HeightVoteSet) map[uint64][]*FaultinessReason {
	faultyProcesses = make(map[uint64][]*FaultinessReason)
	preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound, hvsMap)
	for i := uint64(1); i <= numProcesses; i++ {
		checkForFaultiness(numProcesses, firstDecisionRound, secondDecisionRound, hvsMap[i])
	}
	return faultyProcesses
}

func preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsMap map[uint64]*common.HeightVoteSet) {
	for round := firstDecisionRound; round <= secondDecisionRound; round++ {
		for processIndex := uint64(1); processIndex <= numProcesses; processIndex++ {
			hvs := hvsMap[processIndex]
			// We should check whether process sent its voteset at all
			if hvs == nil {
				addFaultinessReason(NewFaultinessReason(processIndex, 0, errHVSNotSent))
				continue
			}

			vs, vsLoad := hvs.VoteSetMap[round]
			// Maybe some process does not have a voteset for this round
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

func addMissingVotes(hvsList map[uint64]*common.HeightVoteSet, receivedMessages []*common.Message) {
	for _, mes := range receivedMessages {
		senderHeightVoteSet := hvsList[mes.SenderID]
		// Sender did not send its voteset -> just ignore it
		if senderHeightVoteSet == nil {
			continue
		}
		senderHeightVoteSet.AddMessage(mes)
	}
}

func addFaultinessReason(fr *FaultinessReason) {
	reasons, reasonsLoad := faultyProcesses[fr.processID]
	if reasons == nil || !reasonsLoad {
		faultyProcesses[fr.processID] = make([]*FaultinessReason, 0)
		reasons, _ = faultyProcesses[fr.processID]
	}
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

		sizeSentPrevotes := len(vs.SentPrevoteMessages)
		if sizeSentPrevotes != 0 {
			// Multiple prevote messages have been sent - faulty behaviour
			if sizeSentPrevotes > 1 {
				addFaultinessReason(NewFaultinessReason(hvs.OwnerID, round, errMultiplePrevote))
			} else {
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
		}

		sizeSentPrecommits := len(vs.SentPrecommitMessages)
		if sizeSentPrecommits != 0 {
			// Multiple precommit messages have been sent - faulty behaviour
			if sizeSentPrecommits > 1 {
				addFaultinessReason(NewFaultinessReason(hvs.OwnerID, round, errMultiplePrecommit))
			} else {
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
}
