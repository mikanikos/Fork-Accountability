#!/usr/bin/env bash

cd ..
cd monitor
go build
cd ..
cd validator
go build
cd ..

echo "Starting validators
"

./cmd/validator/validator.exe -config="config_1.yaml" &
./cmd/validator/validator.exe -config="config_2.yaml" &
./cmd/validator/validator.exe -config="config_3.yaml" &
./cmd/validator/validator.exe -config="config_4.yaml" &

sleep 2

echo "
Starting monitor
"

./cmd/monitor/monitor.exe -config="config.yaml" &

sleep 3

pkill -f validator
pkill -f monitor