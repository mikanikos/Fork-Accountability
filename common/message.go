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
	SenderID       uint64      `yaml:"sender"`
	Round          uint64      `yaml:"round"`
	Value          int         `yaml:"value"`
	Justifications []*Message  `yaml:"justifications"`
}

// NewMessage creates a new message
func NewMessage(typeMes MessageType, senderID, round uint64, value int, justifications []*Message) *Message {
	return &Message{
		Type:           typeMes,
		SenderID:       senderID,
		Round:          round,
		Value:          value,
		Justifications: justifications,
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
	sb.WriteString(" - ID: ")
	sb.WriteString(strconv.FormatUint(mes.SenderID, 10))
	sb.WriteString(", Round: ")
	sb.WriteString(strconv.FormatUint(mes.Round, 10))
	sb.WriteString(", Value: ")
	sb.WriteString(strconv.Itoa(mes.Value))
	sb.WriteString(", Justifications: ")
	sb.WriteString(fmt.Sprint(mes.Justifications))
	return sb.String()
}
