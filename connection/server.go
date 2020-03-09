package connection

import (
	"fmt"
	"net"
)

// Listen starts listening for incoming connections from the client monitor
func Listen(port string, monitorAddr string) error {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		return fmt.Errorf("Error while trying to listen on given address: %s", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			return fmt.Errorf("Error while trying to accept incoming connection: %s", err)
		}

		// handle connection in a separate goroutine
		go handleConnection(conn, monitorAddr)
	}
}

// handle connection
func handleConnection(conn net.Conn, monitorAddr string) {

	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Handling client connection from " + remoteAddr)

	// verify that client is monitor
	//if remoteAddr == monitorAddr {

	for {
		packet, err := Receive(conn)

		if err != nil {
			break
		}

		switch packet.Code {
		case HvsRequest:
			Send(conn, &Packet{Code: HvsResponse})

		default:
			// ignore and just go on
			fmt.Println("Unknown packet received, just ignore")
		}

	}
	//}

	conn.Close()
}
