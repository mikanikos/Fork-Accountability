package connection

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Server object to handle requests from clients
type Server struct {
	ReceiveChannel chan *ClientData
}

// NewServer creates a new Server
func NewServer() *Server {
	return &Server{ReceiveChannel: make(chan *ClientData, maxChannelSize)}
}

// ClientData is the data sent by the client to be delivered to the Listener
type ClientData struct {
	Packet     *Packet
	Connection *Connection
}

// Listen starts listening for incoming connections from the client
func (server *Server) Listen(address string) error {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return fmt.Errorf("error while trying to listen on given address: %s", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			return fmt.Errorf("error while trying to accept incoming connection: %s", err)
		}

		// handle connection in a separate goroutine
		go server.HandleConnection(&Connection{Conn: conn})
	}
}

// HandleConnection from the given connection
func (server *Server) HandleConnection(connection *Connection) {

	log.Println("Handling client connection from " + connection.Conn.RemoteAddr().String())

	for {
		packet, err := connection.Receive()

		if err != nil {
			if err == io.EOF {
				log.Printf("Client %s closed the connection", connection.Conn.RemoteAddr())
			} else {
				log.Printf("error while trying to receive packet from %s: %s", connection.Conn.RemoteAddr(), err)
			}
			return
		}

		clientData := &ClientData{Packet: packet, Connection: connection}

		// send data to receiving channel without blocking
		go func(data *ClientData) {
			server.ReceiveChannel <- data
		}(clientData)
	}
}
