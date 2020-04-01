package main

import (
	"fmt"
	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
)

const sendTimer = 1

// ConnectionHandler handles the connection with validators
type ConnectionHandler struct {
	connections    []*connection.Connection
	receiveChannel chan *common.HeightVoteSet
}

// NewConnectionHandler creates a new handler for the connection with validators
func NewConnectionHandler(channelSize int) *ConnectionHandler {
	return &ConnectionHandler{
		connections:    make([]*connection.Connection, 0),
		receiveChannel: make(chan *common.HeightVoteSet, channelSize),
	}
}

// method to resolve processes addresses and store connection objects
func (connHandler *ConnectionHandler) connectToValidators(validators []string) error {

	// resolve validator addresses given and connect to them
	for _, val := range validators {
		conn, err := connection.Connect(val)
		if err != nil {
			return fmt.Errorf("error while connecting to one of the validators given: %s", err)
		}
		connHandler.connections = append(connHandler.connections, conn)
	}

	return nil
}

// request async HeightVoteSets from validators
func (connHandler *ConnectionHandler) requestHeightLogs(height uint64) {

	// prepare packet to send
	packet := &connection.Packet{Code: connection.HvsRequest, Height: height}

	// start goroutines to send message and wait for reply for each validator
	for _, conn := range connHandler.connections {
		// periodically send packet to validator until we receive it
		server := connection.NewServer()
		validatorCloseChannel := make(chan bool)

		// receive packets from validator
		go server.HandleConnection(conn)
		go connHandler.handleIncomingClientData(server, validatorCloseChannel)

		// periodic send request to validator
		go conn.PeriodicSend(packet, validatorCloseChannel, sendTimer)
	}
}

// process packet from client (monitor)
func (connHandler *ConnectionHandler) handleIncomingClientData(server *connection.Server, validatorCloseChannel chan bool) {
	// process client data from server channel
	for clientData := range server.ReceiveChannel {
		packet := clientData.Packet

		// check if packet and its data are valid
		if packet != nil && packet.Code == connection.HvsResponse && packet.Hvs != nil {

			// send it to monitor
			go func(p *connection.Packet) {
				connHandler.receiveChannel <- p.Hvs
				close(validatorCloseChannel)
			}(packet)
		}
	}
}
