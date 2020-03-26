package common

import (
	"sync"
)

// HeightLogs contains all messages for all the rounds from each process in a specific height
type HeightLogs struct {
	Logs  map[uint64]*HeightVoteSet
	Mutex sync.RWMutex
}

// NewHeightVoteSet creates a new height logs structure
func NewHeightLogs() *HeightLogs {
	return &HeightLogs{
		Logs: make(map[uint64]*HeightVoteSet),
	}
}

// AddHvs adds a new hvs in the height logs
func (hl *HeightLogs) AddHvs(hvs *HeightVoteSet) {
	hl.Mutex.Lock()
	defer hl.Mutex.Unlock()
	hl.Logs[hvs.OwnerID] = hvs
}
