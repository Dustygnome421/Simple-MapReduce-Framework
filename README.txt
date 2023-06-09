CS4513: The MapReduce Library
=======================================

Note, this document includes a number of design questions that can help your implementation. We highly recommend that you answer each design question **before** attempting the corresponding implementation.
These questions will help you design and plan your implementation and guide you towards the resources you need.
Finally, if you are unsure how to start the project, we recommend you visit office hours for some guidance on these questions before attempting to implement this project.


Team members
-----------------

1. Aman Gupta (agupta9@wpi.edu)
2. Oliver Reera (omreera@wpi.edu)

Design Questions
------------------

(2 point) 1. If there are n input files, and nReduce number of reduce tasks , how does the the MapReduce Library uniquely name the intermediate files?

    The MapReduce Library names the intermediate files with the format: fi-0, fi-1,..., fi-[nReduce-1] for the ith map task.


(1 point) 2. Following the previous question, for the reduce task r, what are the names of files will it work on?

    For the rth reduce task, the names of files it will work on have the format: f0-r, f1-r,..., f[n-1]-r.


(1 point) 3. If the submitted mapreduce job name is "test", what will be the final output file's name?

    The final output file's name will be "mrtmp.test".


(2 point) 4. Based on `mapreduce/test_test.go`, when you run the `TestBasic()` function, how many master and workers will be started? And what are their respective addresses and their naming schemes?

    There will be one master and two workers. The master's address and naming scheme will be "var/tmp/824-[user id]/mr[process id]-master" and the workers' addresses and naming schemes will be "var/tmp/824-[user id]/mr[process id]-worker0" and "var/tmp/824-[user id]/mr[process id]-worker1".


(4 point) 5. In real-world deployments, when giving a mapreduce job, we often start master and workers on different machines (physical or virtual). Describe briefly the protocol that allows master and workers be aware of each other's existence, and subsequently start working together on completing the mapreduce job. Your description should be grounded on the RPC communications.

    The protocol that allows master and workers to be aware of each other's existence is TCP. When the master and workers communicate with RPC over a network, they must use the TCP protocol rather than UNIX-domain sockets as when they are on the same machine.


(2 point) 6. The current design and implementation uses a number of RPC methods. Can you find out all the RPCs and list their signatures? Briefly describe the criteria a method needs to satisfy to be considered a RPC method. (Hint: you can look up at: https://golang.org/pkg/net/rpc/)

    The RPCs consist of: Shutdown(_ *struct{}, res *ShutdownReply), Register(args *RegisterArgs, _ *struct{}), and DoTask(arg *DoTaskArgs, _ *struct{}). From golang.org, the requirements for a method to be considered an RPC method are: the method's type is exported, the method is exported, the method has two arguments, both exported (or builtin) types, the method's second argument is a pointer, and the method has return type error.


Errata
------

Describe any known errors, bugs, or deviations from the requirements.

---
