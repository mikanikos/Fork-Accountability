package monitor

var faultyProcesses map[uint64][]*FaultinessReason

// IdentifyFaultyProcesses detects which processes caused the fork and finds all processes that have bad behavior
func IdentifyFaultyProcesses(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsList []*HeightVoteSet) map[uint64][]*FaultinessReason {
	faultyProcesses = make(map[uint64][]*FaultinessReason)
	preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound, hvsList)
	for i := uint64(0); i < numProcesses; i++ {
		checkForFaultiness(numProcesses, firstDecisionRound, secondDecisionRound, hvsList[i])
	}
	return faultyProcesses
}

func preprocessMessages(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvsList []*HeightVoteSet) {
	for round := firstDecisionRound; round <= secondDecisionRound; round++ {
		for processIndex := uint64(0); processIndex < numProcesses; processIndex++ {
			hvs := hvsList[processIndex]
			// We should check whether process sent its voteset at all
			if hvs == nil {
				addFaultinessReason(NewFaultinessReason(processIndex+1, 0, errHVSNotSent))
				continue
			}

			vs, vsLoad := hvs.voteSetMap[round]
			// Maybe some process does not have a voteset for this round
			if vs == nil || !vsLoad {
				continue
			}

			// Processing of the received prevote messages
			addMissingVotes(hvsList, vs.receivedPrevoteMessages)
			// Processing of the received precommit messages
			addMissingVotes(hvsList, vs.receivedPrecommitMessages)
		}
	}
}

func addMissingVotes(hvsList []*HeightVoteSet, receivedMessages []*Message) {
	for _, mes := range receivedMessages {
		senderHeightVoteSet := hvsList[mes.SenderID-1]
		// Sender did not send its voteset -> just ignore it
		if senderHeightVoteSet == nil {
			continue
		}
		senderHeightVoteSet.addMessage(mes)
	}
}

func addFaultinessReason(fr *FaultinessReason) {
	reasons, reasonsLoad := faultyProcesses[fr.processID]
	if reasons == nil || !reasonsLoad {
		faultyProcesses[fr.processID] = make([]*FaultinessReason, 0)
		reasons, _ = faultyProcesses[fr.processID]
	}
	reasons = append(reasons, fr)
	faultyProcesses[fr.processID] = reasons
}

func checkForFaultiness(numProcesses, firstDecisionRound, secondDecisionRound uint64, hvs *HeightVoteSet) {
	quorum := numProcesses - (numProcesses-1)/3 // quorum = 2f + 1
	precommitSent := false
	var precommitValue *BlockPublish
	precommitRound := uint64(0)

	for round := firstDecisionRound; round <= secondDecisionRound; round++ {
		vs, vsLoad := hvs.voteSetMap[round]
		// Maybe some process does not have a voteset for this round
		if vs == nil || !vsLoad {
			continue
		}

		sizeSentPrevotes := len(vs.sentPrevoteMessages)
		if sizeSentPrevotes != 0 {
			// Multiple prevote messages have been sent - faulty behaviour
			if sizeSentPrevotes > 1 {
				addFaultinessReason(NewFaultinessReason(hvs.ownerID, round, errMultiplePrevote))
			} else {
				// Only one prevote message has been sent
				// If the process had previously sent precommit for some value, it can only send prevote message for different value if it has received 2f + 1 (quorum) prevote messages for that value
				if precommitSent {
					message := vs.sentPrevoteMessages[0]
					// Only if two values are not the same, we should look for 2f + 1 prevote messages
					if !message.TxBlock.equals(precommitValue) {
						if !hvs.thereAreQuorumPrevoteMessagesForPrevote(precommitRound, round, quorum, message) {
							addFaultinessReason(NewFaultinessReason(hvs.ownerID, round, errNotEnoughPrevotesForPrevote))
						}
					}
				}
			}
		}

		sizeSentPrecommits := len(vs.sentPrecommitMessages)
		if sizeSentPrecommits != 0 {
			// Multiple precommit messages have been sent - faulty behaviour
			if sizeSentPrecommits > 1 {
				addFaultinessReason(NewFaultinessReason(hvs.ownerID, round, errMultiplePrecommit))
			} else {
				message := vs.sentPrecommitMessages[0]
				if message.TxBlock != nil && !vs.thereAreQuorumPrevoteMessagesForPrecommit(round, quorum, message) {
					addFaultinessReason(NewFaultinessReason(hvs.ownerID, round, errNotEnoughPrevotesForPreoommit))
				}

				// If not nil is precommited
				if message.TxBlock != nil {
					precommitSent = true
					precommitValue = message.TxBlock
					precommitRound = round
				}
			}
		}
	}
}
