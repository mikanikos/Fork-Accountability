package algorithm

import (
	"strconv"
	"strings"
)

// FaultinessReason stores information about faulty processes
type FaultinessReason struct {
	processID uint64
	round     uint64
	reason    string
}

// NewFaultinessReason creates a new FaultinessReason structure
func NewFaultinessReason(processID, round uint64, reason string) *FaultinessReason {
	return &FaultinessReason{
		processID: processID,
		round:     round,
		reason:    reason,
	}
}

func (fr *FaultinessReason) String() string {
	var sb strings.Builder

	sb.WriteString("Process ")
	sb.WriteString(strconv.FormatUint(fr.processID, 10))
	sb.WriteString(" at round ")
	sb.WriteString(strconv.FormatUint(fr.round, 10))
	sb.WriteString(": ")
	sb.WriteString(fr.reason)

	return sb.String()
}

func (fr *FaultinessReason) equals(other *FaultinessReason) bool {
	return fr.processID == other.processID && fr.round == other.round && fr.reason == other.reason
}
