package main

const (
	configPath = "cmd/monitor/_config/config.yaml"
	reportPath = "cmd/monitor/_report/report.out"

	sendTimer = 5

	maxChannelSize = 100

	successfulStatus = "Monitor: Algorithm completed"
	failStatus       = "Monitor: Algorithm failed because not enough message logs have been received or the message logs received were not sufficient to find at least f+1 faulty processes"
	timeoutStatus    = "Monitor: Algorithm failed because of timeout expiration"
)
