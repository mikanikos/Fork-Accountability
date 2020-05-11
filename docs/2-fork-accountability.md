# Fork accountability

The agreement property of the Tendermint consensus algorithm can be violated if there are more than *f* faulty processes in the system, as we discussed in [1-tendermint-overview.md](./1-tendermint-overview.md).

A well-designed consensus protocol should provide some guarantees when this limit is exceeded in order to allow the system to recover from a possible fault. An important guarantee is **fork accountability**, which is the topic of this document. 

## The problem

A fork is a split in the blockchain network or, with other words, a divergence into different block paths. 
A **fork** happens when two correct validators decide on different blocks in the same height (i.e., there are two commits for different blocks in the same consensus instance). 

In the context of the Tendermint protocol, a fork happens in height *h* when a quorum of messages (*2f + 1* PRECOMMIT messages) are sent for different values: more formally, there are two sets of *2f + 1* PRECOMMIT messages, A and B, such that A's messages are PRECOMMIT votes for value *v* and B's messages are PRECOMMIT votes for value *v' != v* in height *h*.  

The Tendermint conensus algorithm guarantees that a fork (i.e., a violation of the agreement property) cannot happen when the number of faulty processes is less than one third.
We want to address and analyze the problem of defining a Tendermint consensus algorithm that provides the same four properties (validity, agreement, termination, integrity) when the number of faulty processes is less than one third and guarantees, in addition, an efficient detection of faulty processes that can potentially make the consensus fail with a fork when there are more than one third of faulty validators.

In particular, we assume there are more than one third of faulty validators such that the following holds:

    f < num. faulty processes <= 2f

In the Tendermint protocol, all processes should be accountable for their actions and faulty validators should be identified and punished according to the protocol specifications, for example by removing them from the validators set.
We aim to give incentives to validators to behave correctly according to the protocol specifications and detect faulty validators, without mistakenly detecting correct validators. 
 
## Misbehaviors
 
In the previous section, we mentioned that a process can make the consensus fail with a fork. There are, indeed, multiple reasons that can generate a fork in the blockchain, all coming from processes that deviate from the protocol specification in various ways.

We can summarize some of the most important misbehaviors that can possibly cause a fork in a consensus instance according to the Tendermint consensus algorithm rules: 

- A process sends multiple proposals in a round *r* for different values 

- A process sends a PREVOTE/PRECOMMIT message for an invalid value *v* in a round *r*

- A process sends multiple PREVOTE/PRECOMMIT messages in a round *r* for different values (**equivocation**)

- A process sends a PREVOTE/PRECOMMIT message for a value *v* despite having a lock on a different value *v'*

- A process sends a PRECOMMIT for a value *v* without having received at least two thirds PREVOTE messages for *v* 

Processes are considered correct if they do not make any of the above-mentioned misbehaviors.

## Design of an accountability algorithm

We aim to develop an accountability algorithm that, by analyzing the messages received by validators during an execution of a consensus instance, is able to determine which processes respected the protocol and which ones did not and led to a fork. 
In order to prove the misbehavior of the processes that did not respect the protocol, we also want to show the proof of their misbehavior with respect to the Tendermint consensus protocol.

### Interface

The algorithm exposes the following interface (in the event of a fork):

- **detect(P)**: process P is faulty according to the accountability algorithm

### Properties

The algorithm guarantees the following properties:

- **accuracy**: if the accountability algorithm detects a process *p*, then process *p* is faulty

- **f+1-completeness**: the accountability algorithm eventually detects *f + 1* different processes

All the processes that are not detected are considered correct.

## Monitor

We introduce a trusted third-party verification entity which is responsible to run the accountability algorithm called **monitor**.

The monitor is responsible to coordinate and run the accountability algorithm as soon as a fork is detected. 
Since the input of the algorithm are the **message logs** of the validators (i.e., all the information related to the exchanged messages by processes in the consensus instance), the monitor will first request the message logs from each validator and then will run the accountability algorithm.

We assume that all validators keep track of all the messages they sent and received during an execution of the Tendermint consensus algorithm. 
The idea is that a correct process would be able to prove that it was not faulty by showing its activity (*message logs*) to the monitor.
Indeed, the monitor should be able to detect all faulty processes that were responsible to generate a fork by analyzing the message logs from all the correct processes that have no reason for not sending their message logs.

## Fork scenarios

In order to better understand the problem, we analyze the scenarios when the fork can happen and how the monitor would be able to handle the situation when analyzing the log received from each validator.

From now on, the expression "leading to a fork" will be used to generally define a set of actions carried out by a process that caused a fork. 
For example, saying "a process *p* leads to a fork" means that *p* caused a fork by having made one or more faulty behaviours in a specific consensus instance.

### Fork in a single round

A fork happens in a certain round *r* when two decisions are made in the same round *r* for different values.
In other words, this means that a correct process *p* had decided value *v* in round *r* and some other correct process *p'* (*p != p'*) has decided value *v'* (*v != v'*) in round *r*. 

Given the fact that all the messages that led to a fork have been exchanged in round *r*, at least *f + 1* processes must have sent a PREVOTE/PRECOMMIT message for both *v* and *v'*.
Since at least one correct process will receive the messages for both values, the monitor would be able to eventually identify the faulty processes that led to a fork.  

### Fork across different rounds

A fork happens in different rounds when two decisions for different values are made in different rounds. 
In other words, this means that a correct process *p* had decided value *v* in round *r* and some other correct process *p'* (*p != p'*) has decided value *v'* (*v != v'*) in round *r'* such that *r' > r*. 

Given this scenario, it can be noticed that at least *f + 1* processes have sent a PRECOMMIT message for value *v* in round *r* and also sent a PREVOTE message for value *v'* in round *r'* despite having locked value *v* when they sent a PRECOMMIT.

These processes are faulty unless they have a valid justification for sending the PREVOTE message for *v'* after being locked on *v* (i.e., other *2f + 1* PREVOTE messages for value *v'* in a round *r''* s.t. *r < r'' < r'*).
 
On the other hand, if they do have such a justification, it means that there are at least *f + 1* processes that sent a PRECOMMIT message for value *v* in round *r* and also sent a PREVOTE message for value *v'* in round *r''*, despite having locked value *v* when they sent the PRECOMMIT of round *r*.

It is clear that this scenario is similar to the one described at the beginning: the processes that sent a PREVOTE message for *v'* in *r''* are faulty unless they have a valid justification for sending the PREVOTE message for *v'* after being locked on *v*.

The monitor would be able to track down the origin of the misbehaviour by checking recursively that all the processes that sent a PREVOTE message for a value different from the one they were locked on have a valid justification.
Since messages cannot be generated out of nowhere, there must be an iteration of this recursive problem where the monitor would be able to find at least *f + 1* processes that sent a PREVOTE message without a valid justification and led to a fork.

In conclusion, only two possible options remain:

1. at least *f + 1* processes have sent both PREVOTE messages for value *v* and *v'* in the same round *r*

2. at least *f + 1* processes have sent a PRECOMMIT message for value *v* in round *r* and also sent a PREVOTE message for value *v'* in a round *r''* (such that *r < *r'' <= r'*) without having delivered *2f +1* PREVOTE messages for value *v'* in a round *r'''* (where *r <= r''' < r''*)

## Challenges

It is possible that faulty processes send invalid message logs to the monitor. An invalid message log is a message log which contains partial or incomplete information regarding the actual messages sent and received by the owner of the log (this means some messages might have not been reported or listed in the log).
Indeed, a faulty process is able to alter and modify its message log in order to hide a possible misbehavior. In that case, the monitor might not have enough information to detect all the faulty processes that led to a fork. 

We assume that messages exchanged in the consensus algorithm are signed and it is not possible to forge messages (i.e., some process *p* cannot state that it has received a message from process *p'* if *p'* did not send that message). 

Therefore, the following scenarios are possible:
 
1. A process *p* denies having received a message *m* (*p*'s message log does not contain *m* as a message received)  

2. A process *p* denies having sent a message *m* (*p*'s message log does not contain *m* as a message sent) 

The first case will be ignored because it does not contribute to the generation of a fork.

The second point is, instead, more important in this discussion because a process can deny having sent a message *m* in order to hide its misbehavior. 
However, if *m* led to a fork, at least one process *p'* (*p != p'*) received *m*: the monitor can add this missing message in *p*'s message log without needing to take any further action, simply correcting the invalid message log received. 

From now on, the expression "inferring a message log" will be used to define the action carried out by the monitor to correct an invalid message log with the missing messages taken from other message logs received by *p*.

To summarize, the monitor can infer the valid message logs for each process after receiving the message logs from all the correct processes, which are guaranteed to arrive and to be correct.
Once all the message logs received are adjusted and corrected accordingly, the monitor can start analyzing each message log and determine whether the process owner of the message log is faulty or not.

### Communication between the monitor and validators

The communication between the monitor and the validators for receiving message logs of the Tendermint consensus algorithm is a crucial aspect of the monitor execution.

The monitor needs to make sure it will receive the message logs from the correct processes and, depending on how the communication is modelled, some assumptions must be made in order to design a correct fork accountability algorithm.

#### Synchronous model

In the case of a synchronous model, the monitor should wait a certain amount of time (*timeout*) to receive all the message logs before running the accountability algorithm.

After the *timeout*, if some process *p* did not send its message log, the monitor can only consider *p* faulty due to the nature of the communication, even though it does not have all the information to determine its misbehavior with certainty. 
On the other hand, by not considering *p* faulty, the monitor might not be able to find at least *f + 1* faulty processes that led to a fork.

#### Asynchronous model

The approach given above does not work in an asynchronous case, when there are no assumptions regarding the communication between the monitor and the validators. 

In fact, the monitor cannot expose some process *p* if it did not receive *p*'s message log: the monitor could mistakenly detect a correct process because it does not have enough information to determine if *p* had a valid reason to send a PREVOTE/PRECOMMIT message (in other words, accuracy property would be violated).
 
By receiving the message logs from other processes, the monitor can only determine if *p* equivocated. 
However, while analyzing the fork scenarios in the previous section, we have shown that faulty processes can lead to a fork by not equivocating (point 2 in the section "Fork in different rounds") and this makes impossible the design of an asynchronous accountability algorithm in the current Tendermint consensus algorithm.  

In order to design a correct asynchronous accountability algorithm, it is necessary to slightly modify the Tendermint consensus algorithm so that the accuracy property is not violated.

The critical case outlined above is when a process sends a PRECOMMIT message for a value *v* and then, in a later round, it sends a PREVOTE message for a different value than *v*. In this way, a faulty process can lead to a fork without making equivocation.

An asynchronous accountability algorithm would work correctly if processes can justify the sending of a PREVOTE message directly: a simple solution would be attaching the justifications to the PREVOTE message itself so that the accountability algorithm can directly check if the message has been correctly sent.
In this case, the monitor would be able to complete the accountability algorithm correctly (respecting both the completeness and accuracy properties) as soon as at least *f + 1* faulty processes will be detected and by relying on the fact that at least *f + 1* correct processes will send correct message logs.

## Algorithm design
 
Assuming the small modification in the Tendermint consensus algorothm described above, these are the high-level steps carried out by the monitor to run the asynchronous accountability algorithm:

1. Upon a fork, monitor requests message logs from all the validators

2. Monitor runs the accountability algorithm to detect faulty processes upon receiving at least *f + 1* message logs and at every new message logs received when there are at least *f + 1* message logs 

3. While running the accountability algorithm, the monitor analyzes the received logs: for every message *m* that has been received by some other process but is not present in the sender's message log, the monitor attaches *m* to the message log of the sender of *m*

4. For each validator, the monitor scans the received message logs and the inferred message logs and determines for each process whether it is faulty or not by analyzing the history of sent messages and transitions

5. When the monitor detects at least *f + 1* faulty processes, the algorithm completes

Termination is guaranteed by the fact that correct processes will send their messages logs and their message logs will be correct (no sent message will be missing). 
If this condition does not hold, the algorithm would not be able to identify correctly all the faulty processes that led to a fork. However, correct processes have no reason to misbehave or send invalid message logs. 

## Notes

- As some readers might have noticed, the fork accountability problem can be easily solved by simply piggybacking justifications in every single message. 
However, this solution would be inefficient because it considerably increases the amount of information stored and exchanged by processes. The scope of the discussion is to solve the problem efficiently with a minimal amount of changes required on the Tendermint consensus protocol and that is the reason why we do not consider and analyze this possible solution.   

- By excluding the piggybacking solution of the previous point, the fork accountability cannot be solved when the number of faulty proccesses is greater than *2f*: in this case, faulty processes basically control the system and nothing can be done to solve the problem.

- The accountability algorithm is able to also catch other faulty behaviours that did not necessarily lead to a fork - this is a positive side effect of the accountability algorithm described above.

- The monitor needs to analyze the message logs from the round of the first decision to the round of the second decision. Previous message history is not necessary to make the algorithm work. 

- The accountability algorithm runs for a single instance of consensus (i.e., for a single height). However, it can be easily configured to support multiple heights by running the same single-height algorithm on different heights.  
