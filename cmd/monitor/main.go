package main

import (
	"flag"
	"log"

	"github.com/mikanikos/Fork-Accountability/utils"
)

func main() {

	// parse arguments
	configFile := flag.String("config", configRelativePath+"config.yaml", "relative path of the configuration file for the monitor respect to the project folder")
	writeReport := flag.Bool("report", false, "specify if a report should be generated at the end of the execution instead of printing to standard output")

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
