package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mikanikos/Fork-Accountability/accountability"
	"github.com/mikanikos/Fork-Accountability/connection"
	"github.com/mikanikos/Fork-Accountability/utils"
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
	receiveChannel chan *connection.Packet
	// create accountability structure
	accAlgorithm *accountability.Accountability
}

// NewMonitor creates a new monitor
func NewMonitor() *Monitor {
	return &Monitor{
		Validators:     make([]string, 0),
		connections:    make([]*connection.Connection, 0),
		receiveChannel: make(chan *connection.Packet, maxChannelSize),
		accAlgorithm:   accountability.NewAccountability(),
	}
}

// Run monitor algorithm
func (monitor *Monitor) Run(report string) {

	// write logs to file, if desired
	if report != "" {
		f, err := utils.OpenFile(report)
		if err != nil {
			log.Fatalf("Monitor exiting: error opening report file: %s", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// connect to validators for requesting hvs
	err := monitor.connectToValidators()
	if err != nil {
		log.Fatalf("Monitor exiting: couldn't connect to all validators: %s", err)
	}

	log.Println("Monitor successfully connected to all validators")

	// make request for hvs to validators
	monitor.requestHeightLogs()

	log.Println("Monitor started sending packets to all validators")

	// run accountability algorithm
	output := monitor.runAccountabilityAlgorithm()
	log.Println(output)

	log.Println(monitor.accAlgorithm.String())
}

func (monitor *Monitor) runAccountabilityAlgorithm() string {
	numValidators := len(monitor.Validators)

	// lower bound on the number of faulty processes
	threshold := (numValidators-1)/3 + 1 // f+1

	// count the number of responses (regardless of validity) from different validators
	responseCount := 0

	// wait until the specified timer expires
	timer := time.NewTicker(time.Duration(monitor.Timeout) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// exit because timer expired
			return timeoutStatus

		case packet := <-monitor.receiveChannel:

			// check if new packet has been received
			if packet != nil {

				// increments the number of valid packets received
				monitor.accAlgorithm.MessageLogsReceived++
				// add hvs for the validator who sent it
				monitor.accAlgorithm.HeightLogs.AddHvs(packet.ID, packet.Hvs)

				log.Printf("Monitor: received height vote set from process %s\n", packet.ID)
				log.Printf("Monitor received %d message logs\n", monitor.accAlgorithm.MessageLogsReceived)

				// if we have delivered at least f + 1 message logs, run the monitor algorithm
				if monitor.accAlgorithm.CanRun(threshold) {

					log.Println("Monitor started running the accountability algorithm")

					// run monitor and get faulty processes
					monitor.accAlgorithm.Run(uint64(numValidators), monitor.FirstDecisionRound, monitor.SecondDecisionRound)

					log.Printf("Monitor detected %d faulty processes\n", monitor.accAlgorithm.FaultySet.Length())

					// if we have at least f + 1 faulty processes, the algorithm completed
					if monitor.accAlgorithm.IsCompleted(threshold) {
						return successfulStatus
					}
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

	// start goroutines to send message and wait for reply from each validator
	for _, conn := range monitor.connections {

		// receive packets from validator
		go monitor.receiveHvsFromValidator(conn)
	}
}

// receive hvs from validator
func (monitor *Monitor) receiveHvsFromValidator(conn *connection.Connection) {

	packet := &connection.Packet{Code: connection.HvsRequest, Height: monitor.Height}

	// try receive packet until a valid packet is sent
	for {

		log.Printf("Monitor: sending packet to %s", conn.Conn.RemoteAddr().String())

		// sending packet to validator
		err := conn.Send(packet)
		if err != nil {
			log.Printf("Error while sending request to %s: %s", conn.Conn.RemoteAddr().String(), err)
		}

		// wait to receive packet from validator
		packet, err := conn.Receive()
		if err != nil {
			// if connection is closed, exit
			if err == io.EOF {
				log.Printf("Connection has been closed by validator on address %s", conn.Conn.RemoteAddr())
			} else {
				log.Printf("Error while trying to receive packet from %s: %s", conn.Conn.RemoteAddr(), err)
			}

			// notify that a validator disconnected and will not receive any hvs from it
			monitor.receiveChannel <- nil
			return
		}

		// check if packet and its data are valid
		if packet != nil && packet.Code == connection.HvsResponse && packet.Hvs != nil && packet.Height == monitor.Height {

			// notify the monitor that new hvs has arrived
			monitor.receiveChannel <- packet
			return
		}
	}
}
