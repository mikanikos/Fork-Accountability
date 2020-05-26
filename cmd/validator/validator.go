package main

import (
	"log"
	"time"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
)

const debug = true

// Validator struct
type Validator struct {
	ID       string                           `yaml:"id"`
	Address  string                           `yaml:"address"`
	Messages map[uint64]*common.HeightVoteSet `yaml:"messages"`

	// server
	server *connection.Server
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		Messages: make(map[uint64]*common.HeightVoteSet),
		server:   connection.NewServer(),
	}
}

// Run validator
func (validator *Validator) Run(delay uint64) {
	if debug {
		log.Printf("Validator %s at %s: start listening for incoming requests", validator.ID, validator.Address)
	}

	// handle incoming data from clients
	go validator.handleIncomingClientData(delay)

	// start listening for incoming connection from monitor
	err := validator.server.Listen(validator.Address)
	if err != nil {
		log.Fatalf("Validator %s at %s exiting: cannot listen on given address: %s", validator.ID, validator.Address, err)
	}
}

// process packet from client (monitor)
func (validator *Validator) handleIncomingClientData(delay uint64) {
	// process client data from server channel
	for clientData := range validator.server.ReceiveChannel {

		// wait some time to answer back, if requested
		time.Sleep(time.Duration(delay) * time.Second)

		// get packet and connection
		packet := clientData.Packet
		conn := clientData.Connection

		// if it's a request packet, send the response back
		if packet != nil && packet.Code == connection.HvsRequest {

			if debug {
				log.Printf("Validator %s at %s: received request for height vote set for height %d", validator.ID, validator.Address, packet.Height)
			}

			// load height vote set
			hvs, loaded := validator.Messages[packet.Height]

			// if validator does not have any message log for requested height, send error message
			if hvs != nil && loaded {

				// prepare packet
				packet.ID = validator.ID
				packet.Code = connection.HvsResponse
				packet.Hvs = hvs

				if debug {
					log.Printf("Validator %s at %s: sending height vote set requested for height %d to monitor", validator.ID, validator.Address, packet.Height)
				}

				// send response
				err := conn.Send(packet)
				if err != nil && debug {
					log.Printf("Validator %s at %s: error while sending packet back to monitor: %s", validator.ID, validator.Address, err)
				}
			} else {
				if debug {
					log.Printf("Validator %s at %s does not have any message logs for height %d", validator.ID, validator.Address, packet.Height)
				}
			}
		}
	}
}
