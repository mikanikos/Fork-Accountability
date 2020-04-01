package main

import (
	"flag"
	"github.com/mikanikos/Fork-Accountability/utils"
	"log"
)

const configDirectory = "/_config/"

func main() {

	// parse arguments
	configFile := flag.String("config", "", "configuration file path of the monitor")

	// parse arguments
	flag.Parse()

	// parse file
	monitor := &Monitor{}
	err := utils.ParseConfigFile(configDirectory+*configFile, monitor)
	if err != nil {
		log.Fatalf("Monitor exiting: config file not parsed correctly: %s", err)
	}

	// start monitor execution
	monitor.Run()
}
