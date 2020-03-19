package common

import (
	"fmt"
	"strconv"
	"strings"
)

// HeightVoteSet contains all messages for all the rounds of a specific height
type HeightVoteSet struct {
	OwnerID    uint64              `yaml:"id"`
	VoteSetMap map[uint64]*VoteSet `yaml:"heightvoteset"`
}

// NewHeightVoteSet creates a new HeightVoteSet structure
func NewHeightVoteSet(owner uint64) *HeightVoteSet {
	return &HeightVoteSet{
		OwnerID:    owner,
		VoteSetMap: make(map[uint64]*VoteSet),
	}
}

// AddMessage adds a given message to the right voteSet
func (hvs *HeightVoteSet) AddMessage(mes *Message) {

	vs, loaded := hvs.VoteSetMap[mes.Round]
	if vs == nil || !loaded {
		vs = NewVoteSet()
		hvs.VoteSetMap[mes.Round] = vs
	}

	switch mes.Type {
	case Prevote:
		addSentMessage(vs.SentPrevoteMessages, mes)

	case Precommit:
		addSentMessage(vs.SentPrecommitMessages, mes)

	default:
		//  print error
		fmt.Println("Error: message type not known")
	}
}

// ThereAreQuorumPrevoteMessagesForPrevote checks if there are enough prevotes to justify another prevote given a quorum
func (hvs *HeightVoteSet) ThereAreQuorumPrevoteMessagesForPrevote(lockedRound, currentRound, quorum uint64, prevoteMessage *Message) bool {
	for round, voteSet := range hvs.VoteSetMap {
		if (round < lockedRound || round >= currentRound) || (voteSet == nil) {
			continue
		}

		numOfAppropriateMessages := uint64(0)
		for _, receivedPrevoteMessage := range voteSet.ReceivedPrevoteMessages {
			if receivedPrevoteMessage.Value == prevoteMessage.Value {
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

	sb.WriteString("Process " + strconv.FormatUint(hvs.OwnerID, 10))
	sb.WriteString("\n")

	for round, voteSet := range hvs.VoteSetMap {
		sb.WriteString("*** ROUND ")
		sb.WriteString(strconv.FormatUint(round, 10))
		sb.WriteString(" ***")
		sb.WriteString("\n")
		sb.WriteString(voteSet.String())
	}

	return sb.String()
}
