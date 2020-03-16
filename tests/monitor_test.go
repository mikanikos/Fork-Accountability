package monitor

import (
	"strconv"
	"strings"
	"testing"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/monitor"
)

func TestBasicScenario(t *testing.T) {

	// Process P1 - correct
	voteSet1 := common.NewVoteSet(3)
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10))

	heightVoteSet1 := common.NewHeightVoteSet(1)
	heightVoteSet1.VoteSetMap[3] = voteSet1

	// Process P2 - correct
	voteSet2 := common.NewVoteSet(3)
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))

	voteSet22 := common.NewVoteSet(4)
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	// Process P3 - faulty
	voteSet3 := common.NewVoteSet(3)
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet33 := common.NewVoteSet(4)
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))

	heightVoteSet3 := common.NewHeightVoteSet(3)
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	// Process P4 - faulty
	voteSet4 := common.NewVoteSet(3)
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 3, 10))

	voteSet44 := common.NewVoteSet(4)
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	heightVoteSet4 := common.NewHeightVoteSet(4)
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	heightVoteSets := map[uint64]*common.HeightVoteSet{1 :heightVoteSet1, 2: heightVoteSet2, 3: heightVoteSet3, 4: heightVoteSet4}
	faultyMap := monitor.IdentifyFaultyProcesses(4, 3, 4, heightVoteSets)

	var sb strings.Builder

	for processID, reasonsList := range faultyMap {
		sb.WriteString(strconv.FormatUint(processID, 10))
		sb.WriteString(": ")

		for _, reason := range reasonsList {
			sb.WriteString(reason.String())
			sb.WriteString("; ")
		}

		sb.WriteString("\n")
	}

	expectedFaultyVal3 := "3: Process 3 at round 3: The process sent more than one PREVOTE message in a round; Process 3 at round 4: The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued;"
	expectedFaultyVal4 := "4: Process 4 at round 3: The process sent more than one PREVOTE message in a round; Process 4 at round 4: The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued;"

	result := sb.String()
	if !strings.Contains(result, expectedFaultyVal3) || !strings.Contains(result, expectedFaultyVal4) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithMissingHVS(t *testing.T) {

	// Process P1 - correct
	voteSet1 := common.NewVoteSet(3)
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10))

	heightVoteSet1 := common.NewHeightVoteSet(1)
	heightVoteSet1.VoteSetMap[3] = voteSet1

	// Process P2 - correct
	voteSet2 := common.NewVoteSet(3)
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))

	voteSet22 := common.NewVoteSet(4)
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	heightVoteSets := map[uint64]*common.HeightVoteSet{1 :heightVoteSet1, 2: heightVoteSet2, 3: nil, 4: nil}
	faultyMap := monitor.IdentifyFaultyProcesses(4, 3, 4, heightVoteSets)

	var sb strings.Builder

	for processID, reasonsList := range faultyMap {
		sb.WriteString(strconv.FormatUint(processID, 10))
		sb.WriteString(": ")

		for _, reason := range reasonsList {
			sb.WriteString(reason.String())
			sb.WriteString("; ")
		}

		sb.WriteString("\n")
	}

	expectedFaultyVal3 := "3: Process 3 at round 0: The process did not send its HeightVoteSet"
	expectedFaultyVal4 := "4: Process 4 at round 0: The process did not send its HeightVoteSet"

	result := sb.String()
	if !strings.Contains(result, expectedFaultyVal3) || !strings.Contains(result, expectedFaultyVal4) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithMoreThanOnePrecommit(t *testing.T) {

	// Process P1 - correct
	voteSet1 := common.NewVoteSet(3)
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10))

	heightVoteSet1 := common.NewHeightVoteSet(1)
	heightVoteSet1.VoteSetMap[3] = voteSet1

	// Process P2 - faulty
	voteSet2 := common.NewVoteSet(3)
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 20))

	voteSet22 := common.NewVoteSet(4)
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	// Process P3 - faulty
	voteSet3 := common.NewVoteSet(3)
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet33 := common.NewVoteSet(4)
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))

	heightVoteSet3 := common.NewHeightVoteSet(3)
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	// Process P4 - faulty
	voteSet4 := common.NewVoteSet(3)
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 3, 10))

	voteSet44 := common.NewVoteSet(4)
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	heightVoteSet4 := common.NewHeightVoteSet(4)
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	heightVoteSets := map[uint64]*common.HeightVoteSet{1 :heightVoteSet1, 2: heightVoteSet2, 3: heightVoteSet3, 4: heightVoteSet4}
	faultyMap := monitor.IdentifyFaultyProcesses(4, 3, 4, heightVoteSets)

	var sb strings.Builder

	for processID, reasonsList := range faultyMap {
		sb.WriteString(strconv.FormatUint(processID, 10))
		sb.WriteString(": ")

		for _, reason := range reasonsList {
			sb.WriteString(reason.String())
			sb.WriteString("; ")
		}

		sb.WriteString("\n")
	}

	expectedFaultyVal2 := "2: Process 2 at round 3: The process sent more than one PRECOMMIT message in a round"
	expectedFaultyVal3 := "3: Process 3 at round 3: The process sent more than one PREVOTE message in a round; Process 3 at round 4: The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued;"
	expectedFaultyVal4 := "4: Process 4 at round 3: The process sent more than one PREVOTE message in a round; Process 4 at round 4: The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued;"

	result := sb.String()
	if !strings.Contains(result, expectedFaultyVal2) || !strings.Contains(result, expectedFaultyVal3) || !strings.Contains(result, expectedFaultyVal4) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithMoreThanOnePrevote(t *testing.T) {

	// Process P1 - correct
	voteSet1 := common.NewVoteSet(3)
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10))

	heightVoteSet1 := common.NewHeightVoteSet(1)
	heightVoteSet1.VoteSetMap[3] = voteSet1

	// Process P2 - faulty
	voteSet2 := common.NewVoteSet(3)
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 20))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))

	voteSet22 := common.NewVoteSet(4)
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	// Process P3 - faulty
	voteSet3 := common.NewVoteSet(3)
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet33 := common.NewVoteSet(4)
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))

	heightVoteSet3 := common.NewHeightVoteSet(3)
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	// Process P4 - faulty
	voteSet4 := common.NewVoteSet(3)
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 3, 10))

	voteSet44 := common.NewVoteSet(4)
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	heightVoteSet4 := common.NewHeightVoteSet(4)
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	heightVoteSets := map[uint64]*common.HeightVoteSet{1 :heightVoteSet1, 2: heightVoteSet2, 3: heightVoteSet3, 4: heightVoteSet4}
	faultyMap := monitor.IdentifyFaultyProcesses(4, 3, 4, heightVoteSets)

	var sb strings.Builder

	for processID, reasonsList := range faultyMap {
		sb.WriteString(strconv.FormatUint(processID, 10))
		sb.WriteString(": ")

		for _, reason := range reasonsList {
			sb.WriteString(reason.String())
			sb.WriteString("; ")
		}

		sb.WriteString("\n")
	}

	expectedFaultyVal2 := "2: Process 2 at round 3: The process sent more than one PREVOTE message in a round"
	expectedFaultyVal3 := "3: Process 3 at round 3: The process sent more than one PREVOTE message in a round; Process 3 at round 4: The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued;"
	expectedFaultyVal4 := "4: Process 4 at round 3: The process sent more than one PREVOTE message in a round; Process 4 at round 4: The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued;"

	result := sb.String()
	if !strings.Contains(result, expectedFaultyVal2) || !strings.Contains(result, expectedFaultyVal3) || !strings.Contains(result, expectedFaultyVal4) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithNotEnoughPrevoteForPrecommit(t *testing.T) {

	// Process P1 - faulty
	voteSet1 := common.NewVoteSet(3)
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 20))

	heightVoteSet1 := common.NewHeightVoteSet(1)
	heightVoteSet1.VoteSetMap[3] = voteSet1

	// Process P2 - correct
	voteSet2 := common.NewVoteSet(3)
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))

	voteSet22 := common.NewVoteSet(4)
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	// Process P3 - faulty
	voteSet3 := common.NewVoteSet(3)
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet33 := common.NewVoteSet(4)
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))

	heightVoteSet3 := common.NewHeightVoteSet(3)
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	// Process P4 - faulty
	voteSet4 := common.NewVoteSet(3)
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 3, 10))

	voteSet44 := common.NewVoteSet(4)
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	heightVoteSet4 := common.NewHeightVoteSet(4)
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	heightVoteSets := map[uint64]*common.HeightVoteSet{1 :heightVoteSet1, 2: heightVoteSet2, 3: heightVoteSet3, 4: heightVoteSet4}
	faultyMap := monitor.IdentifyFaultyProcesses(4, 3, 4, heightVoteSets)

	var sb strings.Builder

	for processID, reasonsList := range faultyMap {
		sb.WriteString(strconv.FormatUint(processID, 10))
		sb.WriteString(": ")

		for _, reason := range reasonsList {
			sb.WriteString(reason.String())
			sb.WriteString("; ")
		}

		sb.WriteString("\n")
	}

	expectedFaultyVal1 := "1: Process 1 at round 3: The process did not receive 2f + 1 PREVOTE messages for a sent PRECOMMIT message to be issued"
	expectedFaultyVal3 := "3: Process 3 at round 3: The process sent more than one PREVOTE message in a round; Process 3 at round 4: The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued;"
	expectedFaultyVal4 := "4: Process 4 at round 3: The process sent more than one PREVOTE message in a round; Process 4 at round 4: The process had sent PRECOMMIT message, and did not receive 2f + 1 PREVOTE messages for a sent PREVOTE message for another value to be issued;"

	result := sb.String()
	if !strings.Contains(result, expectedFaultyVal1) || !strings.Contains(result, expectedFaultyVal3) || !strings.Contains(result, expectedFaultyVal4) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}