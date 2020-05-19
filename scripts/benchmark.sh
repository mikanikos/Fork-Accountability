#!/usr/bin/env bash

configutations=(15)

cd ..

go build ./...

cd cmd/validator
go build

cd ..
cd monitor
go build
cd ../..
cd scripts

for n in "${configutations[@]}"
do

    echo "Configurations generated
    "

    python config_generator.py -N $n -M 2

    echo "Starting validators
    "
    
    for ((i = 1; i <= n; i++))
    do
        ./../cmd/validator/validator -config="scripts/config_$i.yaml" &
    done

    # give some time to validators to start listening
    sleep 1

    echo "
    Starting monitor
    "

    /usr/bin/time -v ./../cmd/monitor/monitor -config="scripts/config.yaml"

    pkill -f validator
    pkill -f monitor

    rm *.yaml
done