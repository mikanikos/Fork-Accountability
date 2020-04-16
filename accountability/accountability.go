package accountability

import (
	"github.com/mikanikos/Fork-Accountability/common"
	"strings"
)

// MAIN API

// Accountability stores main information about the algorithm and offers an API for running the algorithm and getting related data
type Accountability struct {
	numValidators uint64
	heightLogs    *HeightLogs
	faultySet     *FaultySet
}

// NewAccountability creates a new Accountability structure
func NewAccountability() *Accountability {
	return &Accountability{
		heightLogs: NewHeightLogs(),
		faultySet:  NewFaultySet(),
	}
}

// String returns a string representation (result) of the accountability algorithm
func (acc *Accountability) String() string {
	var sb strings.Builder
	sb.WriteString("Accountability algorithm report generated\n")

	sb.WriteString("________________________________________________________________________________________________________________________\n\n")

	sb.WriteString(acc.heightLogs.String())
	sb.WriteString(acc.faultySet.String())

	sb.WriteString("________________________________________________________________________________________________________________________\n")


	return sb.String()
}

// Init initializes the variables needed for the execution of the accountability algorithm
func (acc *Accountability) Init(numValidators uint64) {
	acc.numValidators = numValidators

}

// IsCompleted returns true if the algorithm has completed, false otherwise
func (acc *Accountability) IsCompleted() bool {
	// if we have at least f + 1 faulty processes, the algorithm has completed
	return acc.GetNumFaulty() >= acc.getValidityThreshold()
}

// CanRun returns true if the algorithm has enough height vote sets to run, false otherwise
func (acc *Accountability) CanRun() bool {
	// if we have delivered at least f + 1 message logs, run the monitor algorithm
	return acc.GetNumLogs() >= acc.getValidityThreshold()
}

// GetNumLogs returns the number of message logs received so far
func (acc *Accountability) GetNumLogs() uint64 {
	return uint64(acc.heightLogs.ReceivedLength())
}

// GetNumFaulty returns the number of faulty processes detected in the last run of the algorithm
func (acc *Accountability) GetNumFaulty() uint64 {
	return uint64(acc.faultySet.Length())
}

// StoreHvs returns true if the hvs was added, false if it was already present
func (acc *Accountability) StoreHvs(processID string, hvs *common.HeightVoteSet) bool {
	return acc.heightLogs.AddHvs(processID, hvs)
}

func (acc *Accountability) getValidityThreshold() uint64 {
	// lower bound on the number of faulty processes and threshold for starting the algorithm
	return (acc.numValidators-1)/3 + 1 // f+1
}

func (acc *Accountability) getQuorumThreshold() uint64 {
	return acc.numValidators - (acc.numValidators-1)/3 // 2f + 1
}
