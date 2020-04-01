package accountability

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Accountability stores all the validators that are faulty and the corresponding faultiness proofs
type Accountability struct {
	faultyProcesses map[uint64][]*Faultiness
	Mutex           sync.Mutex
}

// NewAccountability creates a new Accountability structure
func NewAccountability() *Accountability {
	return &Accountability{
		faultyProcesses: make(map[uint64][]*Faultiness),
	}
}

// AddFaultinessReason in the faultySet if not already present
func (acc *Accountability) AddFaultinessReason(fr *Faultiness) {
	acc.Mutex.Lock()
	defer acc.Mutex.Unlock()

	reasons, reasonsLoad := acc.faultyProcesses[fr.processID]
	// create list of reasons for the process if not present
	if reasons == nil || !reasonsLoad {
		acc.faultyProcesses[fr.processID] = make([]*Faultiness, 0)
		reasons, _ = acc.faultyProcesses[fr.processID]
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
		acc.faultyProcesses[fr.processID] = reasons
	}
}

// string representation of a faulty set
func (acc *Accountability) String() string {
	var sb strings.Builder

	sb.WriteString("Faulty processes are: \n")

	for processID, reasonsList := range acc.faultyProcesses {
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

// equality for faultySet
func (acc *Accountability) Equal(other *Accountability) bool {
	return reflect.DeepEqual(acc, other)
}

// length of the faultySet
func (acc *Accountability) Length() int {
	return len(acc.faultyProcesses)
}
