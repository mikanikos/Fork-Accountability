package accountability

import (
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/mikanikos/Fork-Accountability/common"
)

// HeightLogs contains all messages for all the rounds from each process in a specific height
type HeightLogs struct {
	logs  map[uint64]*common.HeightVoteSet
	mutex sync.Mutex
}

// NewHeightLogs creates a new HeightLogs structure
func NewHeightLogs() *HeightLogs {
	return &HeightLogs{
		logs: make(map[uint64]*common.HeightVoteSet),
	}
}

// AddHvs adds a new hvs in the height HeightLogs
func (hl *HeightLogs) AddHvs(hvs *common.HeightVoteSet) {
	hl.mutex.Lock()
	defer hl.mutex.Unlock()
	hl.logs[hvs.OwnerID] = hvs
}

// string representation of a HeightLogs
func (hl *HeightLogs) String() string {
	var sb strings.Builder

	sb.WriteString("Height logs are: \n")

	for processID, hvs := range hl.logs {
		sb.WriteString(strconv.FormatUint(processID, 10))
		sb.WriteString(": ")
		sb.WriteString(hvs.String())
		sb.WriteString("\n")
	}

	return sb.String()
}

// Equal is an equality method for HeightLogs
func (hl *HeightLogs) Equal(other *HeightLogs) bool {
	return reflect.DeepEqual(hl, other)
}

// Length returns the length of the HeightLogs
func (hl *HeightLogs) Length() int {
	return len(hl.logs)
}
