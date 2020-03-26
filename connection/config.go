package connection

import "github.com/mikanikos/Fork-Accountability/common"

// Packet is a general packet exchanged by validators and monitor
type Packet struct {
	Code uint32
	Hvs  *common.HeightVoteSet
}

// main whisper protocol parameters, from official specs
const (
	HvsRequest  = 0
	HvsResponse = 1

	// lengths in bytes
	maxBufferSize = 1024
)
