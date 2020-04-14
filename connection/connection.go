package connection

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"go.dedis.ch/protobuf"
)

// Connection is a wrapper for a net.Conn
type Connection struct {
	Conn net.Conn
}

// Send sends a packet to a given connection
func (c *Connection) Send(packet *Packet) error {

	// encode message
	messageEncoded, err := protobuf.Encode(packet)

	if err != nil {
		return fmt.Errorf("error while serializing the packet to send: %s", err)
	}

	// send message
	_, err = c.Conn.Write(messageEncoded)
	if err != nil {
		return fmt.Errorf("error while sending packet to %s: %s", c.Conn.RemoteAddr(), err)
	}

	return nil
}

// Receive receives a packet from a given connection
func (c *Connection) Receive() (*Packet, error) {

	packet := &Packet{}
	packetBytes := make([]byte, maxBufferSize)

	n, err := c.Conn.Read(packetBytes)

	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("error while reading from socket: %s", err)
	}

	if n > maxBufferSize {
		return nil, fmt.Errorf("error while reading from socket: message size too large")
	}

	// decode message
	err = protobuf.Decode(packetBytes[:n], packet)

	if err != nil {
		return nil, fmt.Errorf("error while deserializing packet received: %s", err)
	}

	return packet, nil
}

// Connect tried to establish connection given an address
func Connect(address string) (*Connection, error) {
	connClient, err := net.Dial("tcp", address)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to address %s: %s", address, err)
	}

	return &Connection{Conn: connClient}, nil
}

// Close tries to close a given connection
func (c *Connection) Close() {
	err := c.Conn.Close()
	if err != nil {
		log.Printf("Error while closing connection to address %s: %s", c.Conn.RemoteAddr().String(), err)
	}
}

// PeriodicSend periodically send a request at every timer tick
// the sending can be stopped using the channel given
func (c *Connection) PeriodicSend(packet *Packet, closeChannel chan bool, timer uint64) {

	err := c.Send(packet)
	if err != nil {
		log.Printf("Error while sending request to %s: %s", c.Conn.RemoteAddr().String(), err)
	}

	// start timer for repeating request to validator
	repeatTimer := time.NewTicker(time.Duration(timer) * time.Second)
	defer repeatTimer.Stop()

	for {
		select {

		// case status := <-closeChannel:
		// 	// stop because we received the packet from validator
		// 	if status {
		// 		c.Close()
		// 	}
		// 	return

		case <-repeatTimer.C:
			// repeat request
			err := c.Send(packet)
			if err != nil {
				log.Printf("Error while repeating request to %s: %s", c.Conn.RemoteAddr().String(), err)
			}
		}
	}

}
