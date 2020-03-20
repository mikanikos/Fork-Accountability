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

./validator/validator -config="config_1.yaml" -address="127.0.0.1:21990" &
./validator/validator -config="config_2.yaml" -address="127.0.0.1:21991" &
./validator/validator -config="config_3.yaml" -address="127.0.0.1:21992" &
./validator/validator -config="config_4.yaml" -address="127.0.0.1:21993" &

sleep 2

echo "
Starting monitor
"

./monitor/monitor -processes="127.0.0.1:21990,127.0.0.1:21991,127.0.0.1:21992,127.0.0.1:21993" -firstDecisionRound=3 -secondDecisionRound=4 -waitTimeout=5 &

sleep 7

pkill -f validator
pkill -f monitor