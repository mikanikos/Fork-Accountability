package main

import (
	"fmt"
	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
	"net"
	"strings"
	"sync"
	"time"
)

// ConnectionHandler handles the connection with validators
type ConnectionHandler struct {
	connections []net.Conn
}

// NewConnectionHandler creates a new handler for the connection with validators
func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		connections:    make([]net.Conn, 0),
	}
}

// resolve processes addresses
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
			validatorsConn = append(validatorsConn, conn)
		} else {
			return fmt.Errorf("error while connecting to one of the validators given: %s", err)
		}
	}

	return nil
}

func (connHandler *ConnectionHandler) requestHVSWithTimeout(timeout uint) (map[uint64]*common.HeightVoteSet, error) {

	hvsMap := make(map[uint64]*common.HeightVoteSet)

	// prepare and send data request
	err := connHandler.broadcastHVSRequest()

	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}

	for _, conn := range connHandler.connections {

		wg.Add(1)

		// Launch a goroutine to fetch the hvs
		go func(conn net.Conn) {

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

	if connHandler.waitTimeout(&wg, timeout) {
		fmt.Println("timed out waiting for wait group, not all hvs were sent")
	}

	return hvsMap, nil
}

// waitTimeout waits for the waitgroup for the specified max timeout and returns true if waiting timed out
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

func (connHandler *ConnectionHandler) broadcastHVSRequest() error {

	packet := &connection.Packet{Code: connection.HvsRequest}

	for _, conn := range connHandler.connections {
		err := connection.Send(conn, packet)
		if err != nil {
			return fmt.Errorf("Error while sending packet to validator "+conn.RemoteAddr().String()+": %s", err)
		}
	}

	return nil
}
