package main

import (
	"flag"
	"log"
	"time"

	"github.com/mikanikos/Fork-Accountability/utils"
)

func main() {

	// parse arguments
	configFile := flag.String("config", configPath, "path (relative to the project root directory) of the configuration file for the monitor")
	report := flag.String("report", "", "path (relative to the project root directory) of the report to generate at the end of the execution instead of printing logs to standard output")
	asyncMode := flag.Bool("asyncMode", true, "run the accountability algorithm asynchronously")
	delay := flag.Uint64("delay", 0, "time to wait (in seconds) before start running, use for testing")

	// parse arguments
	flag.Parse()

	// parse file
	monitor, err := newMonitorFromConfig(*configFile)
	if err != nil {
		log.Fatalf("Monitor exiting: config file not parsed correctly: %s", err)
	}

	time.Sleep(time.Duration(*delay) * time.Second)

	// start monitor execution
	monitor.Run(*report, *asyncMode)
}

// create a new monitor from config file
func newMonitorFromConfig(configFile string) (*Monitor, error) {
	monitor := NewMonitor()
	err := utils.ParseConfigFile(configFile, monitor)
	return monitor, err
}
