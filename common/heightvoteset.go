package common

import (
	"fmt"
	"strconv"
	"strings"
)

// HeightVoteSet contains all messages for all the rounds of a specific height
type HeightVoteSet struct {
	OwnerID    uint64              `yaml:"ownerID"`
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
		vs.addSentMessage(mes)

	case Precommit:
		vs.addSentMessage(mes)

	default:
		//  print error
		fmt.Println("Error: message type not known")
	}
}

// ThereAreQuorumPrevoteMessagesForPrevote checks if there are enough prevotes to justify another prevote given a quorum
func (hvs *HeightVoteSet) ThereAreQuorumPrevoteMessagesForPrevote(lockedRound, currentRound, quorum uint64, prevoteMessage *Message) bool {
	// if not enough justification, the process is faulty
	if uint64(len(prevoteMessage.Justifications)) < quorum {
		return false
	}

	// go over all justifications provided and check that each one exists and is appropriate
	for _, justification := range prevoteMessage.Justifications {
		// if it's not between the lockedRound and the current round, it's not valid according to the Tendermint algorithm
		if justification.Round < lockedRound || justification.Round >= currentRound {
			return false
		}

		// load vote set
		vs, vsLoaded := hvs.VoteSetMap[justification.Round]

		// if vote set not present, the justification is not real and thus not valid
		if vs == nil || !vsLoaded {
			return false
		}

		foundJustification := false
		for _, receivedPrevoteMessage := range vs.ReceivedPrevoteMessages {
			// find the justification anc check that is equal to the one contained in the prevote message and corresponds to the same value
			if receivedPrevoteMessage.Value == prevoteMessage.Value && justification.Equal(receivedPrevoteMessage) {
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

// string representation of a hvs
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
