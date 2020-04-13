package accountability

import (
	"strings"
	"sync"

	"github.com/mikanikos/Fork-Accountability/common"
)

// HeightLogs contains all messages for all the rounds from each process in a specific height
type HeightLogs struct {
	logs  map[string]*common.HeightVoteSet
	mutex sync.RWMutex
}

// NewHeightLogs creates a new HeightLogs structure
func NewHeightLogs() *HeightLogs {
	return &HeightLogs{
		logs: make(map[string]*common.HeightVoteSet),
	}
}

// AddHvs adds a new hvs in the height HeightLogs
func (hl *HeightLogs) AddHvs(processID string, hvs *common.HeightVoteSet) {
	hl.mutex.Lock()
	defer hl.mutex.Unlock()
	hl.logs[processID] = hvs
}

// string representation of a HeightLogs
func (hl *HeightLogs) String() string {
	hl.mutex.RLock()
	defer hl.mutex.RUnlock()

	var sb strings.Builder

	sb.WriteString("Height logs\n\n")

	for processID, hvs := range hl.logs {
		sb.WriteString("Process ")
		sb.WriteString(processID)
		sb.WriteString("\n\n")
		sb.WriteString(hvs.String())
		sb.WriteString("\n")
	}

	return sb.String()
}

// Length returns the length of the HeightLogs
func (hl *HeightLogs) Length() int {
	hl.mutex.RLock()
	defer hl.mutex.RUnlock()
	return len(hl.logs)
}

// Contains checks if an element in the logs is already present
func (hl *HeightLogs) Contains(id string) bool {
	hl.mutex.RLock()
	defer hl.mutex.RUnlock()
	_, loaded := hl.logs[id]
	return loaded
}
