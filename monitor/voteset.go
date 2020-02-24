package monitor

import (
	"strings"
)

// VoteSet contains all messages for a specific round
type VoteSet struct {
	round                     uint64
	receivedPrevoteMessages   []*Message
	receivedPrecommitMessages []*Message
	sentPrevoteMessages       []*Message
	sentPrecommitMessages     []*Message
}

// NewVoteSet creates a new VoteSet structure
func NewVoteSet(round uint64) *VoteSet {
	return &VoteSet{
		round:                     round,
		receivedPrevoteMessages:   make([]*Message, 0),
		receivedPrecommitMessages: make([]*Message, 0),
		sentPrevoteMessages:       make([]*Message, 0),
		sentPrecommitMessages:     make([]*Message, 0),
	}
}

// AddSentPrevoteMessage adds a message to the sent prevote messages of a round if not present yet
func (vs *VoteSet) addSentPrevoteMessage(mes *Message) {
	contains := false
	for _, m := range vs.sentPrevoteMessages {
		if mes.equals(m) {
			contains = true
			break
		}
	}
	if !contains {
		vs.sentPrevoteMessages = append(vs.sentPrevoteMessages, mes)
	}
}

// AddSentPrecommitMessage adds a message to the sent precommit messages of a round if not present yet
func (vs *VoteSet) addSentPrecommitMessage(mes *Message) {
	contains := false
	for _, m := range vs.sentPrecommitMessages {
		if mes.equals(m) {
			contains = true
			break
		}
	}
	if !contains {
		vs.sentPrecommitMessages = append(vs.sentPrecommitMessages, mes)
	}
}

func (vs *VoteSet) thereAreQuorumPrevoteMessagesForPrecommit(round, quorum uint64, precommit *Message) bool {
	numberOfAppropriateMessages := uint64(0)
	for _, receivedPrevoteMessage := range vs.receivedPrevoteMessages {
		if receivedPrevoteMessage.equalsRoundValue(precommit) {
			numberOfAppropriateMessages++
		}
	}
	return numberOfAppropriateMessages >= quorum
}

func (vs *VoteSet) String() string {
	var sb strings.Builder

	sb.WriteString(messagesToString("*** RECEIVED PREVOTE MESSAGES ***", vs.receivedPrevoteMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("*** RECEIVED PRECOMMIT MESSAGES ***", vs.receivedPrecommitMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("*** SENT PREVOTE MESSAGES ***", vs.sentPrevoteMessages))
	sb.WriteString("\n")
	sb.WriteString(messagesToString("*** SENT PRECOMMIT MESSAGES ***", vs.sentPrecommitMessages))

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
