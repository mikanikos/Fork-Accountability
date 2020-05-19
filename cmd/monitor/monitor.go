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
	var output string
	if asyncMode {
		output = monitor.runAccountabilityAlgorithmAsync()
	} else {
		output = monitor.runAccountabilityAlgorithm()
	}

	if debug {
		log.Println(output)
	}
}

// run monitor algorithm
func (monitor *Monitor) runAccountabilityAlgorithm() string {

	// count the number of responses (regardless of validity) from different validators
	numValidators := len(monitor.Validators)

	// initialize accountability
	monitor.accAlgorithm.Init(uint64(numValidators), false)

	// wait until the specified timer expires
	timer := time.NewTicker(time.Duration(monitor.Timeout) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if debug {
				log.Println("Monitor: running the accountability algorithm")
			}

			start := time.Now()

			// run monitor and get faulty processes
			monitor.accAlgorithm.Run(monitor.FirstDecisionRound, monitor.SecondDecisionRound)

			elapsedTime := time.Since(start)

			log.Println("Monitor: algorithm completed in " + elapsedTime.String())

			// print result of the execution
			if debug {
				log.Println(monitor.accAlgorithm.String())
				log.Printf("Monitor: detected %d faulty processes\n", monitor.accAlgorithm.GetNumFaulty())
			}

			// if we have at least f + 1 faulty processes, the algorithm completed
			if monitor.accAlgorithm.IsCompleted() {
				return successfulStatus
			}

			return failStatus

		case packet := <-monitor.receiveChannel:

			// check if new packet has been received and store it in case
			if monitor.checkResponseValidity(packet) && monitor.accAlgorithm.StoreHvs(packet.ID, packet.Hvs) {
				if debug {
					log.Printf("Monitor: received height vote set from validator with ID %s. %d message logs have been delivered so far\n", packet.ID, monitor.accAlgorithm.GetNumLogs())
				}
			} else {
				if debug {
					log.Printf("Monitor: received invalid packet from validator with ID %s\n", packet.ID)
				}
			}
		}
	}
}

// run monitor algorithm
func (monitor *Monitor) runAccountabilityAlgorithmAsync() string {

	// count the number of responses (regardless of validity) from different validators
	numValidators := len(monitor.Validators)
	responseCount := 0

	// initialize accountability
	monitor.accAlgorithm.Init(uint64(numValidators), true)

	// wait until the specified timer expires
	timer := time.NewTicker(time.Duration(monitor.Timeout) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// exit because timer expired
			return timeoutStatus

		case packet := <-monitor.receiveChannel:

			// check if new packet has been received and store it in case
			if monitor.checkResponseValidity(packet) && monitor.accAlgorithm.StoreHvs(packet.ID, packet.Hvs) {

				if debug {
					log.Printf("Monitor: received height vote set from validator with ID %s. %d message logs have been delivered so far\n", packet.ID, monitor.accAlgorithm.GetNumLogs())
				}

				// if we have delivered at least f + 1 message logs, run the monitor algorithm
				if monitor.accAlgorithm.CanRun() {

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

					// if we have at least f + 1 faulty processes, the algorithm completed
					if monitor.accAlgorithm.IsCompleted() {
						return successfulStatus
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
				// exit because no new hvs will arrive and avoid waiting longer
				return failStatus
			}
		}
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
	if err != nil {
		if debug {
			log.Printf("Monitor: error while sending request to %s: %s", conn.Conn.RemoteAddr().String(), err)
		}
	}

	// wait to receive packet from validator
	packet, err := conn.Receive()
	if err != nil {
		// if connection is closed or there's an error, exit
		if err == io.EOF {
			if debug {
				log.Printf("Monitor: connection has been closed by validator on address %s", conn.Conn.RemoteAddr())
			}
		} else {
			if debug {
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
