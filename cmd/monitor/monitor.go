package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mikanikos/Fork-Accountability/accountability"
	"github.com/mikanikos/Fork-Accountability/connection"
)

const sendTimer = 3

const (
	successfulStatus = "Monitor: Algorithm completed"
	failStatus = "Monitor: Algorithm failed because not enough message logs have been received or the message logs received were not sufficient to find at least f+1 faulty processes"
	timeoutStatus = "Monitor: Algorithm failed because of timeout expiration"
)

// Monitor struct
type Monitor struct {
	Height              uint64   `yaml:"height"`
	FirstDecisionRound  uint64   `yaml:"firstDecisionRound"`
	SecondDecisionRound uint64   `yaml:"secondDecisionRound"`
	Timeout             uint64   `yaml:"timeout"`
	Validators          []string `yaml:"validators"`

	// connections with all the validators
	connections    []*connection.Connection
	receiveChannel chan bool
	// create accountability structure
	accAlgorithm *accountability.Accountability
}

// NewMonitor creates a new monitor
func NewMonitor() *Monitor {
	return &Monitor{
		Validators:     make([]string, 0),
		connections:    make([]*connection.Connection, 0),
		receiveChannel: make(chan bool),
		accAlgorithm:   accountability.NewAccountability(),
	}
}

// Run monitor algorithm
func (monitor *Monitor) Run() {

	// connect to validators for requesting hvs
	err := monitor.connectToValidators()
	if err != nil {
		log.Fatalf("Monitor exiting: couldn't connect to all validators: %s", err)
	}

	// make request for hvs to validators
	monitor.requestHeightLogs()

	// run accountability algorithm
	output := monitor.runAccountabilityAlgorithm()
	log.Println(output)
}

func (monitor *Monitor) runAccountabilityAlgorithm() string {
	numValidators := len(monitor.Validators)

	// lower bound on the number of faulty processes
	minFaulty := (numValidators-1)/3 + 1 // f+1

	// count the number of responses from different validators
	responseCount := 0

	// wait until the specified timer expires
	timer := time.NewTicker(time.Duration(monitor.Timeout) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// exit because timer expired
			return timeoutStatus

		case status := <-monitor.receiveChannel:

			// check status received (new hvs or validator disconnected)
			// if we have delivered at least f + 1 message logs, run the monitor algorithm
			if status && monitor.accAlgorithm.HeightLogs.Length() >= minFaulty {

				// run monitor and get faulty processes
				monitor.accAlgorithm.IdentifyFaultyProcesses(uint64(numValidators), monitor.FirstDecisionRound, monitor.SecondDecisionRound)

				// if we have at least f + 1 faulty processes, the algorithm completed
				if monitor.accAlgorithm.FaultySet.Length() >= minFaulty {
					return successfulStatus
				}
			}

			// count the number of responses from validators
			responseCount++

			if responseCount == numValidators {
				// exit because no new hvs will arrive and avoid waiting longer
				return failStatus
			}
		}
	}
}

// method to resolve processes addresses and store connection objects
func (monitor *Monitor) connectToValidators() error {

	if monitor.Validators == nil || len(monitor.Validators) == 0 {
		return fmt.Errorf("error: no validators given")
	}

	// resolve validator addresses given and connect to them
	for _, val := range monitor.Validators {
		conn, err := connection.Connect(val)
		if err != nil {
			return fmt.Errorf("error while connecting to one of the validators given: %s", err)
		}
		monitor.connections = append(monitor.connections, conn)
	}

	return nil
}

// request async HeightVoteSets from validators
func (monitor *Monitor) requestHeightLogs() {

	// prepare packet to send
	packet := &connection.Packet{Code: connection.HvsRequest, Height: monitor.Height}

	// start goroutines to send message and wait for reply for each validator
	for _, conn := range monitor.connections {

		// channel used to close connection once hvs has been received and stop to send the request
		validatorCloseChannel := make(chan bool)

		// receive packets from validator
		go monitor.receiveHvsFromValidator(conn, validatorCloseChannel)

		// periodic send request to validator
		go conn.PeriodicSend(packet, validatorCloseChannel, sendTimer)
	}
}

// receive hvs from validator
func (monitor *Monitor) receiveHvsFromValidator(conn *connection.Connection, validatorCloseChannel chan bool) {

	// try receive packet until a valid packet is sent
	for {
		packet, err := conn.Receive()

		if err != nil {
			// if connection is closed, exit
			if err == io.EOF {
				log.Printf("Connection has been closed by validator on address %s", conn.Conn.RemoteAddr())
			} else {
				log.Printf("error while trying to receive packet from %s: %s", conn.Conn.RemoteAddr(), err)
			}

			// notify that a validator disconnected and will not receive any hvs from it
			monitor.receiveChannel <- false
			validatorCloseChannel <- false
			return
		}

		// check if packet and its data are valid
		if packet != nil && packet.Code == connection.HvsResponse && packet.Hvs != nil && !monitor.accAlgorithm.HeightLogs.Contains(packet.Hvs.OwnerID) {
			// add hvs for the validator who sent it
			monitor.accAlgorithm.HeightLogs.AddHvs(packet.Hvs)

			// notify the monitor that new hvs has arrived
			monitor.receiveChannel <- true
			validatorCloseChannel <- true
			return
		}
	}
}
