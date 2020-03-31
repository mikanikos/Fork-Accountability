package main

import (
	"flag"
	"fmt"
	"github.com/mikanikos/Fork-Accountability/algorithm"
	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/utils"
	"log"
)

// Monitor struct
type Monitor struct {
	Height              uint64   `yaml:"height"`
	FirstDecisionRound  uint64   `yaml:"firstDecisionRound"`
	SecondDecisionRound uint64   `yaml:"secondDecisionRound"`
	Validators          []string `yaml:"validators"`
}

const configDirectory = "/_config/"

func main() {

	// parse arguments
	configFile := flag.String("config", "", "configuration file path of the monitor")

	// parse arguments
	flag.Parse()

	// parse file
	monitor := &Monitor{}
	err := utils.ParseConfigFile(configDirectory+*configFile, monitor)
	if err != nil {
		log.Fatalf("Monitor exiting: config file not parsed correctly: %s", err)
	}

	numValidators := len(monitor.Validators)
	if numValidators == 0 {
		log.Fatal("Monitor exiting: no validators given")
	}

	// connect to validators for requesting hvs
	connectionHandler := NewConnectionHandler(numValidators)
	err = connectionHandler.connectToValidators(monitor.Validators)
	if err != nil {
		log.Fatalf("Monitor exiting: couldn't connect to all validators: %s", err)
	}

	fmt.Println("Monitor: Connected to validators, waiting for height vote sets")

	// make request for hvs to validators
	logs := common.NewHeightLogs(monitor.Height)
	connectionHandler.requestHeightLogs(monitor.Height)

	// create faulty set structure
	faultySet := algorithm.NewFaultySet()
	// lower bound on the number of faulty processes
	minFaulty := (numValidators-1)/3 + 1 // f+1

	// run until we have at least f+1 faulty processes
	for faultySet.Length() < minFaulty {
		// receive hvs from processes, it blocks the execution until another hvs arrives
		hvs := <-connectionHandler.receiveChannel
		logs.AddHvs(hvs)

		// if we have at least f+1 hvs, run the monitor algorithm
		if len(logs.Logs) >= minFaulty {
			// run monitor and get faulty processes
			faultySet = algorithm.IdentifyFaultyProcesses(uint64(numValidators), monitor.FirstDecisionRound, monitor.SecondDecisionRound, logs)
		}
	}

	fmt.Println(faultySet.String())
	fmt.Println("Monitor: Algorithm completed")
}