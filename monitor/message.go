package monitor

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

// BlockPublish struct
type BlockPublish struct {
	PrevHash    [32]byte
	Transaction int
}

// Message struct
type Message struct {
	Type     MessageType
	SenderID uint64
	Round    uint64
	TxBlock  *BlockPublish
}

// NewMessage creates a new message
func NewMessage(typeMes MessageType, senderID, round uint64, block *BlockPublish) *Message {
	return &Message{
		Type:     typeMes,
		SenderID: senderID,
		Round:    round,
		TxBlock:  block,
	}
}

func (mes *Message) equals(other *Message) bool {
	return mes.Type == other.Type && mes.SenderID == other.SenderID && mes.Round == other.Round && mes.TxBlock.equals(other.TxBlock)
}

func (mes *Message) equalsRoundValue(other *Message) bool {
	return mes.Round == other.Round && mes.TxBlock.equals(other.TxBlock)
}

func (mes *Message) String() string {
	var sb strings.Builder
	sb.WriteString(string(mes.Type))
	sb.WriteString(" - ")
	sb.WriteString(strconv.FormatUint(mes.SenderID, 10))
	sb.WriteString(", ")
	sb.WriteString(strconv.FormatUint(mes.Round, 10))
	sb.WriteString(", ")
	sb.WriteString(strconv.Itoa(mes.TxBlock.Transaction))
	return sb.String()
}

func (block *BlockPublish) equals(other *BlockPublish) bool {
	return block.PrevHash == other.PrevHash && block.Transaction == other.Transaction
}
