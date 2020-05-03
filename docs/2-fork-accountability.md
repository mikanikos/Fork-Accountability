# Fork accountability

As we know, the agreement property of the Tendermint consensus algorithm cannot be violated if there are at most *f* faulty processes in the system.
What happen if more than *f* processes are faulty?

A well-designed consensus protocol should provide some guarantees when this limit is exceeded. The most important guarantee is **fork accountability**, where the processes that caused the consensus to fail can be identified and punished according to the protocol specifications. 

## The problem

We can separate two cases:

- if the number of faulty processes is more than *2f*, we cannot do anything because faulty processes basically control the system and they can do whatever they prefer.

- if the number of faulty processes is less than *2f*, i.e.

      f < num. faulty processes <= 2f
      
  Agreement can still be violated but we're in condition to detect the processes that were faulty and punish them, for example by removing them from the validators set. 
  
A **fork** happens when two correct validators decide on different blocks in the same height, i.e. there are two commits for different blocks at the same height. 
In the context of the Tendermint protocol, a fork happens in height h when a quorum of messages (2f + 1 Precommit messages) are sent for different values, i.e. we have two sets of 2f + 1 Precommit messages, A and B, such that A has value v and B has value v' != v in height h.  

We aim to give incentives to validators to behave correctly according to the protocol specifications and detect faulty validators, without mistakenly detecting correct validators. 
 
## Fork reasons
 
There are multiple reasons that can lead to a fork, all coming from processes that deviate from the protocol specification and not following the behavior of a correct process.

We can summarize the most important ones:

- A process sends multiple proposals in a round *r* for different values

- A process sends a Prevote/Precommit message for an invalid value *v* in a round *r*

- A process sends multiple Prevote/Precommit messages in a round *r* for different values

- A process sends a Prevote/Precommit message for a value *v* despite having a lock on a different value *v'*

- A process sends a Precommit for a value *v* without having received at least 2/3 Prevote messages for *v* 

Given this problem, we want to develop an accountability algorithm that, by analyzing the messages logs from each validator, is able to determine which processes respected the protocol and which ones didn't. 
In order to prove the misbehavior of the processes that didn't respect the protocol, we want to show the proof of their faultiness respect to the Tendermint consensus protocol.

## Accountability algorithm

### Interface and properties

The accountability algorithm analyzes the received message logs and detects faulty processes from the message logs received. 
The algorithm exposes the following interface:

- **detect(P)**: process P is faulty according the accountability algorithm

The algorithm guarantees the following properties:

- **accuracy**: if the accountability algorithm detects a process *p*, then process *p* is faulty

- **f+1-completeness**: the accountability algorithm eventually detects *f + 1* different processes

All the processes that are not detected are considered correct.

## Design of an accountability algorithm

We introduce a trusted third-party verification entity which is responsible to run the accountability algorithm - we call it **monitor**.

The monitor is responsible to coordinate and run the accountability algorithm as soon as a fork is detected. The input of the algorithm are the message logs of the validators, that's why the monitor will first request the message logs from each validator and then will run the accountability algorithm.

We assume that all validators keep track of all the messages they sent and received during the execution of the Tendermint consensus algorithm. 
The idea is that a correct process would be able to prove that it was not faulty by showing its activity (message logs) to the monitor.
Indeed, the monitor is able to detect all the faulty processes that caused a fork (and not only) by analyzing the message logs from all the correct processes.

## Fork scenarios

In order to better understand the problem, let's see when the fork can happen and how the monitor would be able to handle the situation when analyzing the logs for each validator.

### Fork in same round

If a fork happens in a certain round r (two decisions are made in the same round r for value v and v'), we know that at least f+1 processes have sent a Prevote/Precommit message for both v and v'.
Since at least one correct process will receive both messages, the monitor is able to infer all the "hidden" messages and eventually identify the faulty processes who caused the fork.  

### Fork in different rounds

If a fork happens in different rounds (value v is decided in round r1 and value v' is decided in round r2), we know that at least f+1 processes have sent a Precommit message for value v in round r and also sent a Prevote message for value v' in round r2 despite having locked value v when they sent a Precommit.

Therefore, unless these processes have a justification (other 2f+1 Prevote messages for value v' in a round r'' s.t. r1 < r'' < r2), they are faulty. If they do have such a justification, then we now have to look for f+1 processes that sent Prevote messages for value v' despite having locked value v.

As we can see, this is a recursive problem and we need to look back until we find f+1 processes that can't justify the sending of a Prevote/Precommit message for another value despite having locked another one. 
The method is simple: either a process is faulty because can't justify the sending of a message or we go back recursively until we find processes that don't have a justification. It is clear that we'll eventually find an iteration of the recursive problem where at least f+1 processes sent a Prevote message without a valid justification. 

## Challenges

It's possible to receive "altered" message logs and the monitor might not have enough information to detect all the faulty processes that led to a fork. 
However, we rely on the fact that Tendermint messages are signed and it's not possible to forge messages, i.e. some process *p* can't state that it has received a message from process *p'* if *p'* did not send that message. 

So, excluding the previous problem, the following situations are possible:
 
- A process *p* denies having received a message *m* (*p*'s message logs does not contain *m* as a message received): this doesn't really affect the monitor's job because it doesn't contribute to create a fork and we can ignore this scenario  
- A process *p* denies having sent a message *m* (*p*'s message logs does not contain *m* as a message sent): this can lead to a fork so we need to consider the fact that processes can do this to hide their faultiness. However, if *m* led to a fork, at least one correct process has received that message and we can identify the faulty process. 

Given this overview, it's now clear why the monitor would just need the message logs from the correct processes: if a message *m* from a process *p'* has been received by a process *p* and *p'* denies sending m, we can add this missing message in *p'* 's message logs without needing to take any other action. 

Therefore, we know that the monitor can "infer" the valid message logs for each process after receiving the message logs from all the correct processes. The next step would be to simply analyze the logs for each process and determine whether it's faulty or not.
However, we still need to ensure that we're going to receive the message logs from the correct processes and this depends on how we model the communication between the monitor and validators.

### Communication between the monitor and validators

The communication between the monitor and the validators for receiving message logs of the Tendermint consensus algorithm is a crucial aspect of the monitor execution.

For example, what happen when the monitor doesn't receive message logs from one of the validators? Should this validator be considered faulty?

#### Synchronous model

If we model the communication with a synchronous model, the only option would be waiting some time to receive all the message logs before running the accountability algorithm.
After that time, if some process p did not send its message logs, the only thing monitor can do is considering p faulty, even though it doesn't have all the information to determine its faultiness with certainty. 
On the other hand, if it doesn't consider p faulty, the monitor might not be able to find all the faulty processes that led to a fork.

#### Asynchronous model

This method would work but we want to make no assumptions regarding the communication between the monitor and the validators. 
If we don't make any assumption regarding the communication, we would realize that it's impossible to design an accountability algorithm where the communication between monitor and validators is completely asynchronous.

The reason is simple: in an asynchronous setting, the monitor cannot be sure a process P is indeed faulty if it didn't receive P's message logs. 
In this scenario, the accountability algorithm can only determine the faultiness of P by analyzing the messages logs sent by other processes. At this point, the unique faultiness that can be detected from the other processes' message logs is equivocation, i.e. a process sent more than one message with the same type in the same round but with different values.
This is easily visible from the fact that, by receiving messages logs from correct processes, we'll be able to determine if some correct process received more than one message with the same type from another process p (i.e. if some process equivocated).  
However, we've also shown that a process can be faulty in other ways and that would not be enough for designing an asynchronous accountability algorithm that catches all the faulty processes that led to a fork. 
In fact, it can happen that, by not receiving message logs from one message, we might end up mistakenly detecting correct processes because we don't have enough information to determine if a process has enough justifications to have sent a Prevote/Precommit message (accuracy property violated). 

If we want to design a correct accountability algorithm, we need to slightly modify the Tendermint consensus algorithm in order to correctly respect the accuracy property.
As we said above, we need to catch the case where a fork happens without having any process equivocating.

By looking at the other faultiness reasons, we can see that the critical case is when a process sends a Precommit message for a value v and then, in a later round, it sends a Prevote message for a different value than v.
The asynchronous accountability algorithm would work correctly if processes can justify the sending of a Prevote message: a simple solution is to attach the justifications to the Prevote message itself so that the accountability algorithm can directly check if the message has been correctly sent.

This would be enough because, intuitively, there will be at least one correct process that has received both messages and it would be possible to track down the origin and the justifications and find out if a process is really faulty.

Therefore, the monitor will be able to complete the accountability algorithm as soon as at least f+1 processes will be detected.
We also know that at least f+1 messages logs will be received by the monitor (we assume correct processes will send their logs). The algorithm will be able to complete the algorithm and respect the completeness and accuracy properties with the message logs of the correct processes.

## Algorithm design
 
To summarize the above discussion, these are the high-level steps carried out by the monitor to run the asynchronous accountability algorithm:

1. Upon a fork, monitor requests message logs from all the validators

2. Monitor runs the accountability algorithm to detect faulty processes upon receiving at least *f + 1* message logs and at every new message logs received when there are at least *f + 1* message logs 

3. While running the accountability algorithm, the monitor analyzes the received logs and build the correct logs for each process: for every message *m* that has been received by some other process but is not present in the sender's message logs, the monitor attaches *m* to the message logs of the sender of *m*. 

4. For each validator, the monitor scans the received message logs and the "inferred" message logs and determines for each process whether it's faulty or not by analyzing the history of sent messages and transitions

5. When the monitor detects at least *f + 1* faulty processes, the algorithm completes

Termination is guaranteed by the fact that correct processes will send their messages logs (otherwise, why they shouldn't if they have nothing to hide?), and their message logs will be correct (no sent message will be missing). 
If this condition doesn't hold, the algorithm would not be able to identify correctly all the faulty processes that led to a fork. However, correct processes have no reason to misbehave. 

## Notes

- The algorithm can also catch other faulty behaviours that didn't necessarily lead to a fork - that's a positive side effect of the accountability algorithm described above.

- The monitor needs to analyze the message logs from the round of the first decision to the round of the second decision. There's no need to send to the monitor or analyze previous message logs. 

- The accountability algorithm runs for a single instance of consensus, i.e. for a single height. There's no need to run the algorithm on multiple heights.  
