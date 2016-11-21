# lotou

a golang game server framework

I want write it as a game server.


## encoding

encoding is used to encode and decode socket message, here provide to encode type:  
1. a binary encode like go's json, but it encode a struct to binary, and decode from binary stream.  
2. a gob encode, it encode each value's type, can auto parse to interface{} of that type.  

### binary

 [useage](https://github.com/sydnash/lotou/blob/master/encoding/binary/binary_test.go)

### gob

 [useage](https://github.com/sydnash/lotou/blob/master/encoding/gob/type_test.go)
	
## log

log模块用于打印调试和错误信息，目前使用的是异步打印方式，暂不支持同步打印，这种模式会在程序panic的时候有些关键信息无法打印出来，还在考虑是否换位同步打印模式。

##network

network模块实现了一个基于tcp的server和client，暂时还未加入心跳机制，client会在需要发送数据的时候自动建立tcp连接。

## core

core提供service之间通信的桥梁，所有的service都会在core进行注册，然后service之间通过core.send进行消息发送。

## topology

topology，用于支持拓扑结果的服务器，一个master和N个slave，master会左右数据交互中心，不通logic主机上的消息会通过master进行中转。正在实现中...