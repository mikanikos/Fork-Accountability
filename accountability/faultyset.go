package accountability

import (
	"reflect"
	"strings"
	"sync"
)

// FaultySet stores all the validators that are faulty and the corresponding faultiness proofs
type FaultySet struct {
	faultyProcesses map[string][]*Faultiness
	mutex           sync.RWMutex
}

// NewFaultySet creates a new FaultySet structure
func NewFaultySet() *FaultySet {
	return &FaultySet{
		faultyProcesses: make(map[string][]*Faultiness),
	}
}

// AddFaultinessReason in the FaultySet if not already present
func (fs *FaultySet) AddFaultinessReason(fr *Faultiness) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	reasons, reasonsLoad := fs.faultyProcesses[fr.processID]
	// create list of reasons for the process if not present
	if reasons == nil || !reasonsLoad {
		fs.faultyProcesses[fr.processID] = make([]*Faultiness, 0)
		reasons, _ = fs.faultyProcesses[fr.processID]
	}

	// check if already present
	contains := false
	for _, f := range reasons {
		if f.Equal(fr) {
			contains = true
			break
		}
	}
	if !contains {
		reasons = append(reasons, fr)
		fs.faultyProcesses[fr.processID] = reasons
	}
}

// string representation of a faulty set
func (fs *FaultySet) String() string {

	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	var sb strings.Builder

	sb.WriteString("Faulty processes detected\n")

	for _, reasonsList := range fs.faultyProcesses {

		for _, reason := range reasonsList {
			sb.WriteString(reason.String())
			sb.WriteString("; ")
		}

		sb.WriteString("\n")
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
	return reflect.DeepEqual(fs.faultyProcesses, other.faultyProcesses)
}

// Length returns the length of the FaultySet
func (fs *FaultySet) Length() int {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	return len(fs.faultyProcesses)
}

// Clear removes all elements in the FaultySet
func (fs *FaultySet) Clear() {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	fs.faultyProcesses = make(map[string][]*Faultiness)
}
