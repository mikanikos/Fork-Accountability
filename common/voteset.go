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

// add a given message to the correct set of sent messages based on the type
func (vs *VoteSet) addSentMessage(mes *Message) {

	switch mes.Type {
	case Prevote:
		if !contains(vs.SentPrevoteMessages, mes) {
			vs.SentPrevoteMessages = append(vs.SentPrevoteMessages, mes)
		}

	case Precommit:
		if !contains(vs.SentPrecommitMessages, mes) {
			vs.SentPrecommitMessages = append(vs.SentPrecommitMessages, mes)
		}
	}
}

// contains utility for list of messages
func contains(messages []*Message, message *Message) bool {
	contains := false
	for _, m := range messages {
		if message.Equal(m) {
			contains = true
			break
		}
	}
	return contains
}

// String representation of a voteset
func (vs *VoteSet) String() string {
	var sb strings.Builder

	sb.WriteString(messagesToString("[Received prevote messages]", vs.ReceivedPrevoteMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("[Received precommit messages]", vs.ReceivedPrecommitMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("[Sent prevote messages]", vs.SentPrevoteMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("[Sent precommit messages]", vs.SentPrecommitMessages))

	return sb.String()
}

// utility to print list of messages
func messagesToString(description string, messageSet []*Message) string {
	var sb strings.Builder

	sb.WriteString(description)
	sb.WriteString("\n")

	if len(messageSet) == 0 {
		sb.WriteString("\tNo data\n")
	}

	for _, mes := range messageSet {
		sb.WriteString("\t")
		sb.WriteString(mes.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
