package algorithm

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// FaultySet stores all the validators that are faulty and the corresponding faultiness proofs
type FaultySet struct {
	faultyProcesses map[uint64][]*Faultiness
	Mutex           sync.Mutex
}

// NewFaultySet creates a new FaultySet structure
func NewFaultySet() *FaultySet {
	return &FaultySet{
		faultyProcesses: make(map[uint64][]*Faultiness),
	}
}

// AddFaultinessReason in the faultySet if not already present
func (fs *FaultySet) AddFaultinessReason(fr *Faultiness) {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

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
	var sb strings.Builder

	sb.WriteString("Faulty processes are: \n")

	for processID, reasonsList := range fs.faultyProcesses {
		sb.WriteString(strconv.FormatUint(processID, 10))
		sb.WriteString(": ")

		for _, reason := range reasonsList {
			sb.WriteString(reason.String())
			sb.WriteString("; ")
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// equality for faultySets
func (fs *FaultySet) Equal(other *FaultySet) bool {
	return reflect.DeepEqual(fs, other)
}
