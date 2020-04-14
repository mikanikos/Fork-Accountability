package accountability

import (
	"fmt"
	"testing"

	"github.com/mikanikos/Fork-Accountability/utils"

	"github.com/mikanikos/Fork-Accountability/common"
)

// in some tests I assume there are more than 2f faulty processes just to see how the algorithm catches multiple faulty behaviours
// I know this is not possible for the accountability specification

func TestBasicScenario(t *testing.T) {

	numProcesses := 4
	threshold := (numProcesses-1)/3 + 1

	// create accountability struct
	acc := NewAccountability()

	acc.HeightLogs.AddHvs("1", utils.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs("2", utils.GetHvsForDefaultConfig2())
	acc.HeightLogs.AddHvs("3", utils.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs("4", utils.GetHvsForDefaultConfig4())

	if !(acc.HeightLogs.Contains("1") && acc.HeightLogs.Contains("2") && acc.HeightLogs.Contains("3") && acc.HeightLogs.Contains("4")) {
		t.Fatal("Monitor didn't store height logs correctly")
	}

	if acc.HeightLogs.Length() != numProcesses {
		t.Fatal("Monitor didn't store height logs correctly")
	}

	acc.MessageLogsReceived = uint64(numProcesses)
	if !acc.CanRun(threshold) {
		t.Fatal("Monitor should be able to run")
	}

	acc.Run(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 4, faultinessMissingQuorumForPrevote))

	if acc.FaultySet.Length() != expectedFaultySet.Length() {
		t.Fatal("Monitor detected different faulty processes")
	}

	if !acc.IsCompleted(threshold) {
		t.Fatal("Monitor should have completed")
	}

	fmt.Println(acc.String())

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithMissingHVS(t *testing.T) {

	// create accountability struct
	acc := NewAccountability()

	acc.HeightLogs.AddHvs("1", utils.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs("2", utils.GetHvsForDefaultConfig2())

	acc.Run(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 3, faultinessMissingQuorumForPrecommit))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 4, faultinessMissingQuorumForPrecommit))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 4, faultinessMissingQuorumForPrecommit))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithMoreThanOnePrecommit(t *testing.T) {

	// Process P2 - faulty
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, 20, nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "3", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "3", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, 20, nil))

	heightVoteSet2 := common.NewHeightVoteSet()
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	acc := NewAccountability()

	acc.HeightLogs.AddHvs("1", utils.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs("2", heightVoteSet2)
	acc.HeightLogs.AddHvs("3", utils.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs("4", utils.GetHvsForDefaultConfig4())

	acc.Run(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness("2", 3, faultinessMultiplePrecommits))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 4, faultinessMissingQuorumForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithMoreThanOnePrevote(t *testing.T) {

	// Process P2 - faulty
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 20, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, 10, nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "3", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "3", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, 20, nil))

	heightVoteSet2 := common.NewHeightVoteSet()
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	acc := NewAccountability()

	acc.HeightLogs.AddHvs("1", utils.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs("2", heightVoteSet2)
	acc.HeightLogs.AddHvs("3", utils.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs("4", utils.GetHvsForDefaultConfig4())

	acc.Run(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness("2", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 4, faultinessMissingQuorumForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithNotEnoughPrevoteForPrecommit(t *testing.T) {

	// Process P1 - faulty
	voteSet1 := common.NewVoteSet()
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 10, nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 10, nil))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, 10, nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 3, 10, nil))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, 20, nil))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, "1", 3, 20, nil))

	heightVoteSet1 := common.NewHeightVoteSet()
	heightVoteSet1.VoteSetMap[3] = voteSet1

	acc := NewAccountability()

	acc.HeightLogs.AddHvs("1", heightVoteSet1)
	acc.HeightLogs.AddHvs("2", utils.GetHvsForDefaultConfig2())
	acc.HeightLogs.AddHvs("3", utils.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs("4", utils.GetHvsForDefaultConfig4())

	acc.Run(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness("1", 3, faultinessMissingQuorumForPrecommit))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 4, faultinessMissingQuorumForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenario_TestNotEnoughJustifications(t *testing.T) {

	// Process P2 - faulty
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, 10, nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, 20, nil))

	heightVoteSet2 := common.NewHeightVoteSet()
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	// Process P3 - faulty
	voteSet3 := common.NewVoteSet()
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 10, nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 10, nil))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 10, nil))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, "3", 3, 10, nil))

	voteSet33 := common.NewVoteSet()
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, 20, nil))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, 20, nil))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, 20, nil))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, 20, nil))

	heightVoteSet3 := common.NewHeightVoteSet()
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	// Process P4 - faulty
	voteSet4 := common.NewVoteSet()
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 10, nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 10, nil))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 10, nil))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, "4", 3, 10, nil))

	voteSet44 := common.NewVoteSet()
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, 20, nil))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, 20, nil))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, 20, nil))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, 20, nil))

	heightVoteSet4 := common.NewHeightVoteSet()
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	acc := NewAccountability()

	acc.HeightLogs.AddHvs("1", utils.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs("2", heightVoteSet2)
	acc.HeightLogs.AddHvs("3", heightVoteSet3)
	acc.HeightLogs.AddHvs("4", heightVoteSet4)

	acc.Run(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness("2", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 4, faultinessMissingQuorumForPrevote))

	fmt.Println(expectedFaultySet)
	fmt.Println(acc.FaultySet)

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenario_TestFalseJustifications(t *testing.T) {

	// Process P2 - correct
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 10, nil))
	//voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, 10, nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "3", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, 20, nil),
		common.NewMessage(common.Prevote, "3", 3, 20, nil),
		common.NewMessage(common.Prevote, "4", 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, 20, nil))

	heightVoteSet2 := common.NewHeightVoteSet()
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	acc := NewAccountability()

	acc.HeightLogs.AddHvs("1", utils.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs("2", heightVoteSet2)
	acc.HeightLogs.AddHvs("3", utils.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs("4", utils.GetHvsForDefaultConfig4())

	acc.Run(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness("2", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("3", 4, faultinessMissingQuorumForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness("4", 4, faultinessMissingQuorumForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}
