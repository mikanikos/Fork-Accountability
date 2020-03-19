package common

import (
	"strings"
)

// VoteSet contains all messages of a process for a specific round
type VoteSet struct {
	ReceivedPrevoteMessages   []*Message `yaml:"received_prevote"`
	ReceivedPrecommitMessages []*Message `yaml:"received_precommit"`
	SentPrevoteMessages       []*Message `yaml:"sent_prevote"`
	SentPrecommitMessages     []*Message `yaml:"sent_precommit"`
}

// NewVoteSet creates a new VoteSet structure
func NewVoteSet() *VoteSet {
	return &VoteSet{
		ReceivedPrevoteMessages:   make([]*Message, 0),
		ReceivedPrecommitMessages: make([]*Message, 0),
		SentPrevoteMessages:       make([]*Message, 0),
		SentPrecommitMessages:     make([]*Message, 0),
	}
}

func addSentMessage(messages []*Message, mes *Message) {
	contains := false
	for _, m := range messages {
		if mes.equals(m) {
			contains = true
			break
		}
	}
	if !contains {
		messages = append(messages, mes)
	}
}

// ThereAreQuorumPrevoteMessagesForPrecommit checks if there are enough prevotes to justify a precommit given a quorum
func (vs *VoteSet) ThereAreQuorumPrevoteMessagesForPrecommit(round uint64, quorum uint64, precommit *Message) bool {
	numberOfAppropriateMessages := uint64(0)
	for _, receivedPrevoteMessage := range vs.ReceivedPrevoteMessages {
		if receivedPrevoteMessage.Value == precommit.Value && receivedPrevoteMessage.Round == round {
			numberOfAppropriateMessages++
		}
	}
	return numberOfAppropriateMessages >= quorum
}

func (vs *VoteSet) String() string {
	var sb strings.Builder

	sb.WriteString(messagesToString("*** RECEIVED PREVOTE MESSAGES ***", vs.ReceivedPrevoteMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("*** RECEIVED PRECOMMIT MESSAGES ***", vs.ReceivedPrecommitMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("*** SENT PREVOTE MESSAGES ***", vs.SentPrevoteMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("*** SENT PRECOMMIT MESSAGES ***", vs.SentPrecommitMessages))

	return sb.String()
}

func messagesToString(description string, messageSet []*Message) string {
	var sb strings.Builder

	sb.WriteString(description)
	sb.WriteString("\n")
	for _, mes := range messageSet {
		sb.WriteString(mes.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
