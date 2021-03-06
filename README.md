# Fork-Accountability

![](https://github.com/mikanikos/Fork-Accountability/workflows/Build%20and%20Tests/badge.svg)
[![codecov](https://codecov.io/gh/mikanikos/Fork-Accountability/branch/master/graph/badge.svg)](https://codecov.io/gh/mikanikos/Fork-Accountability)

The scope of this project is to implement a simple PoC of a Fork Accountability algorithm in the Tendermint Consensus protocol in the Go language.

The repository includes simple libraries and modules that allow to run experiments and benchmark the algorithm implementation in a minimal test environment:

- a simple connection library to provide an high-level API for establishing connections and communicating among processes based on a client-server architecture; 

- an accountability algorithm implementation based on the theoretical specifications given by the research scientists of the DCL lab at [EPFL](https://www.epfl.ch/en/) and [Informal Systems](https://informal.systems/);  

- monitor implementation that represents the independent verification entity used to run the accountability algorithm;  

- validator implementation that represents the processes that participate in the Tendermint Consensus protocol; 

- sample test scripts to easily run experiments and test both the accountability algorithm and the interaction between monitor and validators.

## Context and overview
Documentation describing the context and the theoretical background can be found in the [docs](docs) folder.   

## Architecture

The project is organized in several files and packages in order to guarantee modularity and extensibility and maintain at the same time a clear and simple structure.

The monitor ([monitor package](cmd/monitor)) is the entity responsible to run the accountability algorithm. It takes as input parameters:

- a path of a yaml configuration file which are necessary to initialize and configure the algorithm

- an (optional) path of a file that can be generated to provide detailed information about the whole execution and, especially, of the accountability algorithm

The monitor is responsible for opening connections with all the validators and initialize the request of the message logs (described below).

Validators ([validator package](cmd/validator)) are simple processes that listen on a given port and each one has its own messages logs that are the result of the execution of the Tendermint consensus protocol. Message logs and listening port are initialized through a configuration file (different for every validator) along with a unique validator identifier.
Every message in the configuration file must be specified with all the corresponding information associated with it (type, round, height, senderId, possible justifications).

The monitor will use the connection library to request message logs from all the addresses (i.e., validator processes) given in the config file. It will wait for responses from each validator and, as soon as a packet arrives, it will store it and send it to the main thread. The main thread will run the fork accountability algorithm if enough messages have been received until that time.
The monitor will repeat the request after a timeout expires and if the message received is not valid. If the validator closes a connection or crashes, the monitor will stop waiting for packets from it and will notify the main thread about the failure in the reception.

The validator, after receiving a valid request packet, will response back if it will have the message logs requested. Otherwise, it will just ignore the request and will not answer the monitor. Optionally, it's possible to configure a response to immediately inform the monitor about the missing log in order to save resources.

The main accountability algorithm is implemented in the accountability package and is described in details in documentation files of the docs folders. Please refer to for a theoretical background or for implementation-specific details.

The connection library implemented in this project wraps the well-known [net library](https://golang.org/pkg/net/) and provides some abstractions to establish a TCP connection, send and receive TCP packets, serialize and de-serialize messages and listen to a specific port.
This library is used by the monitor and the validator to exchange packets for both the request and the sending of the message logs.

## Structure

As an overview, this is the current structure of the project:
    
- [.github](.github): contains Github Actions continuous integration config files

- [accountability](accountability): contains the main accountability algorithm

- [cmd](cmd): contains the binaries for the monitor and the validator. Inside each binary folder, there's a folder with sample config files. 

- [common](common): contains abstractions used throughout the project to better handle the input of the algorithm;

- [connection](connection): contains the connection library used by monitor and validators to communicate.

- [docs](docs): contains markdown files documenting the project and the accountability algorithm from a slightly more theoretical perspective; 

- [scripts](scripts): folder used to group scripts for running experiments in different scenarios; 

- [utils](utils): utilities used for parsing configuration files and for testing the several functionalities of the modules implemented;

Each package contains tests in `*_test.go` files.

## How to run experiments

### Prerequisites
[Go](https://golang.org/dl/) version 1.13 or higher is required. 

To download the project, run the following command 

```
go get -v github.com/mikanikos/Fork-Accountability
```

### Building project

After this, go to the project root directory.

To build all the project files, run the following command:

```
go build -v ./...
```

### Running tests

To run all the tests in the project directory and sub-directories and generate a report on the test coverage, run the following command:

```
go test -v ./... -covermode=count -coverprofile=coverage.out
```

You can then inspect the coverage in your browser by running the following command:

```
go tool cover -html=coverage.out
``` 

It's also possible to run the tests in a specific package with the following command:

```
go test -v ./[package_path] -covermode=count -coverprofile=coverage.out
```

Note that CI/CD is enabled for this project and it's possible to inspect the build status and detailed information about the test coverage directly on Github and Codecov.

### Running the monitor

Go to the [monitor](cmd/monitor) directory inside the [cmd](cmd) directory, compile with the following command:

```
go build
```  

Then, simply run the generated binary. The monitor accepts the following command-line parameters:

- **-config**: path (relative to the project root directory) of the configuration file for the monitor (default "cmd/monitor/_config/config.yaml")

- **-delay**: time to wait (in seconds) before start running, use for testing (default 0)

- **-report**: path (relative to the project root directory) of the report to generate at the end of the execution instead of printing logs to standard output (default "")

The yaml configuration file must have the following parameters in order to provide the monitor with the required information to run the algorithm:

- `height`: it represents the consensus instance where the fork has been detected or the height where the fork accountability algorithm will be run. This parameter will be used to request messages from the validators.

- `firstDecisionRound`: the round where the first decision was made in the consensus algorithm.

- `secondDecisionRound`: the round where the second decision was made in the consensus algorithm.
 
- `timeout`: timer (in seconds) used to exit the execution after the time is expired. It's not needed by the algorithm but it's just a safety measure to prevent a blocking state in case something goes wrong. It can be set at a very high value and will not affect the execution of the algorithm.

- `validators`: is the list of addresses where the validators are listening for incoming monitor requests 

The [_config](cmd/monitor/_config) folder contains some sample config files for the monitor.

### Running the validator

Go to the [validator](cmd/validator) directory inside the [cmd](cmd) package, compile with the following command:

```
go build
```  

Then, simply run the generated binary. The validator accepts the following command-line parameters:

- **-config**: path (relative to the project root directory) of the configuration file for the validator (default "/cmd/validator/_config/config_1.yaml")

- **-delay**: time to wait (in seconds) before replying back to the monitor, use for testing (default 0)

The yaml configuration file must have the following parameters in order to provide the validator with the required information to run correctly:

- `id`: unique id of the validator 
- `address`: address used to listen for incoming requests from the monitor
- `messages`: list of messages organized with the following structure
  
      [height]:
        heightvoteset:
          [round]:
            received_prevote:
              - type: [PREVOTE | PRECOMMIT]
                sender: [sender_id]
                round: [round]
                value:
                  data: [value]
            
            sent_prevote:
              - type: [PREVOTE | PRECOMMIT]
                sender: [sender_id]
                round: [round]
                value:
                  data: [value]
                
            
            received_precommit:
              - type: [PREVOTE | PRECOMMIT]
                sender: [sender_id]
                round: [round]
                value:
                  data: [value]

            
            sent_precommit:
              - type: [PREVOTE | PRECOMMIT]
                sender: [sender_id]
                round: [round]
                value:
                  data: [value]

The value in square brackets are values and they are positive integers except for `type` (PREVOTE or PRECOMMIT) and the `data` fields (it can be any integer value, the type can be changed).

The [_config](cmd/validator/_config) folder contains some sample config files for the validator.


### Running test scripts

It's possible to run bash scripts (in a Unix environment) in order to run more validator instances and the monitor at the same time and easily test different scenarios.
A sample bash script is present in the scripts folder and gives a very minimal example of a simple experiment. 

## Acknowledgments

The project is developed as a semester project in collaboration with the [Distributed Computing Lab](https://dcl.epfl.ch/site/) at [EPFL](https://www.epfl.ch/en/) and [Informal Systems](https://informal.systems/).

Thank you [Jovan](https://github.com/jovankomatovic5) and [Adi](https://github.com/adizere) for providing all the theoretical background and advice during the development of the project.
