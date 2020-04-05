package common

func GetHvsForDefaultConfig1() *HeightVoteSet {
	// Process P1 - correct
	voteSet1 := NewVoteSet()
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, NewMessage(Prevote, 2, 3, 10, nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, NewMessage(Prevote, 3, 3, 10, nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, NewMessage(Prevote, 4, 3, 10, nil))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, NewMessage(Precommit, 1, 3, 10, nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, NewMessage(Precommit, 2, 3, 10, nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, NewMessage(Precommit, 3, 3, 10, nil))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, NewMessage(Prevote, 1, 3, 20, nil))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, NewMessage(Precommit, 1, 3, 10, nil))

	heightVoteSet1 := NewHeightVoteSet(1)
	heightVoteSet1.VoteSetMap[3] = voteSet1

	return heightVoteSet1
}

func GetHvsForDefaultConfig2() *HeightVoteSet {
	// Process P2 - correct
	voteSet2 := NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, NewMessage(Prevote, 2, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, NewMessage(Prevote, 3, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, NewMessage(Prevote, 4, 3, 10, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, NewMessage(Prevote, 1, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, NewMessage(Prevote, 3, 3, 20, nil))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, NewMessage(Prevote, 4, 3, 20, nil))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, NewMessage(Prevote, 2, 3, 10, nil))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, NewMessage(Precommit, 2, 3, 10, nil))

	voteSet22 := NewVoteSet()
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, NewMessage(Prevote, 2, 4, 20, []*Message{
		NewMessage(Prevote, 1, 3, 20, nil),
		NewMessage(Prevote, 3, 3, 20, nil),
		NewMessage(Prevote, 4, 3, 20, nil)}))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, NewMessage(Prevote, 3, 4, 20, nil))
	voteSet22.ReceivedPrevoteMessages = append(voteSet22.ReceivedPrevoteMessages, NewMessage(Prevote, 4, 4, 20, nil))

	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, NewMessage(Precommit, 2, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, NewMessage(Precommit, 3, 4, 20, nil))
	voteSet22.ReceivedPrecommitMessages = append(voteSet22.ReceivedPrecommitMessages, NewMessage(Precommit, 4, 4, 20, nil))

	voteSet22.SentPrevoteMessages = append(voteSet22.SentPrevoteMessages, NewMessage(Prevote, 2, 4, 20, []*Message{
		NewMessage(Prevote, 1, 3, 20, nil),
		NewMessage(Prevote, 3, 3, 20, nil),
		NewMessage(Prevote, 4, 3, 20, nil)}))
	voteSet22.SentPrecommitMessages = append(voteSet22.SentPrecommitMessages, NewMessage(Precommit, 2, 4, 20, nil))

	heightVoteSet2 := NewHeightVoteSet(2)
	heightVoteSet2.VoteSetMap[3] = voteSet2
	heightVoteSet2.VoteSetMap[4] = voteSet22

	return heightVoteSet2
}

func GetHvsForDefaultConfig3() *HeightVoteSet {
	// Process P3 - faulty
	voteSet3 := NewVoteSet()
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, NewMessage(Prevote, 2, 3, 10, nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, NewMessage(Prevote, 3, 3, 10, nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, NewMessage(Prevote, 4, 3, 10, nil))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, NewMessage(Prevote, 3, 3, 10, nil))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, NewMessage(Precommit, 3, 3, 10, nil))

	voteSet33 := NewVoteSet()
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, NewMessage(Prevote, 2, 4, 20, []*Message{
		NewMessage(Prevote, 1, 3, 20, nil),
		NewMessage(Prevote, 3, 3, 20, nil),
		NewMessage(Prevote, 4, 3, 20, nil)}))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, NewMessage(Prevote, 3, 4, 20, nil))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, NewMessage(Prevote, 4, 4, 20, nil))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, NewMessage(Prevote, 3, 4, 20, nil))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, NewMessage(Precommit, 3, 4, 20, nil))

	heightVoteSet3 := NewHeightVoteSet(3)
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	return heightVoteSet3
}

func GetHvsForDefaultConfig4() *HeightVoteSet {
	// Process P4 - faulty
	voteSet4 := NewVoteSet()
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, NewMessage(Prevote, 2, 3, 10, nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, NewMessage(Prevote, 3, 3, 10, nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, NewMessage(Prevote, 4, 3, 10, nil))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, NewMessage(Prevote, 4, 3, 10, nil))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, NewMessage(Precommit, 4, 3, 10, nil))

	voteSet44 := NewVoteSet()
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, NewMessage(Prevote, 2, 4, 20, []*Message{
		NewMessage(Prevote, 1, 3, 20, nil),
		NewMessage(Prevote, 3, 3, 20, nil),
		NewMessage(Prevote, 4, 3, 20, nil)}))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, NewMessage(Prevote, 3, 4, 20, nil))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, NewMessage(Prevote, 4, 4, 20, nil))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, NewMessage(Prevote, 4, 4, 20, nil))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, NewMessage(Precommit, 4, 4, 20, nil))

	heightVoteSet4 := NewHeightVoteSet(4)
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	return heightVoteSet4
}
