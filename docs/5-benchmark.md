# Benchmark and evaluation

The environment where the experiments have been carried out is a cloud virtual machine instance on [Google Cloud](https://cloud.google.com/), located in the zone us-central1-a. 
The machine type is a [n1-standard-16](https://cloud.google.com/compute/docs/machine-types#n1_standard_machine_types) (16 vCPUs, 60 GB memory), running the latest Ubuntu 20.04 LTS. The main hardware specifications and details are reported in Table.

The reported benchmarks consist of several executions of the monitor algorithm, both asynchronous and synchronous version, in different fork scenarios and interacting with different configurations of validators.
The fork scenarios are always due to validators sending a PRECOMMIT message without enough PREVOTE messages and have been automatically generated with a script that allows creating valid configuration files for both the monitor and the validators with user-defined parameters.

The parameters chosen for the benchamarks are the number of validators and the number of rounds between the first and the second decision in the fork. 
The Table summarizes the configurations used for the experiments.

The metrics used to evaluate the performance are the following:
- Execution time of the last (and successful) iteration of the accountability algorithm only
- Total elapsed time between the start and the end of the entire command used to execute the monitor algorithm (it includes the whole accountability algorithm plus connection and packet exchange with validators)
- CPU percentage utilization of the command used to execute the monitor algorithm
- Maximum [resident set size (RSS)](https://en.wikipedia.org/wiki/Resident_set_size) allocated for the process running the command used to execute the monitor algorithm

The last three metrics have been computed using the [time command](https://en.wikipedia.org/wiki/Time_%28Unix%29) in its [extended version](https://man7.org/linux/man-pages/man1/time.1.html).

The benchamrk experiments have been repeated three times with very similar results, therefore the obtained data can be considered moderately valid and reliable.
The file containing all the results of the experiments can be found at [benchmark](/benchmarks/benchmark_report_final.txt). 


## Total execution time
As we can see, the amount of time taken to run the entire monitor algorithm is relatively small in almost all cases. 
As expected, a higher number of both validators and rounds increases the time to complete the algorithm: there are not only more iterations (i.e., more messages) to execute in the accountability algorithm but there is also the overhead due to the communication between monitor and validators.
Even though validators have been set up to send back their message logs as soon as they receive the request, the communication protocol (which includes the connection establishment with each validator, the time to send/receive a packet to/from another process etc.) has an impact on the data analyzed. 
This is also the reason why we notice that the number of validators has a greater impact on the time while the number of rounds minimally affects the results, especially for lower values.

We can notice that the synchronous version has a slightly better performance overall than the asynchronous one: this is likely the result of the several executions of the accountability algorithm that must be run every time a new message log from a validator is delivered and there are at least *f + 1* message logs, as long as at least *f + 1* faulty processes are found. 
On the other hand, the synchronous version runs once when the timer expires or when all the expected message log are received.   


## Execution time of the last iteration of the accountability algorithm
In order to better evaluate the accountability algorithm without the communication overhead, we measured the amount of time taken by the accountability algorithm during the last iteration (i.e., when the accountability algorithm completed successfully after receiving all the necessary message logs to detect the faulty processes).
Clearly, the accountability algorithm in the synchronous version always ran once, while we analyzed only the last execution of the accountability algorithm in the asynchronous version.

The outcome is similar to the one described in the previous section, where the higher number of validators and rounds makes the algorithm naturally a bit slower.
However, the number of validators has a more evident impact on the increase of the execution time, while the number of rounds continues to generally make a smaller difference unless the number is very high.

In this measurement the asynchronous case offers better results respect to the synchronous case because there are, in most cases, less message logs to analyze during the last execution.
in fact, the synchronous mode requires to run the algorithm when all messages logs have been delivered (or when the timeout has expired, but this case is not taken into consideration). 
On the other hand, the asynchronous version usually runs the last iteration of the accountability algorithm with a smaller number of message logs and, even though it carries out additional required computations (e.g. inferring the missing height vote sets, checking for justifications etc.) is a bit faster than the synchronous version, excluding the communication overhead.

## CPU and memory utilization
Regarding the more technical metrics, the algorithm runs well without over-using too many computational resources.
CPU utilization has been measured by analyzing the percentage of CPU (respect to multiple cores) used by the monitor algorithm process.
Memory utilization has been measured by analyzing the amount memory allocated to the monitor algorithm process.

Even under more "stressful" conditions, the algorithm does not require more than four cores and 500 megabytes to run.  
Apparently, the asynchronous mode requires slightly more cores but less memory to run. However, the results are machine-specific and, as we can notice, they are not completely uniform across all the experiments carried out. 
Therefore, it is difficult to extract an exact trend for this data. However, we can say that for relatively high computations the algorithm works efficiently with a small amount of resources and is not expensive for limited machines.

