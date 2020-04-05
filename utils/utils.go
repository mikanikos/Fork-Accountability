package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"path"
	"runtime"
	"strconv"
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

// GetFreeAddress asks the kernel for a free open port that is ready to use.
func GetFreeAddress() (string, error) {

	localhost := "127.0.0.1:"

	addr, err := net.ResolveTCPAddr("tcp", localhost+"0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return localhost + strconv.Itoa(l.Addr().(*net.TCPAddr).Port), nil
}

// GetFreeAddress asks the kernel for free open ports that are ready to use.
func GetFreeAddresses(count int) ([]string, error) {

	localhost := "127.0.0.1:"
	ports := make([]string, 0)

	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", localhost+"0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, localhost+strconv.Itoa(l.Addr().(*net.TCPAddr).Port))
	}
	return ports, nil
}
