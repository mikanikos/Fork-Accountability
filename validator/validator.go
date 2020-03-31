package main

import (
	"flag"
	"fmt"
	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/utils"
	"log"
)

const configDirectory = "/_config/"

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
	err := utils.ParseConfigFile(configDirectory + *configFile, validator)
	if err != nil {
		log.Fatalf("Validator exiting: config file not parsed correctly: %s", err)
	}

	fmt.Println("Validator on " + validator.Address + ": start listening for incoming requests")

	// start listening for incoming connection from monitor
	err = validator.Listen()
	if err != nil {
		log.Fatalf("Validator %s exiting: cannot listen on given address: %s", validator.Address, err)
	}
}
