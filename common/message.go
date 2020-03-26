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
	Type     MessageType `yaml:"type"`
	SenderID uint64      `yaml:"sender"`
	Round    uint64      `yaml:"round"`
	Value    int         `yaml:"value"`
}

// NewMessage creates a new message
func NewMessage(typeMes MessageType, senderID, round uint64, value int) *Message {
	return &Message{
		Type:     typeMes,
		SenderID: senderID,
		Round:    round,
		Value:    value,
	}
}

// equality for messages
func (mes *Message) Equal(other *Message) bool {
	return reflect.DeepEqual(mes, other)
}

// string representation of a message
func (mes *Message) String() string {
	var sb strings.Builder
	sb.WriteString(string(mes.Type))
	sb.WriteString(" - ")
	sb.WriteString(strconv.FormatUint(mes.SenderID, 10))
	sb.WriteString(", ")
	sb.WriteString(strconv.FormatUint(mes.Round, 10))
	sb.WriteString(", ")
	sb.WriteString(strconv.Itoa(mes.Value))
	return sb.String()
}
