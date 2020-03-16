package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mikanikos/Fork-Accountability/connection"

	"github.com/mikanikos/Fork-Accountability/monitor"

	"github.com/mikanikos/Fork-Accountability/common"
)

func main() {

	// parsing arguments
	processes := flag.String("processes", "", "comma separated list of all processes addresses in the form ip:port")
	firstDecisionRound := flag.Uint64("firstDecisionRound", 0, "first round when a decision was taken")
	secondDecisionRound := flag.Uint64("secondDecisionRound", 0, "second round when a decision was taken (when fork was detected)")
	waitTimeout := flag.Uint("waitTimeout", 5, "timeout to wait for HeightVoteSet from each process")

	flag.Parse()

	// wait for validators to start (can be added/removed whether validators are already listening on their port or not)
	//time.Sleep(time.Second * time.Duration(5))

	processesConn, err := establishConnections(*processes)

	if err != nil {
		panic(err)
	}

	fmt.Println("Monitor: Establish connection with validators")

	// request hvs from all processes
	hvsMap, err := requestHVSWithTimeout(processesConn, *waitTimeout)

	if err != nil {
		panic(err)
	}

	fmt.Println("Monitor: Got all hvs from validators")

	// run monitor and get faulty processes
	faultyProcesses := monitor.IdentifyFaultyProcesses(uint64(len(processesConn)), *firstDecisionRound, *secondDecisionRound, hvsMap)

	printFaultyProcesses(faultyProcesses)
}

func printFaultyProcesses(faultyMap map[uint64][]*monitor.FaultinessReason) {
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

// resolve processes addresses
func establishConnections(validators string) ([]net.Conn, error) {

	// split list of string addresses only if it's not empty in order to avoid problems
	validatorsList := make([]string, 0)

	if validators != "" {
		validatorsList = strings.Split(validators, ",")
	}

	// resolve peers addresses given
	validatorsConn := make([]net.Conn, 0)
	for _, val := range validatorsList {
		conn, err := net.Dial("tcp", val)
		if err == nil {
			validatorsConn = append(validatorsConn, conn)
		} else {
			return nil, fmt.Errorf("error while connecting to one of the validators given: %s", err)
		}
	}

	return validatorsConn, nil
}

func requestHVSWithTimeout(connections []net.Conn, timeout uint) (map[uint64]*common.HeightVoteSet, error) {

	hvsMap := make(map[uint64]*common.HeightVoteSet)

	// prepare and send data request
	err := broadcastHVSRequest(connections)

	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}

	for _, conn := range connections {

		wg.Add(1)

		// Launch a goroutine to fetch the hvs
		go func(conn net.Conn) {

			packet, err := connection.Receive(conn)

			if err != nil {
				fmt.Printf("Monitor: error while receiving hvs from validator: %s", err)
			} else if packet == nil || packet.Code != connection.HvsResponse || packet.Hvs == nil {
				fmt.Println("Monitor: invalid packet received from " + conn.RemoteAddr().String())
			} else {
				fmt.Println("Monitor: received hvs from " + conn.RemoteAddr().String())
				hvsMap[packet.Hvs.OwnerID] = packet.Hvs
			}
			wg.Done()

		}(conn)
	}

	if waitTimeout(&wg, timeout, connections) {
		return nil, fmt.Errorf("timed out waiting for wait group")
	}

	return hvsMap, nil

}

// waitTimeout waits for the waitgroup for the specified max timeout and returns true if waiting timed out
func waitTimeout(wg *sync.WaitGroup, timeout uint, connections []net.Conn) bool {
	closeChannel := make(chan struct{})

	// start timer for repeating request
	repeatTimer := time.NewTicker(time.Duration(timeout/3) * time.Second)
	defer repeatTimer.Stop()

	go func() {
		defer close(closeChannel)
		wg.Wait()
	}()

	select {

	case <-closeChannel:
		// completed normally
		return false

	case <-repeatTimer.C:
		// repeat request
		_ = broadcastHVSRequest(connections)

	case <-time.After(time.Duration(timeout) * time.Second):
		// timed out
		return true
	}

	return false
}

func broadcastHVSRequest(connections []net.Conn) error {

	packet := &connection.Packet{Code: connection.HvsRequest}

	for _, conn := range connections {
		err := connection.Send(conn, packet)
		if err != nil {
			return fmt.Errorf("Error while sending packet to validator "+conn.RemoteAddr().String()+": %s", err)
		}
	}

	return nil
}
