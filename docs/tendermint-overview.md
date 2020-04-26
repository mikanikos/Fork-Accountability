# Tendermint protocol overview

We are going to present a very short overview of the Tendermint consensus algorithm and the main aspects that are necessary in order to understand the fork accountability algorithm. Please note that this is note an exhaustive and formal description, we invite the reader to check out the Tendermint documentation and the official paper.   

## Consensus algorithm

Tendermint is a partially synchronous BFT consensus protocol. The protocol requires a fixed set of validators that attempt to come to consensus on one value or, more precisely, on a block (a block is a list of transactions) using gossip-based communication. The consensus algorithm proceeds in rounds: in each round there's a validator that proposes a value (proposer) and all the validators then vote during the round on whether to accept the proposed value or move on to the next round (proposers are chosen according to their voting power).

As most of the consensus algorithms, the Tendermint consensus algorithm assumes a correct majority of processes: less than 1/3 of processes are faulty (f), therefore the total number of processes n > 3f. 

During a round processes can exchange different messages and we can summarize the execution in two voting steps: Prevote and Precommit steps. A vote can be for a particular value or for a null value.

A correct process decides on a value v upon receiving a proposal and 2f+1 Precommit messages for v in a round r. A correct process sends a Precommit message for a value v upon receiving a proposal and 2f+1 Prevote messages for v in a round r.

Validators wait some time before sending a Prevote with a null value if they don't receive a valid proposal after a certain time and they send a Precommit message with a null value if they don't receive 2f+1 Prevote messages for v in a round r.

If a correct process receives more than 2/3 Precommit messages for a null value in the same round, it moves to the next round.

After a decision has been made, processes continue to agree on other values on another height (consensus instance) and they repeat the process described above.

In order to ensure that processes will eventually come to a decision, there are some constraints (or rules) that are applied. This is also used to prevent any malicious attempt to cause more than one block to be committed at a given height if the majority assumption is respected. 
When a validator sends a Precommit for a value at round r, we say it is locked on that value. Validators can propose and send Prevote messages for a value they have locked and they can only change the locked value if they receive a more recent proposal with a quorum of Prevotes messages. 

These conditions ensure both safety and liveness of the consensus algorithm because a validator cannot send a Precommit message without sufficient evidence and cannot send a Precommit message for a different value at the same time.

Tendermint is notable for its simplicity and performance and is currently one of the best consensus algorithm used in production.


## Idea behind the proof of the Agreement property

The Tendermint algorithm guarantees the following properties:

- Agreement: No two correct processes decide on different values
- Validity: A value decided is a valid value, respect to the specifications of validity of a message
- Termination: All correct processes eventually decide
- Integrity No correct process decides more than once in a single height 

In the next section we'll give a very intuitive proof of the agreement property for the Tendermint algorithm and show why the majority assumption is important for ensuring this property: therefore, when more than 1/3 of processes are faulty, a fork can happen.

We limit the discussion to the agreement property only its proof can be useful to better understand the fork accountability problem. We invite interested readers to read the official paper for a more formal, complete and rigorous proof.

### Idea

The key idea behind the proof of the agreement property is that any two sets of 2f+1 processes have at least one correct process in common in the Tendermint context.

We know that we have n > 3f, for simplicity we can assume:
 
    n = 3f + 1

We have two sets of 2f + 1 processes, their sum can be written as:
 
    2(2f + 1) = 4f + 2 = 3f + 1 + f + 1 = n + f + 1.
     
That means that the intersection of these two sets contain at least f + 1 processes, therefore at least one correct process.  

Given this result, it's quite easy to understand that, if there are at most f faulty processes, it's not possible to violate the agreement property in the Tendermint consensus algorithm.

In fact, it can never happen that there are two sets of 2f+1 Prevote or Precommit messages in the same round r for different values. If that would be the case, we'll have at least f+1 processes that sent both the messages, which would mean f+1 are faulty. But at most f processes are faulty, so this is impossible.

This comes from the fact that at most one value can be locked in a round and if a correct process decided value v in round r, only v can be locked in all the next rounds on that height and v is the only value that could be decided from now on.
