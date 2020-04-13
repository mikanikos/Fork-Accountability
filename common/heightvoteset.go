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
		sb.WriteString("*** ROUND ")
		sb.WriteString(strconv.FormatUint(round, 10))
		sb.WriteString(" ***")
		sb.WriteString("\n")
		sb.WriteString(voteSet.String())
		sb.WriteString("\n")
	}

	return sb.String()
}
