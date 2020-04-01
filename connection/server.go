package connection

import (
	"fmt"
	"io"
	"net"
)

type PacketConnection struct {
	Packet *Packet
	Conn   net.Conn
}

// Listen starts listening for incoming connections from the client monitor
func Listen(address string, receiveChannel chan *PacketConnection) error {
	listener, err := net.Listen("tcp", address)

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
		go handleConnection(conn, receiveChannel)
	}
}

// handle connection
func handleConnection(conn net.Conn, receiveChannel chan *PacketConnection) {

	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Handling client connection from " + remoteAddr)

	for {
		packet, err := Receive(conn)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Monitor closed the connection")
			} else {
				fmt.Printf("error while trying to receive packet: %s", err)
			}
			break
		}

		// send data to receiving channel
		receiveChannel <- &PacketConnection{
			Packet: packet,
			Conn:   conn,
		}
	}

	_ = conn.Close()
}
