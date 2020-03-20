package algorithm

import (
	"reflect"
	"strconv"
	"strings"
)

// Faultiness stores information about faulty processes
type Faultiness struct {
	processID uint64
	round     uint64
	reason    error
}

// NewFaultySet creates a new Faultiness structure
func NewFaultiness(processID, round uint64, reason error) *Faultiness {
	return &Faultiness{
		processID: processID,
		round:     round,
		reason:    reason,
	}
}

// string representation of a faultiness
func (fr *Faultiness) String() string {
	var sb strings.Builder

	sb.WriteString("Process ")
	sb.WriteString(strconv.FormatUint(fr.processID, 10))
	sb.WriteString(" at round ")
	sb.WriteString(strconv.FormatUint(fr.round, 10))
	sb.WriteString(": ")
	sb.WriteString(fr.reason.Error())

	return sb.String()
}

// equality for faultiness
func (fr *Faultiness) Equal(other *Faultiness) bool {
	return reflect.DeepEqual(fr, other)
}
