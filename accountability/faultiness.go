package accountability

// Faultiness wraps a string for simplicity
type Faultiness string

// Faultiness returns the string representation
func (fr Faultiness) FaultinessReason() string { return string(fr) }

const (
	faultinessMissingHvs                      = Faultiness("The process did not send its HeightVoteSet")
	faultinessMultiplePrevotes                = Faultiness("The process sent more than one PREVOTE message in a round")
	faultinessMultiplePrecommits              = Faultiness("The process sent more than one PRECOMMIT message in a round")
	faultinessMissingQuorumForPrecommit       = Faultiness("The process did not receive 2f + 1 PREVOTE messages for a sent PRECOMMIT message to be issued")
	faultinessMissingQuorumForPrevote         = Faultiness("The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued")
	faultinessMissingJustificationsForPrevote = Faultiness("The process had sent PRECOMMIT message, and did not have enough justifications (2f + 1 PREVOTE messages) in the sent PREVOTE message for another value to be issued")
)
