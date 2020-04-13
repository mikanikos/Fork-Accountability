package accountability

import (
	"reflect"
	"strconv"
	"strings"
)

// FaultinessReason wraps a string for simplicity
type FaultinessReason string

// FaultinessReason returns the string representation
func (fr FaultinessReason) FaultinessReason() string { return string(fr) }

const (
	faultinessMultiplePrevotes          = FaultinessReason("the process sent more than one PREVOTE message in a round")
	faultinessMultiplePrecommits        = FaultinessReason("the process sent more than one PRECOMMIT message in a round")
	faultinessMissingQuorumForPrecommit = FaultinessReason("the process did not receive 2f + 1 PREVOTE messages for a sent PRECOMMIT message to be issued")
	faultinessMissingQuorumForPrevote   = FaultinessReason("the process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued")
)

// Faultiness stores information about faulty processes
type Faultiness struct {
	processID string
	round     uint64
	reason    FaultinessReason
}

// NewFaultiness creates a new Faultiness structure
func NewFaultiness(processID string, round uint64, reason FaultinessReason) *Faultiness {
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
	sb.WriteString(fr.processID)
	sb.WriteString(" at round ")
	sb.WriteString(strconv.FormatUint(fr.round, 10))
	sb.WriteString(": ")
	sb.WriteString(fr.reason.FaultinessReason())

	return sb.String()
}

// Equal is an equality method for faultiness
func (fr *Faultiness) Equal(other *Faultiness) bool {
	return reflect.DeepEqual(fr, other)
}
