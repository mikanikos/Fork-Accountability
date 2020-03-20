#!/usr/bin/env bash

cd ..
go build
cd validator
go build
cd ..

echo "Starting validators
"

./validator/validator -config="config_1.yaml" -address="127.0.0.1:8080" &
./validator/validator -config="config_2.yaml" -address="127.0.0.1:8081" &
./validator/validator -config="config_3.yaml" -address="127.0.0.1:8082" &
./validator/validator -config="config_4.yaml" -address="127.0.0.1:8083" &

sleep 2

echo "
Starting monitor
"

./Fork-Accountability -processes="127.0.0.1:8080,127.0.0.1:8081,127.0.0.1:8082,127.0.0.1:8083" -firstDecisionRound=3 -secondDecisionRound=4 -waitTimeout=5

sleep 1

pkill -f validator
pkill -f Fork-Accountability