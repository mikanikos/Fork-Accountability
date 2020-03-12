package common

import (
	"strconv"
	"strings"
)

// MessageType represents the type of the message
type MessageType string

const (
	// Prevote type
	prevote MessageType = "PREVOTE"
	// Precommit type
	precommit MessageType = "PRECOMMIT"
)

// Message struct
type Message struct {
	Type     MessageType
	SenderID uint64
	Round    uint64
	Value    int
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

func (mes *Message) equals(other *Message) bool {
	return mes.Type == other.Type && mes.SenderID == other.SenderID && mes.Round == other.Round && mes.Value == other.Value
}

func (mes *Message) equalsRoundValue(other *Message) bool {
	return mes.Round == other.Round && mes.Value == other.Value
}

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
