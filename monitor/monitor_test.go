package monitor

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestBasicScenario(t *testing.T) {
	// Process P1 - correct
	voteSet1 := NewVoteSet(3)
	voteSet1.receivedPrevoteMessages = append(voteSet1.receivedPrevoteMessages, NewMessage(prevote, 2, 3, &BlockPublish{Transaction: 10}))
	voteSet1.receivedPrevoteMessages = append(voteSet1.receivedPrevoteMessages, NewMessage(prevote, 3, 3, &BlockPublish{Transaction: 10}))
	voteSet1.receivedPrevoteMessages = append(voteSet1.receivedPrevoteMessages, NewMessage(prevote, 4, 3, &BlockPublish{Transaction: 10}))

	voteSet1.receivedPrecommitMessages = append(voteSet1.receivedPrecommitMessages, NewMessage(precommit, 1, 3, &BlockPublish{Transaction: 10}))
	voteSet1.receivedPrecommitMessages = append(voteSet1.receivedPrecommitMessages, NewMessage(precommit, 2, 3, &BlockPublish{Transaction: 10}))
	voteSet1.receivedPrecommitMessages = append(voteSet1.receivedPrecommitMessages, NewMessage(precommit, 3, 3, &BlockPublish{Transaction: 10}))

	voteSet1.sentPrevoteMessages = append(voteSet1.sentPrevoteMessages, NewMessage(prevote, 1, 3, &BlockPublish{Transaction: 20}))
	voteSet1.sentPrecommitMessages = append(voteSet1.sentPrecommitMessages, NewMessage(precommit, 1, 3, &BlockPublish{Transaction: 10}))

	heightVoteSet1 := NewHeightVoteSet(1)
	heightVoteSet1.voteSetMap[3] = voteSet1

	// Process P2 - correct
	voteSet2 := NewVoteSet(3)
	voteSet2.receivedPrevoteMessages = append(voteSet2.receivedPrevoteMessages, NewMessage(prevote, 2, 3, &BlockPublish{Transaction: 10}))
	voteSet2.receivedPrevoteMessages = append(voteSet2.receivedPrevoteMessages, NewMessage(prevote, 3, 3, &BlockPublish{Transaction: 10}))
	voteSet2.receivedPrevoteMessages = append(voteSet2.receivedPrevoteMessages, NewMessage(prevote, 4, 3, &BlockPublish{Transaction: 10}))

	voteSet2.receivedPrevoteMessages = append(voteSet2.receivedPrevoteMessages, NewMessage(prevote, 1, 3, &BlockPublish{Transaction: 20}))
	voteSet2.receivedPrevoteMessages = append(voteSet2.receivedPrevoteMessages, NewMessage(prevote, 3, 3, &BlockPublish{Transaction: 20}))
	voteSet2.receivedPrevoteMessages = append(voteSet2.receivedPrevoteMessages, NewMessage(prevote, 4, 3, &BlockPublish{Transaction: 20}))

	voteSet2.sentPrevoteMessages = append(voteSet2.sentPrevoteMessages, NewMessage(prevote, 2, 3, &BlockPublish{Transaction: 10}))
	voteSet2.sentPrecommitMessages = append(voteSet2.sentPrecommitMessages, NewMessage(precommit, 2, 3, &BlockPublish{Transaction: 10}))

	voteSet22 := NewVoteSet(4)
	voteSet22.receivedPrevoteMessages = append(voteSet22.receivedPrevoteMessages, NewMessage(prevote, 2, 4, &BlockPublish{Transaction: 20}))
	voteSet22.receivedPrevoteMessages = append(voteSet22.receivedPrevoteMessages, NewMessage(prevote, 3, 4, &BlockPublish{Transaction: 20}))
	voteSet22.receivedPrevoteMessages = append(voteSet22.receivedPrevoteMessages, NewMessage(prevote, 4, 4, &BlockPublish{Transaction: 20}))

	voteSet22.receivedPrecommitMessages = append(voteSet22.receivedPrecommitMessages, NewMessage(precommit, 2, 4, &BlockPublish{Transaction: 20}))
	voteSet22.receivedPrecommitMessages = append(voteSet22.receivedPrecommitMessages, NewMessage(precommit, 3, 4, &BlockPublish{Transaction: 20}))
	voteSet22.receivedPrecommitMessages = append(voteSet22.receivedPrecommitMessages, NewMessage(precommit, 4, 4, &BlockPublish{Transaction: 20}))

	voteSet22.sentPrevoteMessages = append(voteSet22.sentPrevoteMessages, NewMessage(prevote, 2, 4, &BlockPublish{Transaction: 20}))
	voteSet22.sentPrecommitMessages = append(voteSet22.sentPrecommitMessages, NewMessage(precommit, 2, 4, &BlockPublish{Transaction: 20}))

	heightVoteSet2 := NewHeightVoteSet(2)
	heightVoteSet2.voteSetMap[3] = voteSet2
	heightVoteSet2.voteSetMap[4] = voteSet22

	// Process P3 - faulty
	voteSet3 := NewVoteSet(3)
	voteSet3.receivedPrevoteMessages = append(voteSet3.receivedPrevoteMessages, NewMessage(prevote, 2, 3, &BlockPublish{Transaction: 10}))
	voteSet3.receivedPrevoteMessages = append(voteSet3.receivedPrevoteMessages, NewMessage(prevote, 3, 3, &BlockPublish{Transaction: 10}))
	voteSet3.receivedPrevoteMessages = append(voteSet3.receivedPrevoteMessages, NewMessage(prevote, 4, 3, &BlockPublish{Transaction: 10}))

	voteSet3.sentPrevoteMessages = append(voteSet3.sentPrevoteMessages, NewMessage(prevote, 3, 3, &BlockPublish{Transaction: 10}))
	voteSet3.sentPrecommitMessages = append(voteSet3.sentPrecommitMessages, NewMessage(precommit, 3, 3, &BlockPublish{Transaction: 10}))

	voteSet33 := NewVoteSet(4)
	voteSet33.receivedPrevoteMessages = append(voteSet33.receivedPrevoteMessages, NewMessage(prevote, 2, 4, &BlockPublish{Transaction: 20}))
	voteSet33.receivedPrevoteMessages = append(voteSet33.receivedPrevoteMessages, NewMessage(prevote, 3, 4, &BlockPublish{Transaction: 20}))
	voteSet33.receivedPrevoteMessages = append(voteSet33.receivedPrevoteMessages, NewMessage(prevote, 4, 4, &BlockPublish{Transaction: 20}))

	voteSet33.sentPrevoteMessages = append(voteSet33.sentPrevoteMessages, NewMessage(prevote, 3, 4, &BlockPublish{Transaction: 20}))
	voteSet33.sentPrecommitMessages = append(voteSet33.sentPrecommitMessages, NewMessage(precommit, 3, 4, &BlockPublish{Transaction: 20}))

	heightVoteSet3 := NewHeightVoteSet(3)
	heightVoteSet3.voteSetMap[3] = voteSet3
	heightVoteSet3.voteSetMap[4] = voteSet33

	// Process P4 - faulty
	voteSet4 := NewVoteSet(3)
	voteSet4.receivedPrevoteMessages = append(voteSet4.receivedPrevoteMessages, NewMessage(prevote, 2, 3, &BlockPublish{Transaction: 10}))
	voteSet4.receivedPrevoteMessages = append(voteSet4.receivedPrevoteMessages, NewMessage(prevote, 3, 3, &BlockPublish{Transaction: 10}))
	voteSet4.receivedPrevoteMessages = append(voteSet4.receivedPrevoteMessages, NewMessage(prevote, 4, 3, &BlockPublish{Transaction: 10}))

	voteSet4.sentPrevoteMessages = append(voteSet4.sentPrevoteMessages, NewMessage(prevote, 4, 3, &BlockPublish{Transaction: 10}))
	voteSet4.sentPrecommitMessages = append(voteSet4.sentPrecommitMessages, NewMessage(precommit, 4, 3, &BlockPublish{Transaction: 10}))

	voteSet44 := NewVoteSet(4)
	voteSet44.receivedPrevoteMessages = append(voteSet44.receivedPrevoteMessages, NewMessage(prevote, 2, 4, &BlockPublish{Transaction: 20}))
	voteSet44.receivedPrevoteMessages = append(voteSet44.receivedPrevoteMessages, NewMessage(prevote, 3, 4, &BlockPublish{Transaction: 20}))
	voteSet44.receivedPrevoteMessages = append(voteSet44.receivedPrevoteMessages, NewMessage(prevote, 4, 4, &BlockPublish{Transaction: 20}))

	voteSet44.sentPrevoteMessages = append(voteSet44.sentPrevoteMessages, NewMessage(prevote, 4, 4, &BlockPublish{Transaction: 20}))
	voteSet44.sentPrecommitMessages = append(voteSet44.sentPrecommitMessages, NewMessage(precommit, 4, 4, &BlockPublish{Transaction: 20}))

	heightVoteSet4 := NewHeightVoteSet(4)
	heightVoteSet4.voteSetMap[3] = voteSet4
	heightVoteSet4.voteSetMap[4] = voteSet44

	heightVoteSets := []*HeightVoteSet{heightVoteSet1, heightVoteSet2, heightVoteSet3, heightVoteSet4}
	faultyMap := IdentifyFaultyProcesses(4, 3, 4, heightVoteSets)

	var sb strings.Builder

	sb.WriteString("Faulty processes are: \n")

	for processID, reasonsList := range faultyMap {
		sb.WriteString(strconv.FormatUint(processID, 10))
		sb.WriteString(": ")

		for _, reason := range reasonsList {
			sb.WriteString(reason.String())
			sb.WriteString("; ")
		}

		sb.WriteString("\n")
	}

	fmt.Println(sb.String())
}
