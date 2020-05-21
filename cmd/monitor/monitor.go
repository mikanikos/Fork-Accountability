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

	// receive channel for incoming packets
	receiveChannel chan *connection.Packet
	// accountability structure
	accAlgorithm *accountability.Accountability
}

// NewMonitor creates a new monitor
func NewMonitor() *Monitor {
	return &Monitor{
		Validators:     make([]string, 0),
		receiveChannel: make(chan *connection.Packet, maxChannelSize),
		accAlgorithm:   accountability.NewAccountability(),
	}
}

// Run monitor algorithm
func (monitor *Monitor) Run(report string, asyncMode bool) {

	// write logs to file, if desired
	if report != "" {
		f, err := utils.OpenFile(report)
		if err != nil {
			log.Fatalf("Monitor exiting: error opening report file: %s", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	if debug {
		log.Println("Monitor: started running")
	}

	// connect to validators and make request for hvs
	err := monitor.connectToValidators()
	if err != nil {
		log.Fatalf("Monitor exiting: couldn't connect to all validators: %s", err)
	}

	if debug {
		log.Println("Monitor: successfully connected to all validators and requested message logs")
	}

	// run accountability algorithm
	output := monitor.runMonitorAlgorithm(asyncMode)

	if debug {
		log.Println(output)
	}
}

// run monitor algorithm
func (monitor *Monitor) runMonitorAlgorithm(async bool) string {

	// count the number of responses (regardless of validity) from different validators
	numValidators := len(monitor.Validators)
	responseCount := 0

	// initialize accountability
	monitor.accAlgorithm.Init(uint64(numValidators), async)

	// wait until the specified timer expires
	timer := time.NewTicker(time.Duration(monitor.Timeout) * time.Second)
	defer timer.Stop()

loop:
	for {
		select {
		case <-timer.C:
			if async {
				// safety timeout
				return timeoutStatus
			} else {
				// timeout for receiving hvs, exit and run the algorithm
				break loop
			}

		case packet := <-monitor.receiveChannel:

			// check if new packet has been received and store it in case
			if monitor.checkResponseValidity(packet) && monitor.accAlgorithm.StoreHvs(packet.ID, packet.Hvs) {
				if debug {
					log.Printf("Monitor: received height vote set from validator with ID %s. %d message logs have been delivered so far\n", packet.ID, monitor.accAlgorithm.GetNumLogs())
				}

				// run the algorithm, if possible, in the asynchronous version
				if async {
					// if we have delivered at least f + 1 message logs, run the monitor algorithm
					if monitor.accAlgorithm.CanRun() {

						// run algorithm
						monitor.runAccountabilityAlgorithm()

						// if we have at least f + 1 faulty processes, the algorithm completed correctly
						if monitor.accAlgorithm.IsCompleted() {
							return successfulStatus
						}
					}
				}

			} else {
				if debug {
					log.Printf("Monitor: received invalid packet from validator with ID %s\n", packet.ID)
				}
			}

			// increment the number of responses from validators
			responseCount++
			if responseCount == numValidators {
				if async {
					// fail because no new hvs will arrive and the success condition was not met
					return failStatus
				} else {
					// exit and run the algorithm because all hvs have been delivered and avoid waiting longer
					break loop
				}
			}
		}
	}

	// run algorithm
	monitor.runAccountabilityAlgorithm()

	// if we have at least f + 1 faulty processes, the algorithm completed correctly
	if monitor.accAlgorithm.IsCompleted() {
		return successfulStatus
	}

	return failStatus
}

// run accountability algorithm
func (monitor *Monitor) runAccountabilityAlgorithm() {

	if debug {
		log.Println("Monitor: running the accountability algorithm")
	}

	start := time.Now()

	// run monitor and get faulty processes
	monitor.accAlgorithm.Run(monitor.FirstDecisionRound, monitor.SecondDecisionRound)

	elapsedTime := time.Since(start)

	log.Println("Monitor: algorithm completed in " + elapsedTime.String())

	if debug {
		// print result of the execution
		log.Println(monitor.accAlgorithm.String())
		log.Printf("Monitor: detected %d faulty processes\n", monitor.accAlgorithm.GetNumFaulty())
	}
}

// request message logs for a specific height to all validators, return error if it can't connect to some validator
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

		// start goroutines to send message and wait for reply from each validator
		go monitor.receiveHvsFromValidator(conn)
	}

	return nil
}

// receive hvs from validator on a given connection
func (monitor *Monitor) receiveHvsFromValidator(conn *connection.Connection) {

	// prepare packet to send
	packetToSend := &connection.Packet{Code: connection.HvsRequest, Height: monitor.Height}

	if debug {
		log.Printf("Monitor: sending packet to %s", conn.Conn.RemoteAddr().String())
	}

	// sending packet to validator
	err := conn.Send(packetToSend)
	if err != nil && debug {
		log.Printf("Monitor: error while sending request to %s: %s", conn.Conn.RemoteAddr().String(), err)
	}

	// wait to receive packet from validator
	packet, err := conn.Receive()
	if err != nil {
		// if connection is closed or there's an error, exit
		if debug {
			if err == io.EOF {
				log.Printf("Monitor: connection has been closed by validator on address %s", conn.Conn.RemoteAddr())
			} else {
				log.Printf("Monitor: error while trying to receive packet from %s: %s", conn.Conn.RemoteAddr(), err)
			}
		}

		// notify that will not receive any hvs from the validator on this connection
		packet = &connection.Packet{Code: connection.HvsMissing}
	}

	// send packet to main thread
	monitor.receiveChannel <- packet

	// close connection with validator
	conn.Close()
}

// check that the packet received is valid and contains correct information
func (monitor *Monitor) checkResponseValidity(packet *connection.Packet) bool {
	return packet != nil &&
		packet.Code == connection.HvsResponse &&
		packet.Height == monitor.Height &&
		packet.ID != "" &&
		packet.Hvs != nil &&
		packet.Hvs.IsValid(packet.ID)
}
