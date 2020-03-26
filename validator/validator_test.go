package main

import (
	"github.com/mikanikos/Fork-Accountability/common"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_CorrectConfigParsing(t *testing.T) {

	// Process p1 - correct
	voteSet1 := common.NewVoteSet()
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

	// Process p2 - correct
	voteSet2 := common.NewVoteSet()
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 20))
	voteSet2.ReceivedPrevoteMessages = append(voteSet2.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 20))

	voteSet2.SentPrevoteMessages = append(voteSet2.SentPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet2.SentPrecommitMessages = append(voteSet2.SentPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10))

	voteSet22 := common.NewVoteSet()
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
	voteSet3 := common.NewVoteSet()
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10))

	voteSet33 := common.NewVoteSet()
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20))

	heightVoteSet3 := common.NewHeightVoteSet(3)
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	// Process P4 - faulty
	voteSet4 := common.NewVoteSet()
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 3, 10))

	voteSet44 := common.NewVoteSet()
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20))

	heightVoteSet4 := common.NewHeightVoteSet(4)
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	fileName1 := "config_1.yaml"

	hvs1, err := parseConfigFile(fileName1)
	if err != nil {
		t.Fatalf("Failed while parsing config file 1: %s", err)
	}

	fileName2 := "config_2.yaml"

	hvs2, err := parseConfigFile(fileName2)
	if err != nil {
		t.Fatalf("Failed while parsing config file 2: %s", err)
	}

	fileName3 := "config_3.yaml"

	hvs3, err := parseConfigFile(fileName3)
	if err != nil {
		t.Fatalf("Failed while parsing config file 3: %s", err)
	}

	fileName4 := "config_4.yaml"
	hvs4, err := parseConfigFile(fileName4)
	if err != nil {
		t.Fatalf("Failed while parsing config file 4: %s", err)
	}

	eq := reflect.DeepEqual(heightVoteSet1, hvs1)
	if !eq {
		t.Fatal("HVS 1 was not parsed correctly")
	}

	eq = reflect.DeepEqual(heightVoteSet2, hvs2)
	if !eq {
		t.Fatal("HVS 2 was not parsed correctly")
	}

	eq = reflect.DeepEqual(heightVoteSet3, hvs3)
	if !eq {
		t.Fatal("HVS 3 was not parsed correctly")
	}

	eq = reflect.DeepEqual(heightVoteSet4, hvs4)
	if !eq {
		t.Fatal("HVS 4 was not parsed correctly")
	}
}

func Test_WrongConfigFilename(t *testing.T) {
	fileName3 := "config_not_existing.yaml"

	_, err := parseConfigFile(fileName3)
	if err == nil {
		t.Fatalf("Should have failed because filename doesn't exist")
	}
}

// should run in the same folder of the method called
func Test_BadFormattedConfig(t *testing.T) {

	badConfig := "bad_config.yaml"

	_ = ioutil.WriteFile(badConfig, []byte("cjdcjdcjd"), 0644)

	_, err := parseConfigFile(badConfig)
	if err == nil {
		t.Fatalf("Should have failed because file is bad formatted")
	}
	_ = os.Remove(badConfig)
}
