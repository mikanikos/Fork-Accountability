package accountability

import (
	"testing"

	"github.com/mikanikos/Fork-Accountability/common"
)

// in some tests I assume there are more than 2f faulty processes just to see how the algorithm catches multiple faulty behaviours
// I know this is not possible for the accountability specification

func TestBasicScenario(t *testing.T) {

	// create accountability struct
	acc := NewAccountability()

	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig2())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig4())

	acc.IdentifyFaultyProcesses(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 4, faultinessNotEnoughPrevotesForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 4, faultinessNotEnoughPrevotesForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithMoreThanOnePrecommit(t *testing.T) {

	// Process P2 - faulty
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 20, nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 3, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 3, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20, nil))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	acc := NewAccountability()

	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs(heightVoteSet2)
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig4())

	acc.IdentifyFaultyProcesses(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness(2, 3, faultinessMultiplePrecommits))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 4, faultinessNotEnoughPrevotesForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 4, faultinessNotEnoughPrevotesForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithMoreThanOnePrevote(t *testing.T) {

	// Process P2 - faulty
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 20, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10, nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 3, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 3, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20, nil))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	acc := NewAccountability()

	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs(heightVoteSet2)
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig4())

	acc.IdentifyFaultyProcesses(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness(2, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 4, faultinessNotEnoughPrevotesForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 4, faultinessNotEnoughPrevotesForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenarioWithNotEnoughPrevoteForPrecommit(t *testing.T) {

	// Process P1 - faulty
	voteSet1 := common.NewVoteSet()
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10, nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10, nil))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20, nil))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 20, nil))

	heightVoteSet1 := common.NewHeightVoteSet(1)
	heightVoteSet1.VoteSetMap[3] = voteSet1

	acc := NewAccountability()

	acc.HeightLogs.AddHvs(heightVoteSet1)
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig2())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig4())

	acc.IdentifyFaultyProcesses(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness(1, 3, faultinessNotEnoughPrevotesForPrecommit))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 4, faultinessNotEnoughPrevotesForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 4, faultinessNotEnoughPrevotesForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenario_TestNotEnoughJustifications(t *testing.T) {

	// Process P2 - faulty
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10, nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20, nil))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	acc := NewAccountability()

	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs(heightVoteSet2)
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig4())

	acc.IdentifyFaultyProcesses(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness(2, 4, faultinessNotEnoughPrevotesForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 4, faultinessNotEnoughPrevotesForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 4, faultinessNotEnoughPrevotesForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}

func TestBasicScenario_TestFalseJustifications(t *testing.T) {

	// Process P2 - correct
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))
	//voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10, nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 3, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 3, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 4, 20, nil))

	heightVoteSet2 := common.NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	acc := NewAccountability()

	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig1())
	acc.HeightLogs.AddHvs(heightVoteSet2)
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig3())
	acc.HeightLogs.AddHvs(common.GetHvsForDefaultConfig4())

	acc.IdentifyFaultyProcesses(4, 3, 4)

	expectedFaultySet := NewFaultySet()
	expectedFaultySet.AddFaultinessReason(NewFaultiness(2, 4, faultinessNotEnoughPrevotesForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(3, 4, faultinessNotEnoughPrevotesForPrevote))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 3, faultinessMultiplePrevotes))
	expectedFaultySet.AddFaultinessReason(NewFaultiness(4, 4, faultinessNotEnoughPrevotesForPrevote))

	if !acc.FaultySet.Equal(expectedFaultySet) {
		t.Fatal("Monitor failed to detect faulty processes")
	}
}
