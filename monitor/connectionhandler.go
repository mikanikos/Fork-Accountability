package main

import (
	"fmt"
	"net"
	"time"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
)

// ConnectionHandler handles the connection with validators
type ConnectionHandler struct {
	connections    []net.Conn
	receiveChannel chan *common.HeightVoteSet
}

// NewConnectionHandler creates a new handler for the connection with validators
func NewConnectionHandler(channelSize int) *ConnectionHandler {
	return &ConnectionHandler{
		connections:    make([]net.Conn, 0),
		receiveChannel: make(chan *common.HeightVoteSet, channelSize),
	}
}

// method to resolve processes addresses and store connection objects
func (connHandler *ConnectionHandler) connectToValidators(validators []string) error {

	// resolve peers addresses given
	for _, val := range validators {
		conn, err := connection.Connect(val)
		if err == nil {
			fmt.Println("Monitor: connected to " + conn.RemoteAddr().String())
			connHandler.connections = append(connHandler.connections, conn)
		} else {
			return fmt.Errorf("error while connecting to one of the validators given: %s", err)
		}
	}

	return nil
}

// request async HeightVoteSets from validators
func (connHandler *ConnectionHandler) requestHeightLogs(height uint64) {

	// start waiting for every connection
	for _, connVal := range connHandler.connections {

		// Launch a goroutine to fetch the hvs from a validator
		go func(conn net.Conn) {

			// prepare packet
			packet := &connection.Packet{Code: connection.HvsRequest, Height: height}

			for {

				// wait a bit before resending
				time.Sleep(500)

				// send packet
				err := connection.Send(conn, packet)
				if err != nil {
					fmt.Println("Monitor: error while sending packet to validator "+conn.RemoteAddr().String()+": %s", err)
					continue
				}

				// receive data from validator
				packet, err := connection.Receive(conn)

				if err != nil {
					fmt.Printf("Monitor: error while receiving hvs from validator: %s", err)
				} else if packet == nil || packet.Code != connection.HvsResponse || packet.Hvs == nil {
					fmt.Println("Monitor: invalid packet received from " + conn.RemoteAddr().String())
				} else {
					fmt.Println("Monitor: received hvs from " + conn.RemoteAddr().String())

					go func(p *connection.Packet) {
						connHandler.receiveChannel <- p.Hvs
					}(packet)
					break
				}
			}

			err := conn.Close()
			if err != nil {
				fmt.Println("Monitor: error closing connection for " + conn.RemoteAddr().String())
			}

		}(connVal)
	}
}
