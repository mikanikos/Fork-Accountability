package main

import (
	"fmt"
	"github.com/mikanikos/Fork-Accountability/connection"
	"io"
	"net"
)

// Listen starts listening for incoming connections from the client monitor
func (validator *Validator) Listen() error {
	listener, err := net.Listen("tcp", validator.Address)

	if err != nil {
		return fmt.Errorf("error while trying to listen on given address: %s", err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			_ = listener.Close()
			return fmt.Errorf("error while trying to accept incoming connection: %s", err)
		}

		// handle connection in a separate goroutine
		go validator.handleConnection(conn)
	}
}

// handle connection
func (validator *Validator) handleConnection(conn net.Conn) {

	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Handling client connection from " + remoteAddr)

	for {
		packet, err := connection.Receive(conn)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Monitor closed the connection")
			} else {
				fmt.Printf("error while trying to receive packet: %s", err)
			}
			break
		}

		switch packet.Code {
		case connection.HvsRequest:
			fmt.Println("Validator on " + conn.LocalAddr().String() + ": sending hvs to monitor")
			packet.Code = connection.HvsResponse
			packet.Hvs = validator.Messages[packet.Height]
			err := connection.Send(conn, packet)
			if err != nil {
				fmt.Println("Error while sending packet back to monitor")
				break
			}
		}
	}

	_ = conn.Close()
}
