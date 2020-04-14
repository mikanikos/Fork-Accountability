package main

import (
	"flag"
	"log"

	"github.com/mikanikos/Fork-Accountability/utils"
)

func main() {

	// parse arguments
	configFile := flag.String("config", configPath, "path (relative to the project root directory) of the configuration file for the monitor")
	writeReport := flag.String("report", "", "path (relative to the project root directory) of the report to generate at the end of the execution instead of printing logs to standard output")

	// parse arguments
	flag.Parse()

	// parse file
	monitor, err := parseMonitorConfig(*configFile)
	if err != nil {
		log.Fatalf("Monitor exiting: config file not parsed correctly: %s", err)
	}

	// start monitor execution
	monitor.Run(*writeReport)
}

// parse config file for the monitor
func parseMonitorConfig(configFile string) (*Monitor, error) {
	monitor := NewMonitor()
	err := utils.ParseConfigFile(configFile, monitor)
	return monitor, err
}
