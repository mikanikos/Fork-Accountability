package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"runtime"
	"strconv"

	"gopkg.in/yaml.v2"
)

const localhost = "127.0.0.1:"

// ParseConfigFile parses the config file path from the project root directory and returns the validator data
func ParseConfigFile(localPath string, structure interface{}) error {

	projectPath, err := getProjectRootPath()
	if err != nil {
		return fmt.Errorf("error getting project root path: %s", err)
	}

	yamlFile, err := ioutil.ReadFile(path.Join(projectPath, localPath))
	if err != nil {
		return fmt.Errorf("error while reading file given from input: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, structure)
	if err != nil {
		return fmt.Errorf("error while parsing config yaml file: %s", err)
	}

	return nil
}

// GetFreeAddress asks the kernel for a free open port that is ready to use
func GetFreeAddress() (string, error) {

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

// GetFreeAddresses asks the kernel for free open ports that are ready to use
func GetFreeAddresses(count int) ([]string, error) {

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

// OpenFile open a file given its local path from the project root
func OpenFile(localPath string) (*os.File, error) {

	projectPath, err := getProjectRootPath()
	if err != nil {
		return nil, fmt.Errorf("error getting project root path: %s", err)
	}

	return os.OpenFile(path.Join(projectPath, localPath), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

func getProjectRootPath() (string, error) {

	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("No caller information")
	}

	return path.Dir(path.Dir(filePath)), nil
}
