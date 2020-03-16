package monitor

const (
	errHVSNotSent                    = "The process did not send its HeightVoteSet"
	errMultiplePrevote               = "The process sent more than one PREVOTE message in a round"
	errMultiplePrecommit             = "The process sent more than one PRECOMMIT message in a round"
	errNotEnoughPrevotesForPrecommit = "The process did not receive 2f + 1 PREVOTE messages for a sent PRECOMMIT message to be issued"
	errNotEnoughPrevotesForPrevote   = "The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued"
)
