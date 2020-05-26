# Fork accountability related problems 

Fork accountability is not the only problem that must be addressed when a fork occurs. Indeed, a reliable consensus algorithm should also provide guarantees on the detection and recovery from a fork, which are complementary problems to fork accountability. In this section, we analyze these two aspects and provide intuitive solutions to solve both problems in the Tendermint Consensus algorithm.

## Fork detection

A fork detection solution aims to detect a fork as soon as possible in a reliable and efficient way.

We define the following interface for a fork detection algorithm:

- **detect(h)**: fork has occurred at height *h*

The fork detection algorithm guarantees the following properties:

- **accuracy**: if the accountability algorithm detects a fork at height *h*, then a fork occurred at height *h*

- **completeness**: if a fork occurred at height *h*, the algorithm eventually detects the fork

As already said regarding the fork accountability algorithm, we can assume that *n = 3f + 1* and that at most *2f* faulty processes are present in the system.

A simple solution for the fork detection problem would require having a trusted entity (we can keep assuming the monitor itself) to be connected to all validators and periodically check for a fork by asking validators for their current state.
Correct validators will send their state to the monitor as soon as they receive a check request. The monitor would be able to verify the validity of a validator state and will detect a fork as soon as it detects two different commits in the same height.
Once it has detected the fork, the monitor will start the fork accountability algorithm in order to detect the faulty processes, as described before.

A simple illustrative pseudo-code version of this idea is shown here below:

```
// monitor initializes a map string -> State, where State is a protocol-specific type, that allows to keep track of the latest states among all validators
init():
	statesPerValidatorID = init(map: string -> bool)

// monitor periodically checks for forks 
onTimeoutCheck():
    requestStatesFromValidators()   

// monitor delivers a new state from a validator (with unique ID) as soon as a new one arrives
onStateReceived(state, ID):
    
	if isValid(state):
        statesPerValidatorID[ID] = state

		// check if there are two decisions for a single height
		if checkForFork(states):
            startForkAccountabilityAlgorithm()

```

This is not meant to be a complete solution to the problem because it just provides a basic intuitive idea to handle the fork detection. It also lacks a formal theoretical support which is required for the full refinement of the algorithm.
The validation and a possible proof of concept of this algorithm is left to future work.

## Fork recovery

A fork recovery mechanism allows validators to agree on the state of the system after a fork occurred. 
This step should occurr after validators have run the fork accountability algorithm and it includes the ban of faulty processes (detected by a fork accountability algorithm) from the validator set and the possible restart of the blockchain with the new validator set.

It is clear that the Tendermint consensus algorithm cannot proceed further without having enough validators, namely 4 validators, in order to make the algorithm work correctly. In fact, there is no quorum of at least *2f + 1* processes that allow to satisfy the Tendermint consensus rules when there are only 4 validators in the system. 
Therefore, if the number of left processes after the banning phase is lower than 4, the consensus algorithm should block and should be able to restart only when another valid process will join the validator set based on the protocol rules.
Unless this case occurs, processes need to find a way to restart the system from a stable, common point such that invalid decisions and modifications due to the faulty configuration of the system are discarded.

The simplest idea to achieve this is giving validators the ability to save and restore the state of the system before a fork occurred.
In other words, processes should be able to save the current state of the blockchain and their internal state every time before starting a new consensus instance (new height). In this way, if a fork occurs at height h, processes would be able to restore their previous state before starting height *h* and re-execute the consensus instance of that height with the new correct validator set.

A simple illustrative pseudo-code version of this idea is shown here below:

```
// each validator initializes a map int -> State, where State is a protocol-specific type, that allows to keep track of the state before starting a new consensus instance
init():
    validatorSet = getCurrentValidatorSet()
    heightStateMap = init(map: int -> State)

// before starting a consensus instance, each validator saves the current system state
startHeight(h):
    heightStateMap[h] = getCurrentSystemState()	

// after a validator received the fork accountability output (faultyProcesses) about a fork on a certain height
onForkAccountabilityCompletion(faultyProcesses, height):

    validatorSet = validatorSet - faultyProcesses

    if validatorSet >= MIN_CONSENSUS_PARTICIPANTS:
        restoreState(heightStatesMap[height])
        startHeight(height)
    else:
        throwError("The new set of validators doesn't have enough participants to execute the consensus algorithm")

```

Although this solution could be considered valid at first, a future study should be dedicated to the complete validation and study of this algorithm.

## Decentralized fork accountability
The fork accountability solution outlined in the previous sections relies on the monitor which is a centralized entity. It would be desirable executing the whole accountability algorithm among validators in a distributed way, without having any trusted third-party components.

The problem is not trivial and there is not even certainty that is solvable. 
Intuitively, processes need to run a distributed algorithm that allows reaching consensus on the faulty processes in the system.
However, the system itself is compromised because there are more than *f* faulty processes (i.e., there is not a quorum of correct processes that allows running the consensus correctly).

The study and the design of a possible solution to such a problem could be further analyzed in a future work.



