#!/usr/bin/env bash

cd ..
go build ./...


echo "Starting validators
"

./cmd/validator/validator -config="config_1.yaml" &
./cmd/validator/validator -config="config_2.yaml" &
./cmd/validator/validator -config="config_3.yaml" &
./cmd/validator/validator -config="config_4.yaml" &

sleep 2

echo "
Starting monitor
"

./cmd/monitor/monitor -config="config.yaml" &

sleep 3

pkill -f validator
pkill -f monitor