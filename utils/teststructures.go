package utils

import "github.com/mikanikos/Fork-Accountability/common"

// GetHvsForDefaultConfig1 for hvs of the validator with config_1.yaml
func GetHvsForDefaultConfig1() *common.HeightVoteSet {
	// Process P1 - correct
	voteSet1 := common.NewVoteSet()
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "1", 3, common.NewValue(10), nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, common.NewValue(10), nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 3, common.NewValue(10), nil))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, common.NewValue(20), nil))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, "1", 3, common.NewValue(10), nil))

	heightVoteSet1 := common.NewHeightVoteSet()
	heightVoteSet1.VoteSetMap[3] = voteSet1

	return heightVoteSet1
}

// GetHvsForDefaultConfig2 for hvs of the validator with config_2.yaml
func GetHvsForDefaultConfig2() *common.HeightVoteSet {
	// Process P2 - correct
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, common.NewValue(20), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(20), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(20), nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, common.NewValue(10), nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, common.NewValue(20), []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, common.NewValue(20), nil),
		common.NewMessage(common.Prevote, "3", 3, common.NewValue(20), nil),
		common.NewMessage(common.Prevote, "4", 3, common.NewValue(20), nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, common.NewValue(20), nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, common.NewValue(20), nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, common.NewValue(20), nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, common.NewValue(20), nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, common.NewValue(20), nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, common.NewValue(20), []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, common.NewValue(20), nil),
		common.NewMessage(common.Prevote, "3", 3, common.NewValue(20), nil),
		common.NewMessage(common.Prevote, "4", 3, common.NewValue(20), nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, common.NewValue(20), nil))

	heightVoteSet2 := common.NewHeightVoteSet()
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	return heightVoteSet2
}

// GetHvsForDefaultConfig3 for hvs of the validator with config_3.yaml
func GetHvsForDefaultConfig3() *common.HeightVoteSet {
	// Process P3 - faulty
	voteSet3 := common.NewVoteSet()
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, "3", 3, common.NewValue(10), nil))

	voteSet33 := common.NewVoteSet()
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, common.NewValue(20), []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, common.NewValue(20), nil),
		common.NewMessage(common.Prevote, "3", 3, common.NewValue(20), nil),
		common.NewMessage(common.Prevote, "4", 3, common.NewValue(20), nil)}))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, common.NewValue(20), nil))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, common.NewValue(20), nil))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, common.NewValue(20), nil))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, common.NewValue(20), nil))

	heightVoteSet3 := common.NewHeightVoteSet()
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	return heightVoteSet3
}

// GetHvsForDefaultConfig4 for hvs of the validator with config_4.yaml
func GetHvsForDefaultConfig4() *common.HeightVoteSet {
	// Process P4 - faulty
	voteSet4 := common.NewVoteSet()
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, "4", 3, common.NewValue(10), nil))

	voteSet44 := common.NewVoteSet()
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, common.NewValue(20), []*common.Message{
		common.NewMessage(common.Prevote, "1", 3, common.NewValue(20), nil),
		common.NewMessage(common.Prevote, "3", 3, common.NewValue(20), nil),
		common.NewMessage(common.Prevote, "4", 3, common.NewValue(20), nil)}))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, common.NewValue(20), nil))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, common.NewValue(20), nil))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, common.NewValue(20), nil))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, common.NewValue(20), nil))

	heightVoteSet4 := common.NewHeightVoteSet()
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	return heightVoteSet4
}

// GetHvsForDefaultConfig1WithNoJustifications for hvs of the validator with config_1.yaml
func GetHvsForDefaultConfig1WithNoJustifications() *common.HeightVoteSet {
	// Process P1 - correct
	voteSet1 := common.NewVoteSet()
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "1", 3, common.NewValue(10), nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, common.NewValue(10), nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 3, common.NewValue(10), nil))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, common.NewValue(20), nil))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, "1", 3, common.NewValue(10), nil))

	heightVoteSet1 := common.NewHeightVoteSet()
	heightVoteSet1.VoteSetMap[3] = voteSet1

	return heightVoteSet1
}

// GetHvsForDefaultConfig2WithNoJustifications for hvs of the validator with config_2.yaml
func GetHvsForDefaultConfig2WithNoJustifications() *common.HeightVoteSet {
	// Process P2 - correct
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "1", 3, common.NewValue(20), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(20), nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(20), nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 3, common.NewValue(10), nil))

	voteSet22 := common.NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, common.NewValue(20), nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, common.NewValue(20), nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, common.NewValue(20), nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, common.NewValue(20), nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, common.NewValue(20), nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, common.NewValue(20), nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, common.NewValue(20), nil))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, common.NewMessage(common.Precommit, "2", 4, common.NewValue(20), nil))

	heightVoteSet2 := common.NewHeightVoteSet()
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	return heightVoteSet2
}

// GetHvsForDefaultConfig3WithNoJustifications for hvs of the validator with config_3.yaml
func GetHvsForDefaultConfig3WithNoJustifications() *common.HeightVoteSet {
	// Process P3 - faulty
	voteSet3 := common.NewVoteSet()
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, "3", 3, common.NewValue(10), nil))

	voteSet33 := common.NewVoteSet()
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, common.NewValue(20), nil))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, common.NewValue(20), nil))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, common.NewValue(20), nil))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, common.NewValue(20), nil))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, "3", 4, common.NewValue(20), nil))

	heightVoteSet3 := common.NewHeightVoteSet()
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	return heightVoteSet3
}

// GetHvsForDefaultConfig4WithNoJustifications for hvs of the validator with config_4.yaml
func GetHvsForDefaultConfig4WithNoJustifications() *common.HeightVoteSet {
	// Process P4 - faulty
	voteSet4 := common.NewVoteSet()
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 3, common.NewValue(10), nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 3, common.NewValue(10), nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, "4", 3, common.NewValue(10), nil))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, "4", 3, common.NewValue(10), nil))

	voteSet44 := common.NewVoteSet()
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "2", 4, common.NewValue(20), nil))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "3", 4, common.NewValue(20), nil))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, common.NewValue(20), nil))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, "4", 4, common.NewValue(20), nil))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, "4", 4, common.NewValue(20), nil))

	heightVoteSet4 := common.NewHeightVoteSet()
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	return heightVoteSet4
}
