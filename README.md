# Instructions

1. Install the go module on your machine. Run the command:
   - `go install tcp-network`
2. Run the executable. Run the command providing an integer for `X`:
   - `$HOME/go/bin/tcp-network X`

Instruction for how to develop, use, and test the code
information on what the problem was, why it needed to be solved, or what were the consequences of picking that particular solution.

# Documentation
# Product Requirements

MP1 tasked us with designing a simple network simulation, with the added caveat that the messages we send might be delivered out of order on the destination process. According to the [design spec](https://docs.google.com/document/d/1qLuygCkNm5WbI_a-LBhVs95_BlTyEbFuOeM_CtxrDV0/), we were to read the delay from a configuration file organized like so,

`min_delay(ms) max_delay(ms)
ID1 IP1 port1
ID2 IP2 port2
... ... ...`

The simulated network delay described above is bounded by values in the configuration file. The first line of the configuration file has two values, a minimum delay and a maximum delay. We are to randomize the delay between these two values so as to have an unpredictable delay between when messages are sent, and messages are recived. We were not allowed to put the program to sleep to simulate the delay. This meant we could not simply pass messages through channels, because the messages would always arrive in order due to the FIFO nature of go channels. We could not sleep entire functions either, because messages would still arrive in order regardless of the sleep duration.

The other specifications for MP1 included implementing unicast functionality using TCP communication, contained within two functions - unicast_send and unicast_receive. 

`<b>unicast_send(destination, message):</b> sends a message to the destination process`

`<b>unicast_receive(source, message):</b> delivers the message received from the source process`

To send messages, we are to build a terminal interface where a message can be sent from instance to instance. The pattern of these messages is, for example, `send 2 Hello`, more generally the pattern is, `send {ID} {message}`. In the first example, "send" is the keyword to send a message, "2" is the ID of the instance we are sending to and "Hello" is our message.

When a process sends or receives a message, it should write a timestamp to the terminal that includes the time the message was sent/received and the content of the message.

# Proposed Design
We started by breaking up the program into 4 processes.
    - The server process, stored in pkg/server/
    - The client process, stored in pkg/client/
    - The config file and the program to read it, stored in pkg/config/
    - And main.go

We delegated separate responsibilities to each process. Server contains the goroutines and tcp threads that send and recieve messages between server instances. Client reads input from the terminal and forwards it to its corresponding server instance to send the message. Config reads data from the configuration file to assign the min and max delay, as well as ID's, IP's, and Ports to each server. And main, which acts as an entry point to the program and wraps all the separate processes in one file.

# Technical Design

// PUT PROGRAM FLOW DIAGRAM HERE

# Implementation
