#!/usr/bin/env bash

cd ..
go build ./...

RED='\033[0;31m'
NC='\033[0m'
failed="F"
numValidators=4

echo "Starting validators
"

validatorOutFile="validator"
monitorOutFile="monitor"

for i in $(seq 1 $numValidators);
do
  ./cmd/validator/validator -config="cmd/validator/_config/config_" + "$i" + ".yaml" > "./tests/out/$validatorOutFile$i.out" &
done

# give some time to validators to start listening
sleep 2

echo "
Starting monitor
"

./cmd/monitor/monitor -config="cmd/monitor/_config/config.yaml" > "./tests/out/$monitorOutFile.out" &

# give some time for communication and running the algorithm
sleep 2

# Chceck output
msgLine="Algorithm completed"
if ! (grep -q "$msgLine" "$monitorOutFile.out"); then
  failed="T"
fi

pkill -f validator
pkill -f monitor

if [[ "$failed" == "T" ]] ; then
    echo -e "${RED}***FAILED***${NC}"
else
	echo "***PASSED***"
fi