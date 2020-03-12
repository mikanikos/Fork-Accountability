package connection

import (
	"fmt"
	"net"

	"github.com/mikanikos/Fork-Accountability/common"
)

// Listen starts listening for incoming connections from the client monitor
func Listen(port string, hvs *common.HeightVoteSet) error {
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
		go handleConnection(conn, hvs)
	}
}

// handle connection
func handleConnection(conn net.Conn, hvs *common.HeightVoteSet) {

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
			fmt.Println("Validator on " + conn.LocalAddr().String() + ": sending hvs to monitor")
			Send(conn, &Packet{Code: HvsResponse, Hvs: hvs})

		default:
			// ignore and just go on
			fmt.Println("Unknown packet received, just ignore")
		}
	}
	//}

	conn.Close()
}
