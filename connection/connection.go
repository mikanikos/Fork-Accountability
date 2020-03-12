package connection

import (
	"fmt"
	"net"

	"go.dedis.ch/protobuf"
)

// Send sends a packet to a given connection
func Send(conn net.Conn, packet *Packet) error {

	// encode message
	messageEncoded, err := protobuf.Encode(packet)

	if err != nil {
		return fmt.Errorf("Error while serializing the packet to send: %s", err)
	}

	// send message
	_, err = conn.Write(messageEncoded)
	return err
}

// Receive receives a packet from a given connection
func Receive(conn net.Conn) (*Packet, error) {

	packet := &Packet{}
	packetBytes := make([]byte, maxBufferSize)

	n, err := conn.Read(packetBytes)

	if err != nil {
		return nil, fmt.Errorf("Error while reading from socket: %s", err)
	}

	if n > maxBufferSize {
		return nil, fmt.Errorf("Error while reading from socket: message size too large")
	}

	// decode message
	err = protobuf.Decode(packetBytes[:n], packet)

	if err != nil {
		return nil, fmt.Errorf("Error while deserializing packet received: %s", err)
	}

	return packet, nil
}