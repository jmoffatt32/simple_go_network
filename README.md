# Instructions

1. Install the go module on your machine. Run the command:
   - `go install tcp-network`
2. Run the executable. Run the command providing an integer for `ID`:
   - `$HOME/go/bin/tcp-network {ID}`
3. You will have to start at least two processes (in two different terminals) to send messages between instances
4. You'll be presented with a command prompt, enter commands to send messages in the form,
   - `send {ID} {message}` where ID is the number of server instance that you have started, and message is a message you would like to send

# Documentation

# Product Requirements

MP1 tasked us with designing a simple network simulation, with the added caveat that the messages we send might be delivered out of order on the destination process. According to the [design spec](https://docs.google.com/document/d/1qLuygCkNm5WbI_a-LBhVs95_BlTyEbFuOeM_CtxrDV0/), we were to read the delay from a configuration file organized like so,

```
min_delay(ms) max_delay(ms)
ID1 IP1 port1
ID2 IP2 port2
... ... ...
```

The simulated network delay described above is bounded by values in the configuration file. The first line of the configuration file has two values, a minimum delay and a maximum delay. We are to randomize the delay between these two values so as to have an unpredictable delay between when messages are sent, and messages are recived. We were not allowed to put the program to sleep to simulate the delay. This meant we could not simply pass messages through channels, because the messages would always arrive in order due to the FIFO nature of go channels. We could not sleep entire functions either, because messages would still arrive in order regardless of the sleep duration.

The other specifications for MP1 included implementing unicast functionality using TCP communication, contained within two functions - unicast_send and unicast_receive.

`unicast_send(destination, message): sends a message to the destination process`

`unicast_receive(source, message): delivers the message received from the source process`

To send messages, we are to build a terminal interface where a message can be sent from instance to instance. The pattern of these messages is, for example, `send 2 Hello`, more generally the pattern is, `send {ID} {message}`. In the first example, "send" is the keyword to send a message, "2" is the ID of the instance we are sending to and "Hello" is our message.

When a process sends or receives a message, it should write a timestamp to the terminal that includes the time the message was sent/received and the content of the message.

# Proposed Design

- We started by breaking up the program into 4 processes. 
- The server process, stored in pkg/server/ 
- The client process, stored in pkg/client/ 
- The config file and the program to read it, stored in pkg/config/ 
- And main.go

- We delegated separate responsibilities to each process. 
    - Server contains the goroutines and tcp threads that send and recieve messages between server instances. 
    - Client reads input from the terminal and forwards it to its corresponding server instance to send the message. 
    - Config reads data from the configuration file to assign the min and max delay, as well as ID's, IP's, and Ports to each server. 
    - And main, which acts as an entry point to the program and wraps all the separate processes in one file.


# Implementation / Flow of execution

## Main program in main.go

### Flow of execution

- main takes a command line argument to initialize a process of a given ID.
- To do this, the program stores the ID number as a variable and creates a map from the .config file using the FetchConfig function from the config pkg. 
- In addition to the map of the desired process information, main also obtains the min/max delay values from the .config file.
- Next we launch the server as a go routine by calling the Server function from it's subsequent pkg with the proccess address, config map, and delay values as parameters. 
- To give time for the server to boot-up, the program sleeps for brief moment before launching the MainClient from it's subsequent pkg.
- Now the simple_go_network is operational and the client and server can begin communicating with other running processes. 

## Config functionality in config.go

FetchConfig takes a map of strings with string key values, and an array of 2 integers as inputs. The function opens the .config files from it's subsequent package and reads each line. The delay values on line 1 are stored in the array and the rest of the lines are stored in the map, mapping each process ID to their relative address:port before safely closing the file.

    func FetchConfig() (map[string]string, [2]int)

## Client functionality in client.go

listening takes a connection and reads messages from it. When a message is recieved it is printed for the user.

    func listening(c net.Conn)

MainClient takes an address and dials a connection to it. Once established user input is read from the command-line and sent to the network layer.

    func MainClient(address string)
    
### Flow of execution

Begins in MainClient:

- First the client connects to the server from the given address and runs listening as a go routine that is reading any incoming/outgoing messages from Server, noted below.
- Next the user is prompted for input and a for loop continually reads the user input and sends it to the network layer to be used in Server.
- Finally the program continues until a STOP command is supplied to the client.

## Server functionality in server.go

### Network layer 'main' function

The Server function in server.go serves as a main function to initialize needed variables to properly implement unicastSend and unicastRecieve, aswell as the incoming/outgoing routines to communicate with other processes. Server then reads from the application layer to store the destination address,ID, and desired message to be sent later.

    func Server(address string, addrMap map[string]string, delay [2]int)
    
parseInput parses a string and outputs the destination ID and message from the user input fed in by the client.

    func parse_input(raw_input string) (string, string)

### Incoming messages
incomingRoutine waits to accept a connection from the host client and will pass both the sending and recieving client connections to unicastRecieve.
    
    func incoming_routine(l net.Listener, client net.Conn)
    
unicastRecieve reads user input from the sending client, stores the message in a clean format, and outputs the time recieved as well as the 'Recieved message' text before closing the connection to the sending client. 
    
    func unicast_recieve(sending net.Conn, recieving net.Conn)

### Outgoing messages

outgoingRoutine takes the delay min/max from the config file, an channel storing the most recent outgoing message, and a connection from the sending client. This fuction stores the message data from the recieving client via channel in Server, and then proceeds to sleep the routine for a random duration bounded by the delay values and calls unicastSend within the routine.  Finally the function displays the current time with a "Sent message" text.

    func outgoing_routine(delays [2]int, outgoing_messages chan Message, client net.Conn)
    
unicastSend takes the destination addresss and message, dials the destination server. The function then sends the message to the recieving client and subsequently closes the connection.

    func unicast_send(destination string, message string)

### Flow of execution

The flow of execution begins with Server:

- Initially the port value from that was parsed in main.go is stored and formatted properly.
- The server then begins listening on the specified port and opens a connection to possible clients created by user command line inputs. 
- Next the variables needed to run the communication go routines are initialized/passed to the incoming/outgoing functions and the routines are launched. These routines are responsible for handling incoming and outgoing messages, aswell as calling the unicastRecieve/Send functions.

At this stage incoming/outgoing routines are running as threads, waiting to continue execution using values from the host/remote servers:

- Server then opens communication with the main client for the process and reads the user input.
- The parseInput function then parses the input and extracts the destination ID and message from the user to be stored as a message struct.
- The input is then stored in a channel that is then passed to the outgoing routine to send a message as noted above.
- On the contrary, the incoming routine is still running and awaiting a new message in order to recieve it.
- The server then remains operational until a STOP command is read from the main client.
