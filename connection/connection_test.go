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
		err := Listen("127.0.0.1:7070", nil)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
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

	// server
	go func() {
		err := Listen("127.0.0.1:7070", nil)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()

	time.Sleep(time.Second * time.Duration(3))

	connClient, err := Connect("127.0.0.1:7070")
	if err != nil {
		t.Fatalf("Failed to connect to server: %s", err)
	}

	err = Send(connClient, &Packet{Code: HvsRequest})
	if err != nil {
		t.Fatalf("Failed to send packet: %s", err)
	}

	packet, err := Receive(connClient)

	if err != nil {
		t.Fatalf("Failed to receive packet: %s", err)
	}

	if packet.Code != HvsResponse {
		t.Fatal("Failed to send/receive correct packet")
	}
}
