package monitor

import (
	"fmt"
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
