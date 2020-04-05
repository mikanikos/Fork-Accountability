package main

import (
	"bytes"
	"github.com/mikanikos/Fork-Accountability/common"
	"github.com/mikanikos/Fork-Accountability/connection"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func createTestMonitor() *Monitor {
	monitorTest := NewMonitor()
	monitorTest.Height = 1
	monitorTest.Timeout = 30
	monitorTest.FirstDecisionRound = 3
	monitorTest.SecondDecisionRound = 4
	monitorTest.Validators = append(monitorTest.Validators, "127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082", "127.0.0.1:8083")
	return monitorTest
}

func TestMonitor_CorrectConfigParsing(t *testing.T) {

	monitorTest := createTestMonitor()

	configFile := "config.yaml"
	monitorConfig, err := parseMonitorConfig(configDirectory+configFile)
	if err != nil {
		t.Fatalf("Monitor exiting: config file not parsed correctly: %s", err)
	}

	monitorConfig.receiveChannel = nil
	monitorTest.receiveChannel = nil

	// compare the two monitors
	if !reflect.DeepEqual(monitorTest, monitorConfig) {
		t.Fatal("Monitor config file was not parsed correctly")
	}
}

func validatorMock(id uint64, address string, delay uint64, hvs *common.HeightVoteSet, t *testing.T) {
	server := connection.NewServer()

	go func() {
		for clientData := range server.ReceiveChannel {

			packet := clientData.Packet

			// if it's a request packet, send the response back
			if packet != nil && packet.Code == connection.HvsRequest {

				// prepare packet
				packet.Code = connection.HvsResponse
				packet.Hvs = hvs

				time.Sleep(time.Duration(delay) * time.Second)

				err := clientData.Connection.Send(packet)
				if err != nil {
					t.Fatalf("Error while sending packet back to monitor: %s", err)
				}
			}
		}
	}()

	err := server.Listen(address)
	if err != nil {
		t.Fatalf("Failed while start listening: %s", err)
	}
}


func TestMonitor_ConnectToValidatorsSuccessfully(t *testing.T) {

	go validatorMock(1, "127.0.0.1:8080", 0, common.NewHeightVoteSet(1), t)
	go validatorMock(2, "127.0.0.1:8081", 0, common.NewHeightVoteSet(2), t)
	go validatorMock(3, "127.0.0.1:8082", 0, common.NewHeightVoteSet(3), t)
	go validatorMock(4, "127.0.0.1:8083", 0, common.NewHeightVoteSet(4), t)

	time.Sleep(time.Second * time.Duration(1))

	testMonitor := createTestMonitor()

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

	go validatorMock(1, "127.0.0.1:8080", 0, common.NewHeightVoteSet(1), t)
	go validatorMock(4, "127.0.0.1:8083", 0, common.NewHeightVoteSet(4), t)

	time.Sleep(time.Second * time.Duration(1))

	testMonitor := createTestMonitor()

	err := testMonitor.connectToValidators()

	if err == nil {
		t.Fatal("Should have failed because one or more validators are not listening")
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestMonitor_RunFailed(t *testing.T) {

	go validatorMock(1, "127.0.0.1:8080", 0, common.NewHeightVoteSet(1), t)
	go validatorMock(2, "127.0.0.1:8081", 0, common.NewHeightVoteSet(2), t)
	go validatorMock(3, "127.0.0.1:8082", 0, common.NewHeightVoteSet(3), t)
	go validatorMock(4, "127.0.0.1:8083", 0, common.NewHeightVoteSet(4), t)

	time.Sleep(time.Second * time.Duration(1))

	testMonitor := createTestMonitor()

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, failStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_RunTimeout(t *testing.T) {

	testMonitor := createTestMonitor()

	testMonitor.Timeout = 3
	delay := testMonitor.Timeout+2

	go validatorMock(1, "127.0.0.1:8080", delay, common.NewHeightVoteSet(1), t)
	go validatorMock(2, "127.0.0.1:8081", delay, common.NewHeightVoteSet(2), t)
	go validatorMock(3, "127.0.0.1:8082", delay, common.NewHeightVoteSet(3), t)
	go validatorMock(4, "127.0.0.1:8083", delay, common.NewHeightVoteSet(4), t)

	time.Sleep(time.Second * time.Duration(1))

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, timeoutStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_RunSuccessful(t *testing.T) {

	go validatorMock(1, "127.0.0.1:8080", 0, common.GetHvsForDefaultConfig1(), t)
	go validatorMock(2, "127.0.0.1:8081", 0, common.GetHvsForDefaultConfig2(), t)
	go validatorMock(3, "127.0.0.1:8082", 0, common.GetHvsForDefaultConfig3(), t)
	go validatorMock(4, "127.0.0.1:8083", 0, common.GetHvsForDefaultConfig4(), t)

	time.Sleep(time.Second * time.Duration(1))

	testMonitor := createTestMonitor()

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, successfulStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}

func TestMonitor_RunSuccessfulWithDelays(t *testing.T) {

	go validatorMock(1, "127.0.0.1:8080", 1, common.GetHvsForDefaultConfig1(), t)
	go validatorMock(2, "127.0.0.1:8081", 4, common.GetHvsForDefaultConfig2(), t)
	go validatorMock(3, "127.0.0.1:8082", 3, common.GetHvsForDefaultConfig3(), t)
	go validatorMock(4, "127.0.0.1:8083", 8, common.GetHvsForDefaultConfig4(), t)

	time.Sleep(time.Second * time.Duration(1))

	testMonitor := createTestMonitor()

	output := captureOutput(testMonitor.Run)
	if !strings.Contains(output, successfulStatus) {
		t.Fatal("Output of the algorithm was not expected")
	}
}
