package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mikanikos/Fork-Accountability/algorithm"
)

func main() {

	// parsing arguments
	processes := flag.String("processes", "", "comma separated list of all processes addresses in the form ip:port")
	firstDecisionRound := flag.Uint64("firstDecisionRound", 0, "first round when a decision was taken")
	secondDecisionRound := flag.Uint64("secondDecisionRound", 0, "second round when a decision was taken (when fork was detected)")
	waitTimeout := flag.Uint("waitTimeout", 5, "timeout to wait for HeightVoteSet from each process")

	flag.Parse()

	// wait for validators to start (can be added/removed whether validators are already listening on their port or not)
	time.Sleep(time.Second * time.Duration(5))

	// handler for the connection with validators
	connHandler := NewConnectionHandler()

	// connect to validators for requesting hvs
	err := connHandler.connectToValidators(*processes)
	if err != nil {
		fmt.Printf("Monitor exiting: couldn't connect to all validators: %s", err)
		os.Exit(1)
	}

	fmt.Println("Monitor: Connected to validators")

	// request hvs from all processes
	hvsMap, err := connHandler.requestHeightLogs(*waitTimeout)
	if err != nil {
		fmt.Printf("Monitor exiting: error during the request of hvs: %s", err)
		os.Exit(1)
	}

	if len(hvsMap.Logs) == 0 {
		fmt.Print("Monitor exiting: no hvs received within the timeout")
		os.Exit(1)
	}

	fmt.Println("Monitor: Got all hvs from validators")

	// run monitor and get faulty processes
	faultyProcesses := algorithm.IdentifyFaultyProcesses(uint64(len(hvsMap.Logs)), *firstDecisionRound, *secondDecisionRound, hvsMap)

	fmt.Println(faultyProcesses.String())

	fmt.Println("Monitor: Run completed")
}
