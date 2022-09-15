Proc A: TCP Server --> Proc B: TCP Client

Process A (server) needs to be listening
Using process B (client), dial a connection to the sever on process A
-creates `c` on process B
Process A (server) accepts the dial connection from Process B
-creates `c` on process A
The channel `c` can be used as a line of communication between server and client
