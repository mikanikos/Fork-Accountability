# Fork accountability solution

## Main abstractions

In order to simplify the later discussion, here we introduce some abstractions that help to better identify and classify the message logs received and handled by the monitor.

Each **Message** contains the following information:

- **Type** of the message, *type*

- **Unique identifier** of the sender of the message, *id*

- **Round** of the message, *r*

- **Value** of the message, *v*

A **Vote Set** is the set of all Messages that a process sent and received in a specific round *r*.

A **Height Vote Set** is the set of all Vote Set in a specific height *h* or the set of all Messages that a process sent and received in a specific height *h*.

Tendermint messages have several other important fields that are used by the consensus protocol for other purposes. However, we limit the scope of our discussion to these information only because they are essential for running the designed accountability algorithm.       

## Monitor algorithm
 
These are the detailed steps of the execution of the monitor to run the asynchronous accountability algorithm: 
 
1. Upon a fork, monitor sends a request for validators' height vote sets (following a specific communication protocol) for a given height. Monitor starts waiting for incoming message logs.

2. Validators receive the request from the trusted monitor and they send whatever they have to the monitor immediately. 

3. Monitor runs the accountability algorithm upon receiving a new message logs but only when **the total number of different height vote sets received is at least *f + 1***. 
If the threshold is met, the monitor runs the accountability algorithm. Otherwise, the monitor keeps waiting for other packets from the processes that did not reply yet. 
      
4. The accountability algorithm runs in two consecutive phases:

    4. **Pre-processing phase**: given the received message logs so far, the monitor "infers" the missing sent messages and attaches them to the original sender's height vote set.

    4. **Fault-detection phase**: the monitor scans the height vote set from the first round to the last round and checks that the rules of the Tendermint consensus algorithm are violated. If a process violates any of the Tendermint consensus algorithm rules, it is detected as faulty. 

    The output of the accountability algorithm is the list of processes that have been detected as faulty and the proof of their misbehavior.
  
5. If the monitor detected **at least *f + 1* faulty processes** during the last execution of the accountability algorithm, the monitor completes. Otherwise, it keeps waiting for more height vote sets.

## Accountability algorithm

### Pre-processing phase
During this phase, the monitor analyzes all the height vote sets received. 
For each height vote set, the monitor scans all the received vote sets. 
For every message *m* in the received vote set, the monitor checks if the message is present in the set of message sent by the sender of *m*. If it is not present, the monitor adds it to the sender's vote set.   

### Fault-detection phase
During this phase the monitor analyzes each height vote set after the pre-processing phase and determines whether the process is faulty or not.
For every height vote set received, the monitor goes through all the sent messages in each round r from the first round to the last round. The monitor checks for the following faulty behaviours:

- the process equivocated in round r (sent more than one PREVOTE or PRECOMMIT message)

- the process sent a PRECOMMIT message but did not receive *2f + 1* valid PREVOTE messages to justify the sending of the PRECOMMIT message

- the process sent a PREVOTE message and did not have *2f + 1* valid PREVOTE messages as justification inside the PREVOTE

If one or more of these faulty behaviours is found in any of the rounds analyzed, the process is detected as faulty.

Please note that height vote sets that have not been received (therefore, they have been inferred during the pre-processing phase) will be checked for equivocation only.

## Pseudo-code implementation

We now present a simple pseudo-code version of the algorithm in order to make readers better understand some implementation details.
The following algorithm is only a high-level overview of the main steps of the accountability solution implemented, we invite the reader to check out the code for further details.  

#### Monitor algorithm

The monitor has the following interface:

- **runMonitor(V, h, firstRound, lastRound)**: run the accountability algorithm on the set of validators V for height *h* from *firstRound* to *lastRound*, return the list of faulty processes in height *h*

In order to simplify the understanding of the algorithm, the algorithm uses the following high-level methods to model the communication between the monitor and validators:

- **sendRequest(v, h)**: send request for the height vote set of height h to validator v 

- **deliverHVS()**: deliver the next incoming height vote set sent by some validator, return the next height vote set received

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

#### Accountability algorithm

The following is the accountability algorithm used by the monitor to detect faulty processes based on the information given.

The accountability algorithm has the following interface:

- **runAccountability(hvsDelivered, firstRound, lastRound)**: run accountability algorithm on the height vote sets list *hvsDelivered* from *firstRound* to *lastRound*, return the list of faulty processes detected based on the parameters

The algorithm uses the following high-level methods for simplifying the presentation:

- **getHVSFromSender(hvsDelivered, id)**: get the height vote set corresponding to the sender *id* given, return *nil* otherwise 

- **newHvs(id)**: create a new height vote set given an *id*, return the newly-created inferred height vote set

- **getVoteSetFromRound(hvs, r)**: get the vote set from the given height vote set *hvs* of round *r*, return the vote set relative to the height vote set *hvs* given
    
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
                if vs.sent contains more than one PREVOTE or PRECOMMIT in round r:
                    faultyProcesses = faultyProcesses + hvs.sender                
                
                # if the hvs was not inferred, check for correctness of the execution in round r according to tendermint consenus algorithm rules  
                if !hvs.inferred:
                    
                    # check if process sent a PRECOMMIT or a PREVOTE without proper justiifcation (PREVOTE messages must contain valid justifications) 
                    if vs.sent contains a PRECOMMIT or PREVOTE without a valid justification:
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
