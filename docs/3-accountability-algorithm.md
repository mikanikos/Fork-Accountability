# Accountability algorithm

## Main abstractions

In order to simplify the discussion we introduce some abstractions to better identify and classify the message logs received and handled by the monitor.

Each **Message** contains the following information:

- **Type** of the message, *type*
- **Unique identifier** of the sender of the message, *id*
- **Round** of the message, *r*
- **Value** of the message, *v*

A **Vote Set** is the set of all Messages that a process sent and received in a specific round *r*.

A **Height Vote Set** is the set of all Vote Set in a specific height *h* or the set of all Messages that a process sent and received in a specific height *h*.

Tendermint messages have lots of other important information that are used by the consensus protocol but we limit our discussion to the necessary data to run the accountability algorithm.       

## Execution steps
 
We show the detailed steps of the execution of the monitor and the accountability algorithm: 
 
1. Upon a fork, monitor sends a request for validators' height vote sets (following a specific communication protocol) for a given height. Monitor starts waiting for incoming message logs.

2. Validators receives the request from the trusted monitor. If they have the requested height vote set, they send it immediately. 
If they don't have any logs for the requested height, they can simply ignore the request or, alternatively, send a reply to the monitor to notify they do not have any logs for that height.

3. Monitor runs the accountability algorithm upon receiving a new message logs but only when **the total number of different height vote sets received is at least f + 1**. 
If the threshold is met, the monitor runs the accountability algorithm. Otherwise, the monitor keeps waiting for other packets from the processes that did not reply yet. 
      
4. The accountability algorithm itself has two consecutive steps:

    4. **Pre-processing phase**: given the received message logs so far, we "infer" the missing sent messages and we attached them to the original sender's height vote set.

    4. **Fault-detection phase**: we scan the height vote set from the first round to the last round and we check if the process is faulty by applying the rules of the Tendermint consensus algorithm. 

  The output of the accountability algorithm is the list of processes that have been detected as faulty and the proof of their faultiness.
  
5. If the monitor detected **at least f + 1 faulty processes** during the last execution of the accountability algorithm, the monitor completes. Otherwise, it keeps waiting for more height vote sets.

## Accountability algorithm

### Pre-processing phase
During this phase, we analyze all the height vote sets received. For each height vote set, we analyze all the received vote sets. For every message m in the received vote set, we check if the message is present in the sent vote set of the sender of m. If it is not present, we add it.   

### Fault-detection phase
During this phase we analyze each height vote set after the pre-processing phase and we determine whether the process is faulty or not.
For every height vote, we go through all the sent messages in each round r from the first round to the last round. We check for the following faulty behaviours:

- the process equivocated in round r (sent more than one Prevote or Precommit message) - height vote sets that have not been received will be checked for equivocation only

- the process sent a Precommit message but did not receive *2f + 1* valid Prevote messages to justify the sending of the Precommit message

- the process sent a Prevote message and did not have *2f + 1* valid Prevote messages as justification inside the Prevote

If one or more of these faulty behaviours is found in any of the rounds analyzed, the process is detected as faulty.

### Pseudo-code version

We now present a simple pseudo-code version of the algorithm in order to make the reader better understand some specific details.
The following algorithm is only a high-level overview of the main steps of the accountability processing, we invite the reader to check out the code for further details.  

The monitor has the following interface:

- **runMonitor(V, h, firstRound, lastRound)**: run the accountability algorithm on the set of validators V for height h from firstRound to lastRound, returns the list of faulty processes in height h

In order to simplify the understanding of the algorithm, we use the following high-level methods to model the communication between monitor and validators:

- **sendRequest(v, h)**: send request for the height vote set of height h to validator v 
- **deliverHVS()**: deliver the next incoming height vote set sent by some validator, returns the height vote set received

    ```
    # Monitor algorithm
    runMonitor(V, h, firstRound, lastRound):
        
        // compute max number of faulty processes with respect to the number of validators  
        f = (V.length - 1) / 3
        
        # send request to validators
        for v in V:
            sendRequest(v, h)
        
        # store delivered height vote sets
        hvsDelivered = []
        
        # wait to deliver f+1 different height vote sets
        do:
            hvs = deliverHVS()
            hvsDelivered = hvsDelivered + hvs
        while (hvsDelivered.length < f+1)
        
        # run accountability algorithm
        faultyProcesses = runAccountability(hvsDelivered, firstRound, lastRound)
        
        # repeat until we find at least f+1 faultyProcesses
        while (faultyProcesses.length < f+1):
        
            # wait to deliver more HVS
            hvs = deliverHVS()
            hvsDelivered = hvsDelivered + hvs
        
            # run accountability algorithm 
            faultyProcesses = runAccountability(hvsDelivered, firstRound, lastRound)    
           
            
        # return the final output of the algorithm that satisfied exit condition
        return faultyProcesses
    ```


The following is the accountability algorithm used by the monitor to detect faulty processes based on the information given.

The accountability algorithm has the following interface:

- **runAccountability(hvsDelivered, firstRound, lastRound)**: run accountability algorithm on the height vote sets list hvsDelivered from firstRound to lastRound, returns the list of faulty processes detected based on the parameters

It uses the following high-level methods for simplifying the understanding of the algorithm:

- **getHVSFromSender(hvsDelivered, id)**: get the height vote set corresponding to the sender id given, returns nil otherwise 
- **newHvs(id)**: create new height vote set given an id, returns the newly-created inferred height vote set
- **getVoteSetFromRound(hvs, r)**: get the vote set from the height vote set given hvs and the round r, returns the vote set relative of the height vote set hvs given
    
    ```
    # Accountability algorithm
    runAccountability(hvsDelivered, firstRound, lastRound):
    
        # Pre-processing phase
                
        # go through all the height vote sets received
        for each hvs in hvsDelivered:
            # go through all the messages received
            for each m in hvs.received:
                # get height vote set of the sender of the message m
                hvsSender = getHVSFromSender(hvsDelivered, m.sender)
                
                # if not present in the height vote sets received, create it 
                if hvsSender == nil:
                    hvsSender = newHvs(m.sender)
                    # add inferred flag for later processing
                    hvsSender.inferred = true
                    hvsDelivered = hvsDelivered + hvsSender 
                
                # if m is not present in the hvs of the sender, add it
                if m is not in hvsSender:
                    hvsSender = hvsSender + m
                    
        # Fault-detection phase
        
        # set of faulty processes detected, duplicates are discarded
        faultyProcesses = []
        
        # go through all the height vote sets received
        for each hvs in hvsDelivered:
            # from the first to the last round    
            for each round r from firstRound to lastRound: 
                # get vote set from round
                vs = getVoteSetFromRound(hvs, r)
                
                # check for equivocation
                if vs.sent contains more than one Prevote or Precommit in round r:
                    faultyProcesses = faultyProcesses + hvs.sender                
                
                # if the hvs was not inferred, check for correctness of the execution in round r according to tendermint consenus algorithm rules  
                if !hvs.inferred:
                    
                    # check if process sent a precommit or a prevote without proper justiifcation (prevotes must contain valid justifications) 
                    if vs.sent contains a precommit or prevote without a valid justification:
                        faultyProcesses = faultyProcesses + hvs.sender
                        
                        
        return faultyProcesses
    ```        
                
## Implementation-specific details

- The communication between the monitor and the validators is over TCP and is structured as a normal client-server interaction.  

- The monitor checks the received responses for validity and does not accept a height vote set if this is not valid and will keep waiting for a valid response from validators. 
The monitor also keeps track of the received responses: if a validator sent a valid height vote set, the monitor will stop waiting for a response and closes the connection.
If some implementation failure is detected by the monitor (network failure, validator crashes before sending height vote set etc.) the monitor will know that it will not receive a response from the failed validator. 
Therefore, the monitor will stop the execution if a response is not expected from any other process. 

- For the safety of the execution, the monitor has a timeout that aims to prevent a "wait-forever" state in case of some problems (network issue, invalid hvs, error in config files etc.). The timeout is not necessary from a theoretical point of view because the assumptions made guarantee that the accountability algorithm will eventually complete.
