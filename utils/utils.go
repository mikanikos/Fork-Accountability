package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"runtime"
)

// parse config file given as a parameter and returns the validator data
func ParseConfigFile(localPath string, structure interface{}) error {

	_, validatorFileName, _, ok := runtime.Caller(1)
	if !ok {
		panic("No caller information")
	}

	yamlFile, err := ioutil.ReadFile(path.Dir(validatorFileName) + localPath)
	if err != nil {
		return fmt.Errorf("error while reading file given from input: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, structure)
	if err != nil {
		return fmt.Errorf("error while parsing config yaml file: %s", err)
	}

	return nil
}
