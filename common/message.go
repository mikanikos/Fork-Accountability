package common

import (
	"reflect"
	"strconv"
	"strings"
)

// MessageType represents the type of the message
type MessageType string

const (
	// Prevote type
	Prevote MessageType = "PREVOTE"
	// Precommit type
	Precommit MessageType = "PRECOMMIT"
)

// Message struct
type Message struct {
	Type           MessageType `yaml:"type"`
	SenderID       string      `yaml:"sender"`
	Round          uint64      `yaml:"round"`
	Value          *Value      `yaml:"value"`
	Justifications []*Message  `yaml:"justifications"`
}

// NewMessage creates a new message
func NewMessage(typeMes MessageType, senderID string, round uint64, value *Value, justifications []*Message) *Message {
	return &Message{
		Type:           typeMes,
		SenderID:       senderID,
		Round:          round,
		Value:          value,
		Justifications: justifications,
	}
}

// Equal is the equality method for messages
func (mes *Message) Equal(other *Message) bool {
	return reflect.DeepEqual(mes, other)
}

// String representation of a message
func (mes *Message) String() string {
	var sb strings.Builder
	sb.WriteString(string(mes.Type))
	sb.WriteString(" - Sender: ")
	sb.WriteString(mes.SenderID)
	sb.WriteString(", Round: ")
	sb.WriteString(strconv.FormatUint(mes.Round, 10))
	sb.WriteString(", Value: ")
	sb.WriteString(mes.Value.String())
	sb.WriteString(", Justifications: ")

	if len(mes.Justifications) == 0 {
		sb.WriteString("[]\n")
	}

	currentLength := sb.Len() + 4

	for i, just := range mes.Justifications {
		if i != 0 {
			sb.WriteString(strings.Repeat(" ", currentLength))
		}
		sb.WriteString(string(just.Type))
		sb.WriteString(" - Sender: ")
		sb.WriteString(just.SenderID)
		sb.WriteString(", Round: ")
		sb.WriteString(strconv.FormatUint(just.Round, 10))
		sb.WriteString(", Value: ")
		sb.WriteString(just.Value.String())
		sb.WriteString("\n")
	}

	return sb.String()
}
