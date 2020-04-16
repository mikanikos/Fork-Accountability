package accountability

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// FaultySet stores all the validators that are faulty and the corresponding faultiness proofs
// it uses a complex nested map for efficient additions and for keeping the order of elements during the printing
type FaultySet struct {
	faultinessMap map[string]map[uint64]map[Faultiness]struct{}
	mutex sync.RWMutex
}

// NewFaultySet creates a new FaultySet structure
func NewFaultySet() *FaultySet {
	return &FaultySet{
		faultinessMap: make(map[string]map[uint64]map[Faultiness]struct{}),
	}
}

// AddFaultiness in the FaultySet if not already present
func (fs *FaultySet) AddFaultiness(processID string, round uint64, faultiness Faultiness) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	reasonsForProcess, loaded := fs.faultinessMap[processID]
	// create list of reasons for the process if not present
	if reasonsForProcess == nil || !loaded {
		reasonsForProcess = make(map[uint64]map[Faultiness]struct{})
		fs.faultinessMap[processID] = reasonsForProcess
	}

	reasonsForRound, loaded := reasonsForProcess[round]
	// create list of reasons for the round if not present
	if reasonsForRound == nil || !loaded {
		reasonsForRound = make(map[Faultiness]struct{})
		fs.faultinessMap[processID][round] = reasonsForRound
	}

	_, loaded = reasonsForRound[faultiness]
	if !loaded {
		fs.faultinessMap[processID][round][faultiness] = struct{}{}
	}
}

// string representation of a faulty set
func (fs *FaultySet) String() string {

	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	var sb strings.Builder

	sb.WriteString("RESULTS\n\n")

	sb.WriteString("Faulty processes found: ")
	sb.WriteString(strconv.FormatInt(int64(fs.Length()), 10))
	sb.WriteString("\n\n")

	for processID, reasonsForProcess := range fs.faultinessMap {

		sb.WriteString("*** Process ")
		sb.WriteString(processID)
		sb.WriteString(" ***\n\n")

		for round, reasonsForRound := range reasonsForProcess {

			sb.WriteString("- ROUND ")
			sb.WriteString(strconv.FormatUint(round, 10))
			sb.WriteString("\n")

			for reason := range reasonsForRound {
				sb.WriteString(reason.FaultinessReason())
				sb.WriteString("\n")
			}

			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// Equal is an equality method for FaultySet
func (fs *FaultySet) Equal(other *FaultySet) bool {
	if other == nil {
		return false
	}

	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	return reflect.DeepEqual(fs.faultinessMap, other.faultinessMap)
}

// Length returns the length of the FaultySet
func (fs *FaultySet) Length() int {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	return len(fs.faultinessMap)
}

// Clear removes all elements in the FaultySet
func (fs *FaultySet) Clear() {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	fs.faultinessMap = make(map[string]map[uint64]map[Faultiness]struct{})
}
