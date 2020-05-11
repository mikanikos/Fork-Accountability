# Tendermint protocol

We are going to present a short and simple overview of the [Tendermint](https://github.com/tendermint) consensus algorithm and the main aspects that are required to understand the fork accountability problem. 
Please note that this is not an exhaustive and formal discussion. We invite the reader to check out the [Tendermint documentation](https://docs.tendermint.com/master/) and the [official paper](https://arxiv.org/abs/1807.04938) that describe the theoretical aspects in more detail.   

## Consensus algorithm

Tendermint is a partially synchronous Byzantine fault tolerant (BFT) consensus protocol. The protocol requires a fixed set of *n* validators that attempt to come to consensus on one value or, more precisely, on a block (a block is a list of transactions) using gossip-based communication. 

Tendermint is notable for its simplicity and performance and is currently one of the best consensus algorithm used in production.

### Properties 

The Tendermint consensus algorithm guarantees the following properties:

- **Agreement**: No two correct processes decide on different values

- **Validity**: A value decided is a valid value, with respect to the specifications of validity of a message

- **Termination**: All correct processes eventually decide

- **Integrity**: No correct process decides more than once in a single height 

In order to ensure all these properties, the Tendermint consensus algorithm assumes a correct majority of processes such that less than one third of processes are faulty (*f*):
 
    n > 3f
 
When a process *p* is faulty, *p* can behave in an arbitrary way and we have no assumptions or guarantees regarding its actions.
 
In order to simplify the whole discussion, we assume the following from now on:
 
    n = 3f + 1

### Overview

The consensus algorithm proceeds in rounds: in each round there is a validator that proposes a value (proposer) and all the validators then vote during the round on whether to accept the proposed value or move on to the next round (proposers are chosen according to their voting power).

During a round, processes can exchange different types of messages:
 
- **PROPOSAL message**: sent by the proposer of the round to propose a value in a new round

- **PREVOTE message**: sent by validators to vote for a value during the first voting step

- **PRECOMMIT message**: sent by validators to vote for a value during the second voting step
 
We can summarize the execution in a round in **two voting steps**: PREVOTE and PRECOMMIT steps. A vote can be for a particular value or for *nil* (null value).

A correct process decides on a value *v* in a round *r* upon receiving a proposal and *2f + 1* quorum of PRECOMMIT messages for *v* in a round *r*. 
A correct process sends a PRECOMMIT message for a value *v* in a round *r* upon receiving a proposal and *2f + 1* quorum of PREVOTE messages for *v* in a round *r*.

Validators wait some time before sending a PREVOTE for *nil* if they do not receive a valid proposal after a certain time and they send a PRECOMMIT message for *nil* if they do not receive *2f + 1* PREVOTE messages for a value.
If a correct process receives at least *2f + 1* PRECOMMIT messages for *nil* in a round, it moves to the next round.

After a decision has been made, processes continue to agree on other values on another consensus instance (*height*) and they repeat the process described above in order to agree on different transactions.

### Rules

In order to ensure that processes will eventually come to a decision, there are some constraints that are applied. These rules aim to prevent any malicious attempt to cause more than one block to be committed at a given height. 

When a validator *p* sends a PRECOMMIT for a value *v* at round *r*, we say the *p* is *locked* on *v*. Validators can propose and send PREVOTE messages for a value they have locked and they can only change the locked value if they receive a more recent proposal with a quorum of PREVOTE messages. 

These conditions ensure both the safety and liveness properties of the consensus algorithm since a validator cannot send a PRECOMMIT message without sufficient evidence.

## Agreement property

We limit the discussion to the agreement property only because its proof can be useful to better understand the fork accountability problem. 
We invite interested readers to read the official paper for a more formal, complete and rigorous proof of all the properties listed above.

In this section we will give a simple proof of the agreement property for the Tendermint algorithm and show why the correct majority assumption given before (*n > 3f*, or simply *n = 3f + 1*) is important for ensuring this property.

The key idea behind the proof of the agreement property is that any two sets of *2f + 1* processes have at least one correct process in common with respect to the Tendermint consensus algorithm.

For simplicity, we assume that *n = 3f + 1*.
Since there are two sets of *2f + 1* processes, their sum can be written as:
 
    2(2f + 1) = 4f + 2 = 3f + 1 + f + 1 = n + f + 1
     
That means that the intersection of these two sets contains at least *f + 1* processes or, in other words, at least one correct process.  

From this result, we will show in the next sections that it is not possible to violate the agreement property in the Tendermint consensus algorithm if there are at most *f* faulty processes.

### Idea behind the proof

The idea behind the proof is the impossibility of having two sets of *2f + 1* PREVOTE or PRECOMMIT messages in the same round *r* for different values. If we assume the opposite by contradiction, then it would be possible to have at least *f + 1* processes that sent both the messages, which would mean *f + 1* are faulty. But at most *f* processes can be faulty in Tendermint, so the contradiction.

Moreover, at most one value can be locked in a round and, if a correct process decided value *v* in round *r*, *v* will be locked for all the next rounds after *r* on that height. This guarantees it is not possible to have another quorum for committing another value different than *v*.
### Proof

Assume that a correct process *p* decides value *v* in round *r* and height *h*. We want to prove that any other correct process *p'* in some round *r' >= r* of height *h* decides *v'* such that *v'* = *v*.

We have two cases:

- if *r' = r*: 

    *p'* has decided value *v'* in round *r*, therefore it has received at least *2f + 1* PRECOMMIT messages for value *v'* in round *r*. 
    Similarly, *p* has decided value *v* in the same round *r*, therefore it has received at least *2f + 1* PRECOMMIT messages for value *v* in round *r*. 
    As it has been shown previously, any two sets of *2f + 1* messages intersect in at least one correct process. Since a correct process only sends a single PRECOMMIT message, it must be that *v = v'*.  

- if *r' > r*:
        
    *p* has decided value *v* in the same round *r*, therefore it has received at least *2f + 1* PRECOMMIT messages for value *v* in round *r*.
    Since the number of faulty processes is at most *f*, at least *f + 1* correct processes have locked value *v* by round *r* and, by algorithm locking rules, they will send PREVOTE messages only for value *v* or *nil* in subsequent rounds of height h.
    
    In a similar way, *p'* has decided value *v'* in round *r' > r*, so it has received at least *2f + 1* PRECOMMIT messages for value *v'* in round *r' > r*. 
    Since the number of faulty processes is at most *f*, at least *f + 1* correct processes have locked value *v'* by round *r' > r* and these processes must have received at least 2f + 1 PREVOTE messages for value v' by round r'.
    
    From the fact that the intersection of *2f + 1* processes has at least one correct process, at least one correct process that locked value *v* in round *r* also sent a PREVOTE message for value *v'* in a later round *r' > r*. 
    Since this is impossible, it can only be that *v = v'*.

Since these are the only two possible cases, the above reasoning proves that the agreement property is satisfied.

## Conclusion

The Tendermint consensus protocol is a simple and efficient algorithm that guarantees the four properties mentioned above. These guarantees are valid as long as faulty processes are at most one third of the total number of validators in the system.

However, when this assumption is not valid anymore, processes can decide on different values in the same consensus instance and the blockchain can potentially diverge into two potential decision paths, making the entire protocol break.

Although it is not possible to provide all the four properties in this scenario, it can be desirable having some other "safety" guarantees that allow the system to self-recover from the possible fault. 

The problem of fork accountability deals with this situation and will be addressed in the next section.