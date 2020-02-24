package monitor

import (
	"fmt"
	"strconv"
	"strings"
)

// HeightVoteSet contains all messages for all the rounds of a specific height
type HeightVoteSet struct {
	ownerID    uint64
	voteSetMap map[uint64]*VoteSet
}

// NewHeightVoteSet creates a new HeightVoteSet structure
func NewHeightVoteSet(owner uint64) *HeightVoteSet {
	return &HeightVoteSet{
		ownerID:    owner,
		voteSetMap: make(map[uint64]*VoteSet),
	}
}

func (hvs *HeightVoteSet) addMessage(mes *Message) {
	switch mes.Type {
	case prevote:
		hvs.voteSetMap[mes.Round].addSentPrevoteMessage(mes)

	case precommit:
		hvs.voteSetMap[mes.Round].addSentPrecommitMessage(mes)

	default:
		//  print error
		fmt.Println("Error: message type not known")
	}
}

func (hvs *HeightVoteSet) thereAreQuorumPrevoteMessagesForPrevote(lockedRound, currentRound, quorum uint64, prevoteMessage *Message) bool {
	for round, voteSet := range hvs.voteSetMap {
		if (round < lockedRound || round >= currentRound) || (voteSet == nil) {
			continue
		}

		numOfAppropriateMessages := uint64(0)
		for _, receivedPrevoteMessage := range voteSet.receivedPrevoteMessages {
			if receivedPrevoteMessage.TxBlock.equals(prevoteMessage.TxBlock) {
				numOfAppropriateMessages++
			}
		}

		if numOfAppropriateMessages >= quorum {
			return true
		}
	}

	return false
}

func (hvs *HeightVoteSet) String() string {
	var sb strings.Builder

	sb.WriteString("Process " + strconv.FormatUint(hvs.ownerID, 10))
	sb.WriteString("\n")

	for round, voteSet := range hvs.voteSetMap {
		sb.WriteString("*** ROUND ")
		sb.WriteString(strconv.FormatUint(round, 10))
		sb.WriteString(" ***")
		sb.WriteString("\n")
		sb.WriteString(voteSet.String())
	}

	return sb.String()
}
