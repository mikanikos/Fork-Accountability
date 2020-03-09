package validator

import (
	"flag"

	"github.com/mikanikos/Fork-Accountability/connection"
)

func main() {

	// parse arguments
	port := flag.String("port", "8080", "port of this validator")
	monitorAddress := flag.String("monitorAddr", "", "monitor address")
	configFile := flag.String("config", "", "configuration file")

	// parse arguments
	flag.Parse()

	// parse file
	parseConfigFile(*configFile)

	// setup hvs

	// start listening for incoming connection from monitor
	connection.Listen(":"+*port, *monitorAddress)
}

func parseConfigFile(fileName string) {
	// TODO
}
