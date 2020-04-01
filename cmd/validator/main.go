package main

import (
	"flag"
	"github.com/mikanikos/Fork-Accountability/utils"
	"log"
)

const configDirectory = "/_config/"

func main() {

	// parse arguments
	configFile := flag.String("config", "", "configuration file path of the validator")

	// parse arguments
	flag.Parse()

	// parse file
	validator := NewValidator()
	err := utils.ParseConfigFile(configDirectory+*configFile, validator)
	if err != nil {
		log.Fatalf("Validator exiting: config file not parsed correctly: %s", err)
	}

	// start validator execution
	validator.Run()
}
