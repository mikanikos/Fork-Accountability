package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/utils"
)

func Test_CorrectConfigParsing_1(t *testing.T) {

	// create hvs for validator
	voteSet1 := common.NewVoteSet()
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet1.ReceivedPrevoteMessages = append(voteSet1.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))

	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10, nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 2, 3, 10, nil))
	voteSet1.ReceivedPrecommitMessages = append(voteSet1.ReceivedPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10, nil))

	voteSet1.SentPrevoteMessages = append(voteSet1.SentPrevoteMessages, common.NewMessage(common.Prevote, 1, 3, 20, nil))
	voteSet1.SentPrecommitMessages = append(voteSet1.SentPrecommitMessages, common.NewMessage(common.Precommit, 1, 3, 10, nil))

	heightVoteSet1 := common.NewHeightVoteSet(1)
	heightVoteSet1.VoteSetMap[3] = voteSet1

	// create validator structure
	validatorTest := &Validator{}
	validatorTest.ID = 1
	validatorTest.Address = "127.0.0.1:8080"
	validatorTest.Messages = make(map[uint64]*common.HeightVoteSet)
	validatorTest.Messages[1] = heightVoteSet1

	// parse config file
	validatorConfig := &Validator{}
	fileName1 := "config_1.yaml"
	err := utils.ParseConfigFile(configDirectory+fileName1, validatorConfig)
	if err != nil {
		t.Fatalf("Failed while parsing config file 1: %s", err)
	}

	// compare the two validators
	if !reflect.DeepEqual(validatorTest, validatorConfig) {
		t.Fatal("Validator 1 config file was not parsed correctly")
	}
}

func Test_CorrectConfigParsing_2(t *testing.T) {

	// create hvs for validator
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

	// create validator structure
	validatorTest := &Validator{}
	validatorTest.ID = 2
	validatorTest.Address = "127.0.0.1:8081"
	validatorTest.Messages = make(map[uint64]*common.HeightVoteSet)
	validatorTest.Messages[1] = heightVoteSet2

	// parse config file
	validatorConfig := &Validator{}
	fileName2 := "config_2.yaml"
	err := utils.ParseConfigFile(configDirectory+fileName2, validatorConfig)
	if err != nil {
		t.Fatalf("Failed while parsing config file 2: %s", err)
	}

	// compare the two validators
	if !reflect.DeepEqual(validatorTest, validatorConfig) {
		t.Fatal("Validator 2 config file was not parsed correctly")
	}
}

func Test_CorrectConfigParsing_3(t *testing.T) {

	// create hvs for validator
	voteSet3 := common.NewVoteSet()
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet3.ReceivedPrevoteMessages = append(voteSet3.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))

	voteSet3.SentPrevoteMessages = append(voteSet3.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet3.SentPrecommitMessages = append(voteSet3.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 3, 10, nil))

	voteSet33 := common.NewVoteSet()
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 3, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20, nil))
	voteSet33.ReceivedPrevoteMessages = append(voteSet33.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20, nil))

	voteSet33.SentPrevoteMessages = append(voteSet33.SentPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20, nil))
	voteSet33.SentPrecommitMessages = append(voteSet33.SentPrecommitMessages, common.NewMessage(common.Precommit, 3, 4, 20, nil))

	heightVoteSet3 := common.NewHeightVoteSet(3)
	heightVoteSet3.VoteSetMap[3] = voteSet3
	heightVoteSet3.VoteSetMap[4] = voteSet33

	// create validator structure
	validatorTest := &Validator{}
	validatorTest.ID = 3
	validatorTest.Address = "127.0.0.1:8082"
	validatorTest.Messages = make(map[uint64]*common.HeightVoteSet)
	validatorTest.Messages[1] = heightVoteSet3

	// parse config file
	validatorConfig := &Validator{}
	fileName3 := "config_3.yaml"
	err := utils.ParseConfigFile(configDirectory+fileName3, validatorConfig)
	if err != nil {
		t.Fatalf("Failed while parsing config file 3: %s", err)
	}

	// compare the two validators
	if !reflect.DeepEqual(validatorTest, validatorConfig) {
		t.Fatal("Validator 3 config file was not parsed correctly")
	}
}

func Test_CorrectConfigParsing_4(t *testing.T) {

	// create hvs for validator
	voteSet4 := common.NewVoteSet()
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 3, 10, nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 3, 10, nil))
	voteSet4.ReceivedPrevoteMessages = append(voteSet4.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))

	voteSet4.SentPrevoteMessages = append(voteSet4.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 3, 10, nil))
	voteSet4.SentPrecommitMessages = append(voteSet4.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 3, 10, nil))

	voteSet44 := common.NewVoteSet()
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 2, 4, 20, []*common.Message{
		common.NewMessage(common.Prevote, 1, 3, 20, nil),
		common.NewMessage(common.Prevote, 3, 3, 20, nil),
		common.NewMessage(common.Prevote, 4, 3, 20, nil)}))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 3, 4, 20, nil))
	voteSet44.ReceivedPrevoteMessages = append(voteSet44.ReceivedPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20, nil))

	voteSet44.SentPrevoteMessages = append(voteSet44.SentPrevoteMessages, common.NewMessage(common.Prevote, 4, 4, 20, nil))
	voteSet44.SentPrecommitMessages = append(voteSet44.SentPrecommitMessages, common.NewMessage(common.Precommit, 4, 4, 20, nil))

	heightVoteSet4 := common.NewHeightVoteSet(4)
	heightVoteSet4.VoteSetMap[3] = voteSet4
	heightVoteSet4.VoteSetMap[4] = voteSet44

	// create validator structure
	validatorTest := &Validator{}
	validatorTest.ID = 4
	validatorTest.Address = "127.0.0.1:8083"
	validatorTest.Messages = make(map[uint64]*common.HeightVoteSet)
	validatorTest.Messages[1] = heightVoteSet4

	// parse config file
	validatorConfig := &Validator{}
	fileName4 := "config_4.yaml"
	err := utils.ParseConfigFile(configDirectory+fileName4, validatorConfig)
	if err != nil {
		t.Fatalf("Failed while parsing config file 4: %s", err)
	}

	// compare the two validators
	if !reflect.DeepEqual(validatorTest, validatorConfig) {
		t.Fatal("Validator 4 config file was not parsed correctly")
	}
}

func Test_WrongConfigFilename(t *testing.T) {
	fileName3 := "config_not_existing.yaml"
	validatorConfig := &Validator{}
	err := utils.ParseConfigFile(configDirectory+fileName3, validatorConfig)
	if err == nil {
		t.Fatalf("Should have failed because filename doesn't exist")
	}
}

// should run in the same folder of the method called
func Test_BadFormattedConfig(t *testing.T) {

	badConfig := "bad_config.yaml"
	validatorConfig := &Validator{}

	_ = ioutil.WriteFile(badConfig, []byte("cjdcjdcjd"), 0644)

	err := utils.ParseConfigFile(badConfig, validatorConfig)

	_ = os.Remove(badConfig)

	if err == nil {
		t.Fatalf("Should have failed because file is bad formatted")
	}
}
