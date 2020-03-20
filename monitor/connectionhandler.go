package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
)

// ConnectionHandler handles the connection with validators
type ConnectionHandler struct {
	connections []net.Conn
}

// NewConnectionHandler creates a new handler for the connection with validators
func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		connections: make([]net.Conn, 0),
	}
}

// method to resolve processes addresses and store connection objects
func (connHandler *ConnectionHandler) connectToValidators(validators string) error {

	// split list of string addresses only if it's not empty in order to avoid problems
	validatorsList := make([]string, 0)

	if validators != "" {
		validatorsList = strings.Split(validators, ",")
	}

	// resolve peers addresses given
	validatorsConn := make([]net.Conn, 0)
	for _, val := range validatorsList {
		conn, err := connection.Connect(val)
		if err == nil {
			fmt.Println("Monitor: connected to " + conn.RemoteAddr().String())
			validatorsConn = append(validatorsConn, conn)
		} else {
			return fmt.Errorf("error while connecting to one of the validators given: %s", err)
		}
	}

	return nil
}

// request HeightVoteSets from validators with a max timeout
// if a validator doesn't send its hvs, the monitor will consider it faulty
func (connHandler *ConnectionHandler) requestHVSWithTimeout(timeout uint) (map[uint64]*common.HeightVoteSet, error) {

	hvsMap := make(map[uint64]*common.HeightVoteSet)

	// prepare and send data request
	err := connHandler.broadcastHVSRequest()
	if err != nil {
		return nil, err
	}

	// wait group to wait for responses
	wg := sync.WaitGroup{}

	// start waiting for every connection
	for _, conn := range connHandler.connections {
		wg.Add(1)

		// Launch a goroutine to fetch the hvs
		go func(conn net.Conn) {
			// receive data from validator
			packet, err := connection.Receive(conn)

			if err != nil {
				fmt.Printf("Monitor: error while receiving hvs from validator: %s", err)
			} else if packet == nil || packet.Code != connection.HvsResponse || packet.Hvs == nil {
				fmt.Println("Monitor: invalid packet received from " + conn.RemoteAddr().String())
			} else {
				fmt.Println("Monitor: received hvs from " + conn.RemoteAddr().String())
				hvsMap[packet.Hvs.OwnerID] = packet.Hvs
			}
			wg.Done()

		}(conn)
	}

	// wait routine, it completes after the timeout or as soon as we receive all the hvs
	if connHandler.waitTimeout(&wg, timeout) {
		fmt.Println("timed out waiting for wait group, not all hvs were sent")
	}

	return hvsMap, nil
}

// waitTimeout waits for the WaitGroup for the specified max timeout and returns true if waiting timed out
func (connHandler *ConnectionHandler) waitTimeout(wg *sync.WaitGroup, timeout uint) bool {
	closeChannel := make(chan struct{})

	// start timer for repeating request three times
	repeatTimer := time.NewTicker(time.Duration(timeout/3) * time.Second)
	defer repeatTimer.Stop()

	go func() {
		defer close(closeChannel)
		wg.Wait()
	}()

	for {
		select {

		case <-closeChannel:
			// completed normally
			return true

		case <-repeatTimer.C:
			// repeat request
			_ = connHandler.broadcastHVSRequest()

		case <-time.After(time.Duration(timeout) * time.Second):
			// timed out
			return true
		}
	}
}

// broadcast hvs request to all validators
func (connHandler *ConnectionHandler) broadcastHVSRequest() error {

	// prepare packet
	packet := &connection.Packet{Code: connection.HvsRequest}

	// broadcast message
	for _, conn := range connHandler.connections {
		err := connection.Send(conn, packet)
		if err != nil {
			return fmt.Errorf("Error while sending packet to validator "+conn.RemoteAddr().String()+": %s", err)
		}
	}

	return nil
}
