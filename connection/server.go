package connection

import (
	"fmt"
	"io"
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

// Listen starts listening for incoming connections from the client monitor
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

// handle connection from client
func (server *Server) HandleConnection(connection *Connection) {

	fmt.Println("Handling client connection from " + connection.Conn.RemoteAddr().String())

	defer connection.Close()

	for {
		packet, err := connection.Receive()

		if err != nil {
			if err == io.EOF {
				fmt.Println("Monitor closed the connection")
			} else {
				fmt.Printf("error while trying to receive packet: %s", err)
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