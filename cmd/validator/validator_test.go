package main

import (
	"github.com/mikanikos/Fork-Accountability/connection"
	"github.com/mikanikos/Fork-Accountability/utils"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_CorrectConfigParsing_1(t *testing.T) {

	// create validator structure
	validatorTest := NewValidator()
	validatorTest.ID = 1
	validatorTest.Address = "127.0.0.1:8080"
	validatorTest.Messages[1] = utils.GetHvsForDefaultConfig1()

	// parse config file
	fileName1 := "config_1.yaml"
	validatorConfig, err := parseValidatorConfig(configDirectory + fileName1)
	if err != nil {
		t.Fatalf("Failed while parsing config file 1: %s", err)
	}

	validatorConfig.Server = nil
	validatorTest.Server = nil

	// compare the two validators
	if !reflect.DeepEqual(validatorTest, validatorConfig) {
		t.Fatal("Validator 1 config file was not parsed correctly")
	}
}

func Test_CorrectConfigParsing_2(t *testing.T) {

	// create validator structure
	validatorTest := NewValidator()
	validatorTest.ID = 2
	validatorTest.Address = "127.0.0.1:8081"
	validatorTest.Messages[1] = utils.GetHvsForDefaultConfig2()

	// parse config file
	fileName2 := "config_2.yaml"
	validatorConfig, err := parseValidatorConfig(configDirectory + fileName2)
	if err != nil {
		t.Fatalf("Failed while parsing config file 2: %s", err)
	}

	validatorConfig.Server = nil
	validatorTest.Server = nil

	// compare the two validators
	if !reflect.DeepEqual(validatorTest, validatorConfig) {
		t.Fatal("Validator 2 config file was not parsed correctly")
	}
}

func Test_CorrectConfigParsing_3(t *testing.T) {

	// create validator structure
	validatorTest := NewValidator()
	validatorTest.ID = 3
	validatorTest.Address = "127.0.0.1:8082"
	validatorTest.Messages[1] = utils.GetHvsForDefaultConfig3()

	// parse config file
	fileName3 := "config_3.yaml"
	validatorConfig, err := parseValidatorConfig(configDirectory + fileName3)
	if err != nil {
		t.Fatalf("Failed while parsing config file 3: %s", err)
	}

	validatorConfig.Server = nil
	validatorTest.Server = nil

	// compare the two validators
	if !reflect.DeepEqual(validatorTest, validatorConfig) {
		t.Fatal("Validator 3 config file was not parsed correctly")
	}
}

func Test_CorrectConfigParsing_4(t *testing.T) {

	// create validator structure
	validatorTest := NewValidator()
	validatorTest.ID = 4
	validatorTest.Address = "127.0.0.1:8083"
	validatorTest.Messages[1] = utils.GetHvsForDefaultConfig4()

	// parse config file
	fileName4 := "config_4.yaml"
	validatorConfig, err := parseValidatorConfig(configDirectory + fileName4)
	if err != nil {
		t.Fatalf("Failed while parsing config file 4: %s", err)
	}

	validatorConfig.Server = nil
	validatorTest.Server = nil

	// compare the two validators
	if !reflect.DeepEqual(validatorTest, validatorConfig) {
		t.Fatal("Validator 4 config file was not parsed correctly")
	}
}

func Test_WrongConfigFilename(t *testing.T) {
	fileName3 := "config_not_existing.yaml"
	_, err := parseValidatorConfig(configDirectory + fileName3)
	if err == nil {
		t.Fatalf("Should have failed because filename doesn't exist")
	}
}

func Test_BadFormattedConfig(t *testing.T) {

	badConfig := "bad_config.yaml"
	_ = ioutil.WriteFile(badConfig, []byte("cjdcjdcjd"), 0644)

	_, err := parseValidatorConfig(badConfig)

	_ = os.Remove(badConfig)

	if err == nil {
		t.Fatalf("Should have failed because file is bad formatted")
	}
}

func Test_ValidatorRun(t *testing.T) {

	validatorTest := NewValidator()
	validatorTest.ID = 1
	freeAddress, err := utils.GetFreeAddress()
	if err != nil {
		t.Fatal("Error while getting a free port")
	}

	validatorTest.Address = freeAddress
	validatorCompleted := false

	go func() {
		validatorTest.Run(0)
		validatorCompleted = true
	}()

	time.Sleep(time.Duration(2) * time.Second)

	// client connects
	connClient, err := connection.Connect(validatorTest.Address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %s", err)
	}

	// client sends packet
	err = connClient.Send(&connection.Packet{Code: connection.HvsRequest})
	if err != nil {
		t.Fatalf("Failed to send packet on client: %s", err)
	}

	time.Sleep(time.Duration(2) * time.Second)

	if validatorCompleted {
		t.Fatal("Validator exited unexpectedly")
	}
}
