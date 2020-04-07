#!/usr/bin/env bash

cd ..
cd cmd/validator
go build

cd ..
cd monitor
go build
cd ../..

numValidators=4

echo "Starting validators
"

for i in $(seq 1 $numValidators);
do
  ./cmd/validator/validator -config="/cmd/validator/_config/config_$i.yaml" &
done

# give some time to validators to start listening
sleep 2

echo "
Starting monitor
"

./cmd/monitor/monitor -config="/cmd/monitor/_config/config.yaml" &

# give some time for communication and running the algorithm
sleep 5

pkill -f validator
pkill -f monitor