package main

import (
	"log"
	"time"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
)

// Validator struct
type Validator struct {
	ID       string                           `yaml:"id"`
	Address  string                           `yaml:"address"`
	Messages map[uint64]*common.HeightVoteSet `yaml:"messages"`

	server *connection.Server // server
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
	log.Println("Validator on " + validator.Address + ": start listening for incoming requests")

	// handle incoming data from client monitor
	go validator.handleIncomingClientData(delay)

	// start listening for incoming connection from monitor
	err := validator.server.Listen(validator.Address)
	if err != nil {
		log.Fatalf("Validator %s exiting: cannot listen on given address: %s", validator.Address, err)
	}
}

// process packet from client (monitor)
func (validator *Validator) handleIncomingClientData(delay uint64) {
	// process client data from server channel
	for clientData := range validator.server.ReceiveChannel {

		time.Sleep(time.Duration(delay) * time.Second)

		packet := clientData.Packet

		// if it's a request packet, send the response back
		if packet != nil && packet.Code == connection.HvsRequest {
			log.Println("Validator on " + validator.Address + ": sending hvs to monitor")

			// prepare packet
			packet.Code = connection.HvsResponse
			packet.Hvs = validator.Messages[packet.Height]
			packet.ID = validator.ID

			err := clientData.Connection.Send(packet)
			if err != nil {
				log.Printf("Error while sending packet back to monitor: %s", err)
			}
		}
	}
}
