package main

import (
	"fmt"
	"github.com/mikanikos/Fork-Accountability/accountability"
	"github.com/mikanikos/Fork-Accountability/common"
	"log"
)

// Monitor struct
type Monitor struct {
	Height              uint64   `yaml:"height"`
	FirstDecisionRound  uint64   `yaml:"firstDecisionRound"`
	SecondDecisionRound uint64   `yaml:"secondDecisionRound"`
	Validators          []string `yaml:"validators"`
}

func (monitor *Monitor) Run() {
	numValidators := len(monitor.Validators)
	if numValidators == 0 {
		log.Fatal("Monitor exiting: no validators given")
	}

	// connect to validators for requesting hvs
	connectionHandler := NewConnectionHandler(numValidators)
	err := connectionHandler.connectToValidators(monitor.Validators)
	if err != nil {
		log.Fatalf("Monitor exiting: couldn't connect to all validators: %s", err)
	}

	fmt.Println("Monitor: Connected to validators, waiting for height vote sets")

	// make request for hvs to validators
	connectionHandler.requestHeightLogs(monitor.Height)

	// run accountability algorithm
	monitor.RunAccountabilityAlgorithm(connectionHandler)
}

func (monitor *Monitor) RunAccountabilityAlgorithm(connectionHandler *ConnectionHandler) {
	numValidators := len(monitor.Validators)
	logs := common.NewHeightLogs(monitor.Height)

	// create faulty set structure
	acc := accountability.NewAccountability()
	// lower bound on the number of faulty processes
	minFaulty := (numValidators-1)/3 + 1 // f+1

	// run until we have at least f+1 faulty processes
	for acc.Length() < minFaulty {
		// receive hvs from processes, it blocks the execution until another hvs arrives
		hvs := <-connectionHandler.receiveChannel
		logs.AddHvs(hvs)

		// if we have at least f+1 hvs, run the monitor algorithm
		if len(logs.Logs) >= minFaulty {
			// run monitor and get faulty processes
			acc.IdentifyFaultyProcesses(uint64(numValidators), monitor.FirstDecisionRound, monitor.SecondDecisionRound, logs)
		}
	}

	fmt.Println(acc.String())
	fmt.Println("Monitor: Algorithm completed")
}
