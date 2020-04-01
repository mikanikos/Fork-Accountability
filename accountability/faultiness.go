package accountability

import (
	"reflect"
	"strconv"
	"strings"
)

type FaultinessReason string

func (fr FaultinessReason) FaultinessReason() string { return string(fr) }

const (
	FaultinessHVSNotSent                    = FaultinessReason("the process did not send its HeightVoteSet")
	FaultinessMultiplePrevotes              = FaultinessReason("the process sent more than one PREVOTE message in a round")
	FaultinessMultiplePrecommits            = FaultinessReason("the process sent more than one PRECOMMIT message in a round")
	FaultinessNotEnoughPrevotesForPrecommit = FaultinessReason("the process did not receive 2f + 1 PREVOTE messages for a sent PRECOMMIT message to be issued")
	FaultinessNotEnoughPrevotesForPrevote   = FaultinessReason("the process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued")
)

// Faultiness stores information about faulty processes
type Faultiness struct {
	processID uint64
	round     uint64
	reason    FaultinessReason
}

// NewAccountability creates a new Faultiness structure
func NewFaultiness(processID, round uint64, reason FaultinessReason) *Faultiness {
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
	sb.WriteString(fr.reason.FaultinessReason())

	return sb.String()
}

// equality for faultiness
func (fr *Faultiness) Equal(other *Faultiness) bool {
	return reflect.DeepEqual(fr, other)
}
