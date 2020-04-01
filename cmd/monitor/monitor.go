package main

import (
	"fmt"
	"github.com/mikanikos/Fork-Accountability/accountability"
	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
	"log"
)

const (
	sendTimer      = 1
	maxChannelSize = 100
)

// Monitor struct
type Monitor struct {
	Height              uint64   `yaml:"height"`
	FirstDecisionRound  uint64   `yaml:"firstDecisionRound"`
	SecondDecisionRound uint64   `yaml:"secondDecisionRound"`
	Validators          []string `yaml:"validators"`

	connections    []*connection.Connection
	receiveChannel chan *common.HeightVoteSet
}

// NewMonitor creates a new monitor
func NewMonitor() *Monitor {
	return &Monitor{
		Validators:     make([]string, 0),
		connections:    make([]*connection.Connection, 0),
		receiveChannel: make(chan *common.HeightVoteSet, maxChannelSize),
	}
}

func (monitor *Monitor) Run() {
	if monitor.Validators == nil || len(monitor.Validators) == 0 {
		log.Fatal("Monitor exiting: no validators given")
	}

	// connect to validators for requesting hvs
	err := monitor.connectToValidators(monitor.Validators)
	if err != nil {
		log.Fatalf("Monitor exiting: couldn't connect to all validators: %s", err)
	}

	// make request for hvs to validators
	monitor.requestHeightLogs(monitor.Height)

	// run accountability algorithm
	monitor.runAccountabilityAlgorithm()
}

func (monitor *Monitor) runAccountabilityAlgorithm() {
	numValidators := len(monitor.Validators)

	// create faulty set structure
	acc := accountability.NewAccountability()
	// lower bound on the number of faulty processes
	minFaulty := (numValidators-1)/3 + 1 // f+1

	// run until we have at least f+1 faulty processes
	for acc.FaultySet.Length() < minFaulty {
		// receive hvs from processes, it blocks the execution until another hvs arrives
		hvs := <-monitor.receiveChannel
		acc.HeightLogs.AddHvs(hvs)

		// if we have at least f+1 hvs, run the monitor algorithm
		if acc.HeightLogs.Length() >= minFaulty {
			// run monitor and get faulty processes
			acc.IdentifyFaultyProcesses(uint64(numValidators), monitor.FirstDecisionRound, monitor.SecondDecisionRound)
		}
	}

	fmt.Println(acc.HeightLogs.String())
	fmt.Println()
	fmt.Println(acc.FaultySet.String())
	fmt.Println("Monitor: Algorithm completed")

}

// method to resolve processes addresses and store connection objects
func (monitor *Monitor) connectToValidators(validators []string) error {

	// resolve validator addresses given and connect to them
	for _, val := range validators {
		conn, err := connection.Connect(val)
		if err != nil {
			return fmt.Errorf("error while connecting to one of the validators given: %s", err)
		}
		monitor.connections = append(monitor.connections, conn)
	}

	return nil
}

// request async HeightVoteSets from validators
func (monitor *Monitor) requestHeightLogs(height uint64) {

	// prepare packet to send
	packet := &connection.Packet{Code: connection.HvsRequest, Height: height}

	// start goroutines to send message and wait for reply for each validator
	for _, conn := range monitor.connections {
		// periodically send packet to validator until we receive it
		server := connection.NewServer()
		validatorCloseChannel := make(chan bool)

		// receive packets from validator
		go server.HandleConnection(conn)
		go monitor.handleIncomingClientData(server, validatorCloseChannel)

		// periodic send request to validator
		go conn.PeriodicSend(packet, validatorCloseChannel, sendTimer)
	}
}

// process packet from client (monitor)
func (monitor *Monitor) handleIncomingClientData(server *connection.Server, validatorCloseChannel chan bool) {
	// process client data from server channel
	for clientData := range server.ReceiveChannel {
		packet := clientData.Packet

		// check if packet and its data are valid
		if packet != nil && packet.Code == connection.HvsResponse && packet.Hvs != nil {

			// send it to monitor
			go func(p *connection.Packet) {
				monitor.receiveChannel <- p.Hvs
				close(validatorCloseChannel)
			}(packet)
		}
	}
}
