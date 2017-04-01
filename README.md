Lotou is a lightweight framework for game server implemented in golang.
It utilizes the communication between nodes by abstracting exchangings of messages 
so that the multi-nodes and multi-services can be easily achieved.

If you are interested in developing Lotou, feel free to contact me via 

   QQ 157621271 
   
   Wechat ID daijun_1234. 
   
   QQ group is also available 362112175. (Fill verification message with `Lotou`)
   
   twitter @sydnash1
   
Advices and geeks are always welcomed. 

[中文博客](http://blog.csdn.net/sydnash/article/details/66033983)
-------------------
# lotou
Lotou is inspired by the framework SKYNET written by cloudwu.
Within this framework, all services communicates each other by messaging. 
Functions supplied by services are not exposed by API calling. 
In current version, module core is responsible for routing messages, 
and every service has its own message chan for receiving and sending. 

# features

1. message based on communication between service(idea from share memory through communication in golang)
2. multicore support
3. cluster support 
4. extremely simple user interface

# main functional package

## core
core supplies the communication between services. 
All services are registered to `core`, and send message to others by core.send.

## binary

`binary` encodes data into binary stream for communication between servers and clients.
Usage of this encoding convention is somewhat like the usage of json marshalling. 

gob(not gob in golang's own lib)
Gob encoding aims to sustain nodes' communication. 
Communication of services can be trans-nodes, so the encoding messages into binary stream is necessary.
Message switching within nodes do not need to be encoded for better performance. 

All primitive types in golang can be encoded. 
Any combination of parameters is easily supported 
as long as the sender and the receiver complies with the same signature

self-defined struct can also be supported, but the struct must be registered in advanced.
If different nodes need to use the same struct, the registration of those structs must 
be registered in a fixed and unique sequence. 

At present, encoding performances for slice and map are not so good. 
In order to encode type of `[]interface{}`, each elements of which are encoded separately. 
To be more specific, type info for each elements are all marshalled. 
So if there are elements of one certain same type, type info are marshalled redundantly.

Marshalling of gob(Lotou) is achieved by reflect.

## log

Module of log is for printing debug and error information. 
Currently, log is Synchronized mode, if we have too many log which needs to wirte to file, it may block the main logic.
(it may be changed to asynchronized mode some day)

## network
network implements the routes between nodes in TCP. 
Roles of `server` and `client` work as what they are called. 
Heartbeat is not implemented yet. 
`client` establishes connection to its `server` at the first time of sending message. 

## topology
Multi-node servers are based on topoloy. 
There are two types of nodes in Lotou 
and different types of nodes are clusted into the same network.
Within a certain network there is only one master node and serveral slave nodes. 
Services can be run in any type of node and there are two types of services:
Local service and global service. 
Local service are ones within the same node known to each other by name.
Global service are ones run by different nodes
 and they are exposed to all services of all nodes by being registered to master.
 
 
 
 
