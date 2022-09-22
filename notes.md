# Brainstorming

## Prompt Notes:

Proc A: TCP Server --> Proc B: TCP Client

Process A (server) needs to be listening
Using process B (client), dial a connection to the sever on process A
-creates `c` on process B
Process A (server) accepts the dial connection from Process B
-creates `c` on process A
The channel `c` can be used as a line of communication between server and client

## Design notes:

config.txt :: First line is min-max delay for processes, each subsequent line is an individual process with format "ID# IP# Port#".
main.go :: Go file that will contain the methods unicast_send and unicast_recieve, aswell as handling of the various process layers.

### main.go structure?

unicast_send function(destination, message)
unicast_recieve function(source, message)

main function:

- Open config.txt file, parse the file to get arguments.

  - each iteration we store the arguments needed then pass to respective processes?

  - Application layer go routines ::
    - Handle user input from command line, will pass these arguments to Network layer via channels?
  - Network layer go routines ::
    - Each process will have both a recive and send channel, so each process is a client and server at the same time.
    - Handle the delay by generating a random number from min-max delay time and wait for that generated time.

### Server/Client Configuration

- Each process has a server and a client
- The servers of each provess implement `unicast_send` and `unicast_receive` functions
  - When a server is started it dials a connection to all of the currently running processes
  - `unicast_send`
    - Send the message that is input and delievers it to the destination server along with the source of the message
  - `unicast_receive`
    - Receives an incoming message and delivers it to the client to display along with the source address of the message
- The client of each process connects directly to its respective server and uses its server to send and receive messages
