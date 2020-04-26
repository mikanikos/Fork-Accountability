# Fork accountability

As we know, the agreement property of the Tendermint consensus algorithm cannot be violated if there are at most f faulty processes in the system.
What happen if more than f processes are faulty?

A well designed consensus protocol should provide some guarantees when this limit is exceeded at a certain limit. The most important guarantee is fork-accountability, where the processes that caused the consensus to fail can be identified and punished. 

We can separate two cases:

- if the number of faulty processes is more than 2f, we cannot do anything because faulty processes basically control the system and they can do whatever they prefer.
- if the number of faulty processes is less than 2f, i.e.

      f < num. faulty processes <= 2f
      
  Agreement can still be violated but we're in condition to detect the processes that were faulty and punish them, for example by removing them from the validators set. 
  
A fork happens when two correct validators decide on different blocks in the same height, i.e. there are two commits for different blocks at the same height. We aim to give incentives to validators to behave correctly and according to the protocol specifications and we want to detect detect faulty validators and not mistakenly detect correct validators. 
 
There are multiple reasons that can lead to a fork, all coming from processes that deviate from the protocol specification and not following the behavior of a correct process.

We can summarize the most important ones:

- A process sends multiple proposals in a round r for different values
- A process sends a Prevote/Precommit message for an invalid value v in a round r
- A process sends multiple Prevote/Precommit messages in a round r for different values
- A process sends a Prevote/Precommit message for a value v despite gaving a lock on a different value v'
- A process sends a Precommit for a value v without having received at least 2/3 Prevote messages for v 

Given this problem, we want to develop an accountability algorithm that, by analyzing the messages logs from each validator, is able to determine which processes respected the protocol and which ones didn't. In order to prove the misbehavior of the processes that didn't respect the protocol, we want to show the proof of its faultiness (see reasons above) respect to the Tendermint consensus protocol.