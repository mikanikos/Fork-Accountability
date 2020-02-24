package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/mikanikos/Fork-Accountability/monitor"
)

func main() {

	// parsing arguments
	processes := flag.String("processes", "", "comma separated list of all processes addresses in the form ip:port")
	firstDecisionRound := flag.Uint64("firstDecisionRound", 0, "first round when a decision was taken")
	secondDecisionRound := flag.Uint64("secondDecisionRound", 0, "second round when a decision was taken (when fork was detected)")
	waitTimeout := flag.Uint("waitTimeout", 5, "timeout to wait for HeightVoteSet from each process")

	flag.Parse()

	processesAddresses := resolveAddresses(*processes)

	// request hvs from all processes
	hvsMap := requestHVSWithTimeout(processesAddresses, *waitTimeout)

	hvsList := make([]*monitor.HeightVoteSet, len(hvsMap))
	i := 0
	for _, hvs := range hvsMap {
		hvsList[i] = hvs
		i++
	}

	// run monitor and get faulty processes
	faultyProcesses := monitor.IdentifyFaultyProcesses(uint64(len(processesAddresses)), *firstDecisionRound, *secondDecisionRound, hvsList)

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
func resolveAddresses(validators string) []*net.UDPAddr {

	// split list of string addresses only if it's not empty in order to avoid problems
	validatorsList := make([]string, 0)
	if validators != "" {
		validatorsList = strings.Split(validators, ",")
	}

	// resolve peers addresses given
	validatorsAddresses := make([]*net.UDPAddr, 0)
	for _, val := range validatorsList {
		address, err := net.ResolveUDPAddr("udp4", val)
		if err == nil {
			validatorsAddresses = append(validatorsAddresses, address)
		}
	}

	return validatorsAddresses
}

func requestHVSWithTimeout(addresses []*net.UDPAddr, timeout uint) map[string]*monitor.HeightVoteSet {

	// prepare data request
	// packet := &GossipPacket{DataRequest: &DataRequest{Origin: gossiper.name, Destination: peer, HashValue: hash, HopLimit: uint32(hopLimit)}}

	// // send request
	// go gossiper.sendRequest(packet, &packet.DataRequest.HopLimit, packet.DataRequest.Destination)

	hvsMap := make(map[string]*monitor.HeightVoteSet)

	// start timer for repeating request
	timer := time.NewTicker(time.Duration(timeout) * time.Second)
	defer timer.Stop()

	for {
		select {
		// incoming reply for this request
		// case replyPacket := <-replyChan:

		// 	// save data
		// 	gossiper.fileHandler.hashDataMap.LoadOrStore(hex.EncodeToString(hash), &replyPacket.Data)

		// stop after timeout
		case <-timer.C:
			return hvsMap
		}
	}
}
