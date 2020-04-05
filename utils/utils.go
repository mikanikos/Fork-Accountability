package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"path"
	"runtime"
)

// parse config file given as a parameter and returns the validator data
func ParseConfigFile(localPath string, structure interface{}) error {

	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	yamlFile, err := ioutil.ReadFile(path.Dir(path.Dir(filePath)) + localPath)
	if err != nil {
		return fmt.Errorf("error while reading file given from input: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, structure)
	if err != nil {
		return fmt.Errorf("error while parsing config yaml file: %s", err)
	}

	return nil
}

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
