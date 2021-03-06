package connection

import "github.com/mikanikos/Fork-Accountability/common"

// Packet is a general packet exchanged by validators and monitor
type Packet struct {
	Code   uint32
	ID     string
	Height uint64
	Hvs    *common.HeightVoteSet
}

// main whisper protocol parameters, from official specs
const (
	debug = true

	HvsRequest  = 0
	HvsResponse = 1
	HvsMissing  = 3

	// lengths in bytes
	maxBufferSize = 60000

	maxChannelSize = 100

	readDeadline  = 20
	writeDeadline = 20
)
