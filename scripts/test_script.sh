#!/usr/bin/env bash

cd ..

go build ./...

cd cmd/validator
go build

cd ..
cd monitor
go build
cd ../..

numValidators=4

echo "Starting validators
"

# for i in $(seq 1 $numValidators);
# do
#   ./cmd/validator/validator -config="/cmd/validator/_config/config_$i.yaml" &
# done

./cmd/validator/validator -config="/cmd/validator/_config/config_1.yaml" -delay=1 &
./cmd/validator/validator -config="/cmd/validator/_config/config_2.yaml" -delay=5 &
./cmd/validator/validator -config="/cmd/validator/_config/config_3.yaml" -delay=1 &
./cmd/validator/validator -config="/cmd/validator/_config/config_4.yaml" -delay=5 &

# give some time to validators to start listening
sleep 1

echo "
Starting monitor
"

./cmd/monitor/monitor -config="/cmd/monitor/_config/config.yaml" -report="true" &

# give some time for communication and running the algorithm
sleep 15

pkill -f validator
pkill -f monitor