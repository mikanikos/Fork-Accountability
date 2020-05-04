# Tendermint protocol overview

We are going to present a very short overview of the [Tendermint](https://github.com/tendermint) consensus algorithm and the main aspects that are necessary in order to understand the fork accountability algorithm. 
Please note that this is not an exhaustive and formal description, we invite the reader to check out the [Tendermint documentation](https://docs.tendermint.com/master/) and the [official paper](https://arxiv.org/abs/1807.04938) about the consensus algorithm.   

## Consensus algorithm

Tendermint is a partially synchronous BFT consensus protocol. The protocol requires a fixed set of validators that attempt to come to consensus on one value or, more precisely, on a block (a block is a list of transactions) using gossip-based communication. The consensus algorithm proceeds in rounds: in each round there is a validator that proposes a value (proposer) and all the validators then vote during the round on whether to accept the proposed value or move on to the next round (proposers are chosen according to their voting power).

As most of the consensus algorithms, the Tendermint consensus algorithm assumes a correct majority of processes: less than one third of processes are faulty (*f*), therefore the total number of processes *n* is strictly greater than *3f* (*n > 3f*). For simplicity, we can asssume that:
 
    n = 3f + 1

During a round, processes can exchange different messages and we can summarize the execution in two voting steps: Prevote and Precommit steps. A vote can be for a particular value or for a null value.

A correct process decides on a value *v* in a round *r* upon receiving a proposal and *2f + 1* quorum of Precommit messages for *v* in a round *r*. 
A correct process sends a Precommit message for a value *v* in a round *r* upon receiving a proposal and *2f + 1* quorum of Prevote messages for *v* in a round *r*.

Validators wait some time before sending a Prevote with a null value if they do not receive a valid proposal after a certain time and they send a Precommit message with a null value if they do not receive *2f + 1* Prevote messages for a value.
If a correct process receives more than 2/3 Precommit messages for a null value in the same round, it moves to the next round.

After a decision has been made, processes continue to agree on other values on another height (consensus instance) and they repeat the process described above in order to agree on different transactions.

In order to ensure that processes will eventually come to a decision, there are some constraints that are applied. These rules are used to prevent any malicious attempt to cause more than one block to be committed at a given height. 
When a validator *p* sends a Precommit for a value *v* at round *r*, we say the *p* is *locked* on *v*. Validators can propose and send Prevote messages for a value they have locked and they can only change the locked value if they receive a more recent proposal with a quorum of Prevote messages. 

These conditions ensure both the safety and liveness properties of the consensus algorithm because a validator cannot send a Precommit message without sufficient evidence and cannot send a Precommit message for a different value at the same time.

Tendermint is notable for its simplicity and performance and is currently one of the best consensus algorithm used in production.

## Properties of consensus 

The Tendermint consensus algorithm guarantees the following properties:

- **Agreement**: No two correct processes decide on different values

- **Validity**: A value decided is a valid value, with respect to the specifications of validity of a message

- **Termination**: All correct processes eventually decide

- **Integrity**: No correct process decides more than once in a single height 

We limit the discussion to the agreement property only because its proof can be useful to better understand the fork accountability problem. 
We invite interested readers to read the official paper for a more formal, complete and rigorous proof of all the properties listed above.

## Agreement property

In this section we will give a simple proof of the agreement property for the Tendermint algorithm and show why the correct majority assumption we gave before (*n > 3f*, or simply *n = 3f + 1*) is important for ensuring this property.

The key idea behind the proof of the agreement property is that any two sets of *2f + 1* processes have at least one correct process in common with respect to the Tendermint context.

We know that Tendermint assumes that *n > 3f*, for simplicity we assumed that *n = 3f + 1*.
We have two sets of *2f + 1* processes, their sum can be written as:
 
    2(2f + 1) = 4f + 2 = 3f + 1 + f + 1 = n + f + 1
     
That means that the intersection of these two sets contains at least *f + 1* processes, therefore at least one correct process.  

Given this result, it is quite easy to understand that, if there are at most *f* faulty processes, it is not possible to violate the agreement property in the Tendermint consensus algorithm, as we show in the next section.

### Idea behind the proof

The idea behind the proof is that it can never happen that there are two sets of *2f + 1* Prevote or Precommit messages in the same round *r* for different values. If that would be the case, we would have at least *f + 1* processes that sent both the messages, which would mean *f + 1* are faulty. But at most *f* processes are faulty, so this is impossible.
Moreover, at most one value can be locked in a round and if a correct process decided value *v* in round *r*, *v* will be locked for all the next rounds after *r* on that height. This guarantees it is not possible to have another quorum for committing another value different than *v*.
Therefore, it is clear that *v* will be the only value that could be decided from that moment onwards.

### Proof

Let us assume that a correct process p decides value v in round r of height h. We want to prove that any other correct process p' in some round r' >= r of height h decides v' such that v' = v.

We have two cases:

- if r' = r: 

    p' has decided value v' in round r, so it has received at least 2f + 1 Precommit messages for value v' in round r. 
    p has decided value v in the same round r, so it has received at least 2f + 1 Precommit messages for value v in round r. 
    As it has been shown previously, any two sets of 2f + 1 messages intersect in at least one correct process. Since a correct process only sends a single Precommit message, it must be that v = v'.  

- if r ' > r:
        
    p has decided value v in the same round r, so it has received at least 2f + 1 Precommit messages for value v in round r.
    Since the number of faulty processes is at most f, at least f+1 correct processes have locked value v by round r and, by algorithm rules, they will send Prevote messages only for value v or nil in subsequent rounds of height h.
    
    p' has decided value v' in round r' > r, so it has received at least 2f + 1 Precommit messages for value v' in round r' > r. 
    Since the number of faulty processes is at most f, at least f+1 correct processes have locked value v' by round r' > r and these processes must have received at least 2f + 1 Prevote messages for value v' by round r'.
    
    From the fact that the intersection of 2f + 1 processes has at least one correct process, at least one correct process that locked value v in round r also sent a Prevote message for value v' in a later round r' > r. 
    Since this is impossible, it can only be that v = v'.

Since these are the two only possible cases, the above reasoning proves that the agreement property is satisfied.


In conclusion, when more than one third of processes are faulty, the agreement property of the consensus algorithm could be violated. Otherwise, the agreement property is always satisfied.