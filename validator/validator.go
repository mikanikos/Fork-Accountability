package main

import (
	"flag"
	"fmt"
	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
	"github.com/mikanikos/Fork-Accountability/utils"
	"log"
)

const (
	configDirectory = "/_config/"
	maxChannelSize  = 100
)

// Validator struct
type Validator struct {
	ID       uint64                           `yaml:"id"`
	Address  string                           `yaml:"address"`
	Messages map[uint64]*common.HeightVoteSet `yaml:"messages"`
}

func main() {

	// parse arguments
	configFile := flag.String("config", "", "configuration file path of the validator")

	// parse arguments
	flag.Parse()

	// parse file
	validator := &Validator{}
	err := utils.ParseConfigFile(configDirectory+*configFile, validator)
	if err != nil {
		log.Fatalf("Validator exiting: config file not parsed correctly: %s", err)
	}

	fmt.Println("Validator on " + validator.Address + ": start listening for incoming requests")

	receiveChannel := make(chan *connection.PacketConnection, maxChannelSize)
	go validator.processPacketFromMonitor(receiveChannel)

	// start listening for incoming connection from monitor
	err = connection.Listen(validator.Address, receiveChannel)
	if err != nil {
		log.Fatalf("Validator %s exiting: cannot listen on given address: %s", validator.Address, err)
	}
}

// process packet from client (monitor)
func (validator *Validator) processPacketFromMonitor(receiveChannel chan *connection.PacketConnection) {
	for packetConn := range receiveChannel {
		packet := packetConn.Packet
		conn := packetConn.Conn
		// if it's a request packet, send the response back
		if packet.Code == connection.HvsRequest {
			fmt.Println("Validator on " + validator.Address + ": sending hvs to monitor")
			packet.Code = connection.HvsResponse
			// get the messages in the requested height
			packet.Hvs = validator.Messages[packet.Height]
			err := connection.Send(conn, packet)
			if err != nil {
				fmt.Println("Error while sending packet back to monitor")
			}
		}
	}
}
