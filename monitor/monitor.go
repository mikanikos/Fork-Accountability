package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
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
		panic(fmt.Errorf("Error while connecting to validators %s: ", err))
		return
	}

	fmt.Println("Monitor: Establish connection with validators")

	// request hvs from all processes
	hvsMap, err := connHandler.requestHVSWithTimeout(*waitTimeout)
	if err != nil {
		panic(err)
	}

	fmt.Println("Monitor: Got all hvs from validators")

	// run monitor and get faulty processes
	faultyProcesses := algorithm.IdentifyFaultyProcesses(uint64(len(*processes)), *firstDecisionRound, *secondDecisionRound, hvsMap)

	printFaultyProcesses(faultyProcesses)
}

/**
	Print faulty processes
 */
func printFaultyProcesses(faultyMap map[uint64][]*algorithm.FaultinessReason) {
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
