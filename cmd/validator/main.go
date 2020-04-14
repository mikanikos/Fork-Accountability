package main

import (
	"flag"
	"log"

	"github.com/mikanikos/Fork-Accountability/utils"
)

const configDirectory = "/cmd/validator/_config/"

func main() {

	// parse arguments
	configFile := flag.String("config", configDirectory+"config_1.yaml", "relative path of the configuration file for the validator respect to the project folder")
	delay := flag.Uint64("delay", 0, "time to wait (in seconds) before replying back to the monitor, use for testing")

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
