package main

import (
	"fmt"
	"log"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
)

// Validator struct
type Validator struct {
	ID       uint64                           `yaml:"id"`
	Address  string                           `yaml:"address"`
	Messages map[uint64]*common.HeightVoteSet `yaml:"messages"`

	Server *connection.Server // server
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		Messages: make(map[uint64]*common.HeightVoteSet),
		Server:   connection.NewServer(),
	}
}

// Run validator
func (validator *Validator) Run() {
	fmt.Println("Validator on " + validator.Address + ": start listening for incoming requests")

	// handle incoming data from client monitor
	go validator.handleIncomingClientData()

	// start listening for incoming connection from monitor
	err := validator.Server.Listen(validator.Address)
	if err != nil {
		log.Fatalf("Validator %s exiting: cannot listen on given address: %s", validator.Address, err)
	}
}

// process packet from client (monitor)
func (validator *Validator) handleIncomingClientData() {
	// process client data from server channel
	for clientData := range validator.Server.ReceiveChannel {
		packet := clientData.Packet

		// if it's a request packet, send the response back
		if packet.Code == connection.HvsRequest {
			fmt.Println("Validator on " + validator.Address + ": sending hvs to monitor")

			// prepare packet
			packet.Code = connection.HvsResponse
			packet.Hvs = validator.Messages[packet.Height]

			err := clientData.Connection.Send(packet)
			if err != nil {
				log.Println("Error while sending packet back to monitor")
			}
		}
	}
}
