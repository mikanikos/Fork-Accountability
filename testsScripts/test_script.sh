#!/usr/bin/env bash

cd ..
go build
cd validator
go build
cd ..

echo "Starting validators
"

./validator/validator -config="validator/config_1.yaml" -address="127.0.0.1:8080" &
./validator/validator -config="validator/config_2.yaml" -address="127.0.0.1:8081" &
./validator/validator -config="validator/config_3.yaml" -address="127.0.0.1:8082" &
./validator/validator -config="validator/config_4.yaml" -address="127.0.0.1:8083" &

sleep 3

echo "
Starting monitor
"

./Fork-Accountability -processes="127.0.0.1:8080,127.0.0.1:8081,127.0.0.1:8082,127.0.0.1:8083" -firstDecisionRound=3 -secondDecisionRound=4 -waitTimeout=5

sleep 3

pkill -f validator
pkill -f Fork-Accountability