package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"

	"gopkg.in/yaml.v2"
)

const configDirectory = "/_config/"

func main() {

	// parse arguments
	address := flag.String("address", "127.0.0.1:8080", "address where this validator will start listening for requests from the monitor")
	//monitorAddress := flag.String("monitorAddr", "", "monitor address")
	configFile := flag.String("config", "", "configuration file path of the validator")

	// parse arguments
	flag.Parse()

	// parse file
	hvs, err := parseConfigFile(*configFile)
	if err != nil {
		fmt.Printf("Validator %s exiting: config file not parsed correctly: %s", *address, err)
		os.Exit(1)
	}

	fmt.Println("Validator on " + *address + ": start listening for incoming requests")

	// start listening for incoming connection from monitor
	err = connection.Listen(*address, hvs)
	if err != nil {
		fmt.Printf("Validator %s exiting: cannot listen on given address: %s", *address, err)
		os.Exit(1)
	}
}

// parse config file given as a parameter and returns the hvs
func parseConfigFile(fileName string) (*common.HeightVoteSet, error) {

	_, validatorFileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	yamlFile, err := ioutil.ReadFile(path.Dir(validatorFileName) + configDirectory + fileName)
	if err != nil {
		return nil, fmt.Errorf("error while reading file given from input: %s", err)
	}

	hvs := &common.HeightVoteSet{}

	err = yaml.Unmarshal(yamlFile, hvs)
	if err != nil {
		return nil, fmt.Errorf("error while parsing config yaml file: %s", err)
	}

	return hvs, nil
}
