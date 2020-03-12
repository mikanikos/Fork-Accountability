package connection

import (
	"net"
	"testing"
	"time"
)

func Test_basic_server_client_interaction(t *testing.T) {

	// server
	go func() {
		err := Listen("127.0.0.1:7070", nil)
		if err != nil {
			t.Fatalf("Failed while start listening: %s", err)
		}
	}()

	time.Sleep(time.Second * time.Duration(3))

	connClient, err := net.Dial("tcp", "127.0.0.1:7070")

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
