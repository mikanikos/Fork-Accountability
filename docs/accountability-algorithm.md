# Accountability algorithm

## Main abstractions

In order to simplify the discussion we introduce some abstractions to better identify and classify the message logs received and handled by the monitor.

Each **Message** contains the following information:

- **Type** of the message, *type*
- **Unique identifier** of the sender of the message, *id*
- **Round** of the message, *r*
- **Value** of the message, *v*

A **Vote Set** is the set of all Messages that a process sent and received in a specific round *r*.

An **Height Vote Set** is the set of all Vote Set in a specific height *h*, i.e. the set of all Messages that a process sent and received in a specific height *h*.

Tendermint messages have lots of other important information that are used by the consensus protocol but we limit our discussion to the necessary data to run the accountability algorithm.       

## Execution steps
 
We show the detailed steps of the execution of the monitor and the accountability algorithm: 
 
1. Upon a fork, monitor sends a request for validators' height vote sets (following a specific communication protocol) for a given height. Monitor starts waiting for incoming message logs.

2. Validators receives the request from the trusted monitor. If they have the requested height vote set, they send it immediately. 
If they don't have any logs for the requested height, they can simply ignore the request or, alternatively, send a reply back to the monitor to notify they don't have any logs for that height.

3. Monitor runs the accountability algorithm upon receiving a new message logs but only when **the total number of different height vote sets received is at least f + 1**. 
If the threshold is met, the monitor runs the accountability algorithm. Otherwise, the monitor keeps waiting for other packets from the processes that didn't reply back yet. 
      
4. The accountability algorithm itself has two consecutive steps:

    4. **Pre-processing phase**: given the received message logs so far, we "infer" the missing sent messages and we attached them to the original sender's height vote set.

    4. **Fault-detection phase**: we scan the height vote set from the first round to the last round and we check if the process is faulty by applying the rules of the Tendermint consensus algorithm. 

  The output of the accountability algorithm is the list of processes that have been detected as faulty and the proof of their faultiness.
  
5. If the monitor detected **at least f + 1 faulty processes** during the last execution of the accountability algorithm, the monitor completes. Otherwise, it keeps waiting for more height vote sets.

## Accountability algorithm steps

### Pre-processing phase
During this phase, we analyze all the height vote sets received. For each height vote set, we analyze all the received vote sets. For every message m in the received vote set, we check if the message is present in the sent vote set of the sender of m. If it's not present, we add it.   

### Fault-detection phase
During this phase we analyze each height vote set after the pre-processing phase and we determine whether the process is faulty or not.
For every height vote, we go through all the sent messages in each round r from the first round to the last round. We check for the following faulty behaviours:

- the process equivocated in round r (sent more than one Prevote or Precommit message) - height vote sets that have not been received will be checked for equivocation only

- the process sent a Precommit message but didn't receive *2f + 1* valid Prevote messages to justify the sending of the Precommit message

- the process sent a Prevote message and didn't have *2f + 1* valid Prevote messages as justification inside the Prevote

If one or more of these faulty behaviours is found in any of the rounds analyzed, the process is detected as faulty.

## Implementation-specific details

- The communication between the monitor and the validators is over TCP and is structured as a normal client-server interaction.  

- The monitor checks the received responses for validity and doesn't accept an height vote set if this is not valid and will keep waiting for a valid response from validators. 
The monitor also keeps track of the received responses: if a validator sent a valid height vote set, the monitor will stop waiting for a response and closes the connection.
If some implementation failure is detected by the monitor (network failure, validator crashes before sending height vote set etc.) the monitor will know that it will not receive a response from the failed validator. 
Therefore, the monitor will stop the execution if a response is not expected from any other process. 

- For the safety of the execution, the monitor has a timeout that aims to prevent a "wait-forever" state in case of some problems (network issue, invalid hvs, error in config files etc.). The timeout is not necessary from a theoretical point of view because the assumptions made guarantee that the accountability algorithm will eventually complete.
