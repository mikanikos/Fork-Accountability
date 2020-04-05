package main

import (
	"flag"
	"log"

	"github.com/mikanikos/Fork-Accountability/utils"
)

const configDirectory = "/cmd/validator/_config/"

func main() {

	// parse arguments
	configFile := flag.String("config", configDirectory+"config_1.yaml", "configuration file path of the validator")
	delay := flag.Uint64("delay", 0, "time to wait (in seconds) before replying back to the monitor, used for testing")

	// parse arguments
	flag.Parse()

	// parse file
	validator, err := parseValidatorConfig(*configFile)
	if err != nil {
		log.Fatalf("Validator exiting: config file not parsed correctly: %s", err)
	}

	// start validator execution
	validator.Run(*delay)
}

// parse config file for the validator
func parseValidatorConfig(configFile string) (*Validator, error) {
	validator := NewValidator()
	err := utils.ParseConfigFile(configFile, validator)
	return validator, err
}
