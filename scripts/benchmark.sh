#!/usr/bin/env bash

# number of validators
n_configutations=(4)

# number of rounds
m_configutations=(3)

cd ..

go build ./...

cd cmd/validator
go build

cd ..
cd monitor
go build
cd ../..
cd scripts

benchmark_report=benchmark_report.out

rm $benchmark_report
touch $benchmark_report

echo "Starting benchmarks"

for n in "${n_configutations[@]}"
do
    for m in "${m_configutations[@]}"
    do

        echo "Generating config files for $n validators and $m rounds
" >> $benchmark_report

        python config_generator.py -N $n -M $m
        
        for ((i = 1; i <= n; i++))
        do
            ./../cmd/validator/validator -config="scripts/config_$i.yaml" &
        done

        # give some time to validators to start listening
        sleep 1

        echo "
Async mode
" >> $benchmark_report

        /usr/bin/time -v ./../cmd/monitor/monitor -config="scripts/config.yaml" -asyncMode true >> $benchmark_report 2>&1

        echo "
Sync mode
" >> $benchmark_report


        /usr/bin/time -v ./../cmd/monitor/monitor -config="scripts/config.yaml" -asyncMode false >> $benchmark_report 2>&1
                
        echo "
-------------------------------------------------------------------------------------------------------------------------
" >> $benchmark_report

        pkill -f validator
        pkill -f monitor

        rm *.yaml
    done
done

echo "All benchmarks completed"