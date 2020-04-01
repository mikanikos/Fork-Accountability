package connection

import (
	"testing"
	"time"
)

// run tests individually because of persistent connections between tests

func Test_ServerInitialization(t *testing.T) {

	// server
	go func() {
		err := Listen("127.0.0.1:7070", nil)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()
	time.Sleep(time.Duration(2) * time.Second)
}

func Test_ServerWrongAddressForListen(t *testing.T) {

	// server
	go func() {
		err := Listen("", nil)
		if err == nil {
			t.Fatalf("Should have failed listening: %s", err)
		}
	}()

	time.Sleep(time.Duration(2) * time.Second)
}

func Test_ClientInitialization(t *testing.T) {

	// server
	go func() {
		err := Listen("127.0.0.1:7070", nil)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()

	_, err := Connect("127.0.0.1:7070")
	if err != nil {
		t.Fatalf("Failed to connect to server: %s", err)
	}
}

func Test_ClientFailingToConnect(t *testing.T) {

	_, err := Connect("127.0.0.1:7070")
	if err == nil {
		t.Fatal("Connection should have not been successful")
	}
}

func Test_ClientSendsMessage(t *testing.T) {

	// server
	go func() {
		err := Listen("127.0.0.1:7070", nil)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()

	connClient, err := Connect("127.0.0.1:7070")
	if err != nil {
		t.Fatalf("Failed to connect to server: %s", err)
	}

	err = Send(connClient, &Packet{Code: HvsRequest})
	if err != nil {
		t.Fatalf("Failed sending message: %s", err)
	}
}

func Test_ServerClientInteraction(t *testing.T) {

	receiveChannel := make(chan *PacketConnection)

	// server start listening
	go func() {
		err := Listen("127.0.0.1:7070", receiveChannel)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()

	time.Sleep(time.Second * time.Duration(3))

	// client connects
	connClient, err := Connect("127.0.0.1:7070")
	if err != nil {
		t.Fatalf("Failed to connect to server: %s", err)
	}

	// client sends packet
	err = Send(connClient, &Packet{Code: HvsRequest})
	if err != nil {
		t.Fatalf("Failed to send packet on client: %s", err)
	}

	// server receives packet
	packetFromClient := <-receiveChannel
	packetFromClient.Packet.Code = HvsResponse

	// server sends packet back with modified code
	err = Send(packetFromClient.Conn, packetFromClient.Packet)
	if err != nil {
		t.Fatalf("Failed to send packet on server: %s", err)
	}

	// client receives packet
	packet, err := Receive(connClient)
	if err != nil {
		t.Fatalf("Failed to receive packet: %s", err)
	}

	// check if packet is the one expected
	if packet.Code != HvsResponse {
		t.Fatal("Failed to send/receive correct packet")
	}
}
