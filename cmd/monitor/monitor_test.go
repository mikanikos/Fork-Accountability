package main

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
	"github.com/mikanikos/Fork-Accountability/utils"
)

func createTestMonitor() *Monitor {
	monitorTest := NewMonitor()
	monitorTest.Height = 1
	monitorTest.Timeout = 30
	monitorTest.FirstDecisionRound = 3
	monitorTest.SecondDecisionRound = 4

	addresses, err := utils.GetFreeAddresses(4)
	if err != nil {
		return nil
	}

	monitorTest.Validators = append(monitorTest.Validators, addresses...)
	return monitorTest
}

func TestMonitor_CorrectConfigParsing(t *testing.T) {

	monitorTest := createTestMonitor()

	configFile := "config.yaml"
	monitorConfig, err := parseMonitorConfig(configRelativePath + configFile)
	if err != nil {
		t.Fatalf("Monitor exiting: config file not parsed correctly: %s", err)
	}

	monitorConfig.receiveChannel = nil
	monitorConfig.accAlgorithm = nil

	monitorTest.receiveChannel = nil
	monitorTest.accAlgorithm = nil

	monitorTest.Validators = []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082", "127.0.0.1:8083"}

	// compare the two monitors
	if !reflect.DeepEqual(monitorTest, monitorConfig) {
		t.Fatal("Monitor config file was not parsed correctly")
	}
}

func validatorMock(address string, delay uint64, hvs *common.HeightVoteSet) {
	server := connection.NewServer()

	go func() {
		for clientData := range server.ReceiveChannel {

			time.Sleep(time.Duration(delay) * time.Second)

			packet := clientData.Packet

			// if it's a request packet, send the response back
			if packet != nil && packet.Code == connection.HvsRequest {

				// prepare packet
				packet.Code = connection.HvsResponse
				packet.Hvs = hvs

				err := clientData.Connection.Send(packet)
				if err != nil {
					log.Printf("Error while sending packet back to monitor: %s", err)
				}
			}
		}
	}()

	err := server.Listen(address)
	if err != nil {
		log.Printf("Failed while start listening: %s", err)
	}
}

func TestMonitor_ConnectToValidatorsSuccessfully(t *testing.T) {

	testMonitor := createTestMonitor()

	go validatorMock(testMonitor.Validators[0], 0, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[1], 0, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[2], 0, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[3], 0, common.NewHeightVoteSet())

	time.Sleep(time.Second * time.Duration(1))

	err := testMonitor.connectToValidators()

	if err != nil {
		t.Fatal("Failed to connect to validators")
	}
}

func TestMonitor_ConnectToValidatorsNoValidatorsGiven(t *testing.T) {

	time.Sleep(time.Second * time.Duration(1))

	testMonitor := createTestMonitor()
	testMonitor.Validators = nil

	err := testMonitor.connectToValidators()

	if err == nil {
		t.Fatal("Should have failed because one or more validators are not listening")
	}
}

func TestMonitor_ConnectToValidators_Fail(t *testing.T) {

	testMonitor := createTestMonitor()

	go validatorMock(testMonitor.Validators[0], 0, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[3], 0, common.NewHeightVoteSet())

	time.Sleep(time.Second * time.Duration(1))

	err := testMonitor.connectToValidators()

	if err == nil {
		t.Fatal("Should have failed because one or more validators are not listening")
	}
}

func captureOutput(f func(bool)) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f(false)
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestMonitor_RunFailed(t *testing.T) {

	testMonitor := createTestMonitor()

	go validatorMock(testMonitor.Validators[0], 0, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[1], 0, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[2], 0, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[3], 0, common.NewHeightVoteSet())

	time.Sleep(time.Second * time.Duration(1))

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, failStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_RunTimeout(t *testing.T) {

	testMonitor := createTestMonitor()

	testMonitor.Timeout = 3
	delay := testMonitor.Timeout + 2

	go validatorMock(testMonitor.Validators[0], delay, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[1], delay, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[2], delay, common.NewHeightVoteSet())
	go validatorMock(testMonitor.Validators[3], delay, common.NewHeightVoteSet())

	time.Sleep(time.Second * time.Duration(1))

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, timeoutStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_RunSuccessful(t *testing.T) {

	testMonitor := createTestMonitor()

	go validatorMock(testMonitor.Validators[0], 0, utils.GetHvsForDefaultConfig1())
	go validatorMock(testMonitor.Validators[1], 0, utils.GetHvsForDefaultConfig2())
	go validatorMock(testMonitor.Validators[2], 0, utils.GetHvsForDefaultConfig3())
	go validatorMock(testMonitor.Validators[3], 0, utils.GetHvsForDefaultConfig4())

	time.Sleep(time.Second * time.Duration(1))

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, successfulStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_RunSuccessfulWithDelays(t *testing.T) {

	testMonitor := createTestMonitor()

	go validatorMock(testMonitor.Validators[0], 1, utils.GetHvsForDefaultConfig1())
	go validatorMock(testMonitor.Validators[1], 4, utils.GetHvsForDefaultConfig2())
	go validatorMock(testMonitor.Validators[2], 3, utils.GetHvsForDefaultConfig3())
	go validatorMock(testMonitor.Validators[3], 6, utils.GetHvsForDefaultConfig4())

	time.Sleep(time.Second * time.Duration(2))

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, successfulStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_RunSuccessfulWithAllFaultyFirst(t *testing.T) {

	testMonitor := createTestMonitor()

	go validatorMock(testMonitor.Validators[0], 3, utils.GetHvsForDefaultConfig1())
	go validatorMock(testMonitor.Validators[1], 3, utils.GetHvsForDefaultConfig2())
	go validatorMock(testMonitor.Validators[2], 1, utils.GetHvsForDefaultConfig3())
	go validatorMock(testMonitor.Validators[3], 1, utils.GetHvsForDefaultConfig4())

	time.Sleep(time.Second * time.Duration(1))

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, successfulStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_RunSuccessfulWithAllFaultyLast(t *testing.T) {

	testMonitor := createTestMonitor()

	go validatorMock(testMonitor.Validators[0], 1, utils.GetHvsForDefaultConfig1())
	go validatorMock(testMonitor.Validators[1], 1, utils.GetHvsForDefaultConfig2())
	go validatorMock(testMonitor.Validators[2], 3, utils.GetHvsForDefaultConfig3())
	go validatorMock(testMonitor.Validators[3], 3, utils.GetHvsForDefaultConfig4())

	time.Sleep(time.Second * time.Duration(1))

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, successfulStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_WriteReport(t *testing.T) {

	testMonitor := createTestMonitor()

	go validatorMock(testMonitor.Validators[0], 1, utils.GetHvsForDefaultConfig1())
	go validatorMock(testMonitor.Validators[1], 4, utils.GetHvsForDefaultConfig2())
	go validatorMock(testMonitor.Validators[2], 3, utils.GetHvsForDefaultConfig3())
	go validatorMock(testMonitor.Validators[3], 6, utils.GetHvsForDefaultConfig4())

	time.Sleep(time.Second * time.Duration(2))

	_ = os.Remove(reportDirectory+reportFile)

	testMonitor.Run(true)

	if _, err := os.Stat(reportDirectory+reportFile); os.IsNotExist(err) {
		t.Fatal("Monitor didn't generate report")
	}
}