package common

import (
	"strconv"
	"strings"
)

// HeightVoteSet contains all messages for all the rounds of a specific height
type HeightVoteSet struct {
	VoteSetMap map[uint64]*VoteSet `yaml:"heightvoteset"`
}

// NewHeightVoteSet creates a new HeightVoteSet structure
func NewHeightVoteSet() *HeightVoteSet {
	return &HeightVoteSet{
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

	vs.addSentMessage(mes)
}

// String representation of a hvs
func (hvs *HeightVoteSet) String() string {
	var sb strings.Builder

	for round, voteSet := range hvs.VoteSetMap {
		sb.WriteString("- ROUND ")
		sb.WriteString(strconv.FormatUint(round, 10))
		sb.WriteString("\n")
		sb.WriteString(voteSet.String())
		sb.WriteString("\n")
	}

	return sb.String()
}

// IsValid checks if the height vote set is valid and contains valid data
func (hvs *HeightVoteSet) IsValid(ID string) bool {
	for round, vs := range hvs.VoteSetMap {
		for _, mess := range vs.SentPrevoteMessages {
			if mess.Type != Prevote || mess.Round != round || mess.SenderID != ID {
				return false
			}
		}

		for _, mess := range vs.ReceivedPrevoteMessages {
			if mess.Type != Prevote || mess.Round != round {
				return false
			}
		}

		for _, mess := range vs.SentPrecommitMessages {
			if mess.Type != Precommit || mess.Round != round || mess.SenderID != ID {
				return false
			}
		}

		for _, mess := range vs.ReceivedPrecommitMessages {
			if mess.Type != Precommit || mess.Round != round {
				return false
			}
		}
	}

	return true
}
