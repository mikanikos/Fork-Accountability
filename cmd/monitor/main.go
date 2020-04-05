package main

import (
	"flag"
	"github.com/mikanikos/Fork-Accountability/utils"
	"log"
)

const configDirectory = "/cmd/monitor/_config/"

func main() {

	// parse arguments
	configFile := flag.String("config", configDirectory +"config.yaml", "configuration file path of the monitor")

	// parse arguments
	flag.Parse()

	// parse file
	monitor, err := parseMonitorConfig(*configFile)
	if err != nil {
		log.Fatalf("Monitor exiting: config file not parsed correctly: %s", err)
	}

	// start monitor execution
	monitor.Run()
}

// parse config file for the monitor
func parseMonitorConfig(configFile string) (*Monitor, error) {
	monitor := NewMonitor()
	err := utils.ParseConfigFile(configFile, monitor)
	return monitor, err
}
