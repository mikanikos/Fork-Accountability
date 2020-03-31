package common

import (
	"sync"
)

// HeightLogs contains all messages for all the rounds from each process in a specific height
type HeightLogs struct {
	Height uint64
	Logs   map[uint64]*HeightVoteSet
	Mutex  sync.RWMutex
}

// NewHeightVoteSet creates a new height logs structure
func NewHeightLogs(height uint64) *HeightLogs {
	return &HeightLogs{
		Height: height,
		Logs:   make(map[uint64]*HeightVoteSet),
	}
}

// AddHvs adds a new hvs in the height logs
func (hl *HeightLogs) AddHvs(hvs *HeightVoteSet) {
	hl.Mutex.Lock()
	defer hl.Mutex.Unlock()
	hl.Logs[hvs.OwnerID] = hvs
}
