package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"

	"gopkg.in/yaml.v2"
)

func main() {

	// parse arguments
	address := flag.String("address", "127.0.0.1:8080", "address where this validator will start listening for requests from the monitor")
	//monitorAddress := flag.String("monitorAddr", "", "monitor address")
	configFile := flag.String("config", "", "configuration file of the validator")

	// parse arguments
	flag.Parse()

	// parse file
	hvs, err := parseConfigFile(*configFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Validator on " + *address + ": started listening for incomiming requests")

	// start listening for incoming connection from monitor
	connection.Listen(*address, hvs)
}

func parseConfigFile(fileName string) (*common.HeightVoteSet, error) {

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error while reading file given from input: %s", err)
	}

	hvs := &common.HeightVoteSet{}

	err = yaml.Unmarshal(yamlFile, hvs)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing config yaml file: %s", err)
	}

	fmt.Println(hvs.String())

	return hvs, nil
}
