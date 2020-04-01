package connection

import (
	"testing"
	"time"
)

// run tests individually because of persistent connections between tests

func Test_ServerInitialization(t *testing.T) {

	// server
	go func() {
		err := NewServer().Listen("127.0.0.1:7070")
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()
	time.Sleep(time.Duration(2) * time.Second)
}

func Test_ServerWrongAddressForListen(t *testing.T) {

	// server
	go func() {
		err := NewServer().Listen("")
		if err == nil {
			t.Fatalf("Should have failed listening: %s", err)
		}
	}()

	time.Sleep(time.Duration(2) * time.Second)
}

func Test_ClientInitialization(t *testing.T) {

	address := "127.0.0.1:7071"

	// server
	go func() {
		err := NewServer().Listen(address)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()

	time.Sleep(time.Duration(1) * time.Second)

	_, err := Connect(address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %s", err)
	}
}

func Test_ClientFailingToConnect(t *testing.T) {

	_, err := Connect("")
	if err == nil {
		t.Fatal("Connection should have not been successful")
	}
}

func Test_ClientSendsMessage(t *testing.T) {

	address := "127.0.0.1:9090"

	// server
	go func() {
		err := NewServer().Listen(address)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()

	time.Sleep(time.Duration(2) * time.Second)

	connClient, err := Connect(address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %s", err)
	}

	err = connClient.Send(&Packet{Code: HvsRequest})
	if err != nil {
		t.Fatalf("Failed sending message: %s", err)
	}
}

func Test_ServerClientInteraction(t *testing.T) {

	address := "127.0.0.1:6060"

	server := NewServer()

	// server start listening
	go func() {
		err := server.Listen(address)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()

	time.Sleep(time.Second * time.Duration(3))

	// client connects
	connClient, err := Connect(address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %s", err)
	}

	// client sends packet
	err = connClient.Send(&Packet{Code: HvsRequest})
	if err != nil {
		t.Fatalf("Failed to send packet on client: %s", err)
	}

	// server receives packet
	packetFromClient := <-server.ReceiveChannel
	packetFromClient.Packet.Code = HvsResponse

	// server sends packet back with modified code
	err = packetFromClient.Connection.Send(packetFromClient.Packet)
	if err != nil {
		t.Fatalf("Failed to send packet on server: %s", err)
	}

	// client receives packet
	packet, err := connClient.Receive()
	if err != nil {
		t.Fatalf("Failed to receive packet: %s", err)
	}

	// check if packet is the one expected
	if packet.Code != HvsResponse {
		t.Fatal("Failed to send/receive correct packet")
	}
}
