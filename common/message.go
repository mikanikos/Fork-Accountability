package common

import (
	"fmt"
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
	Value          uint64      `yaml:"value"`
	Justifications []*Message  `yaml:"justifications"`
}

// NewMessage creates a new message
func NewMessage(typeMes MessageType, senderID string, round uint64, value uint64, justifications []*Message) *Message {
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
	sb.WriteString(" - ID: ")
	sb.WriteString(mes.SenderID)
	sb.WriteString(", Round: ")
	sb.WriteString(strconv.FormatUint(mes.Round, 10))
	sb.WriteString(", Value: ")
	sb.WriteString(strconv.FormatUint(mes.Value, 10))
	sb.WriteString(", Justifications: ")
	sb.WriteString(fmt.Sprint(mes.Justifications))
	return sb.String()
}
