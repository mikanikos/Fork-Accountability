package main

const (
	configRelativePath = "/cmd/monitor/_config/"
	reportDirectory    = "/_report/"
	reportFile         = "report.out"

	sendTimer = 3

	successfulStatus = "Monitor: Algorithm completed"
	failStatus       = "Monitor: Algorithm failed because not enough message logs have been received or the message logs received were not sufficient to find at least f+1 faulty processes"
	timeoutStatus    = "Monitor: Algorithm failed because of timeout expiration"
)
