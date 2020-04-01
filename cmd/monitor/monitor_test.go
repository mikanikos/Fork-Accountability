package main

import (
	"github.com/mikanikos/Fork-Accountability/utils"
	"log"
	"reflect"
	"testing"
)

func Test_CorrectConfigParsing(t *testing.T) {

	monitorTest := &Monitor{
		Height:              1,
		FirstDecisionRound:  3,
		SecondDecisionRound: 4,
		Validators:          []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082", "127.0.0.1:8083"},
	}

	configFile := "config.yaml"
	monitorConfig := &Monitor{}
	err := utils.ParseConfigFile(configDirectory+configFile, monitorConfig)
	if err != nil {
		log.Fatalf("Monitor exiting: config file not parsed correctly: %s", err)
	}

	// compare the two validators
	if !reflect.DeepEqual(monitorTest, monitorConfig) {
		t.Fatal("Validator 1 config file was not parsed correctly")
	}
}
