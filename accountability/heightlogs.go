package accountability

import (
	"strings"
	"sync"

	"github.com/mikanikos/Fork-Accountability/common"
)

// HeightLogs contains all messages for all the rounds from each process in a specific height
type HeightLogs struct {
	messageLogs     map[string]*common.HeightVoteSet
	receivedLogsMap map[string]bool
	mutex           sync.RWMutex
}

// NewHeightLogs creates a new HeightLogs structure
func NewHeightLogs() *HeightLogs {
	return &HeightLogs{
		messageLogs:     make(map[string]*common.HeightVoteSet),
		receivedLogsMap: make(map[string]bool),
	}
}

// AddHvs adds a new hvs in the height HeightLogs
func (hl *HeightLogs) AddHvs(processID string, hvs *common.HeightVoteSet) bool {
	hl.mutex.Lock()
	defer hl.mutex.Unlock()

	value, _ := hl.receivedLogsMap[processID]
	if !value {
		hl.messageLogs[processID] = hvs
		hl.receivedLogsMap[processID] = true
		return true
	}

	return false
}

// string representation of a HeightLogs
func (hl *HeightLogs) String() string {
	hl.mutex.RLock()
	defer hl.mutex.RUnlock()

	var sb strings.Builder

	sb.WriteString("HEIGHT LOGS\n\n")

	for processID, hvs := range hl.messageLogs {
		sb.WriteString("*** Process ")
		sb.WriteString(processID)
		sb.WriteString(" ***\n\n")
		sb.WriteString(hvs.String())

		sb.WriteString("------------------------------------------------------------------------------------------------------------------------\n\n")
	}

	return sb.String()
}

// Length returns the length of the HeightLogs
func (hl *HeightLogs) Length() int {
	hl.mutex.RLock()
	defer hl.mutex.RUnlock()
	return len(hl.messageLogs)
}

// ReceivedLength returns the number of received logs so far
func (hl *HeightLogs) ReceivedLength() int {
	hl.mutex.RLock()
	defer hl.mutex.RUnlock()
	numReceived := 0
	for _, val := range hl.receivedLogsMap {
		if val {
			numReceived++
		}
	}
	return numReceived
}
