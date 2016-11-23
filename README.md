# lotou

首先：本人初学go语言，希望用go语言来学习游戏服务器编程，欢迎指正，  
如果有兴趣，可以联系我  
QQ：157621271  
wechat：daijun_1234  
申请请写lotou  

a golang game server framework

I want write it as a game server.

lotou模仿了一部分skynet的模式，服务分为多个service，service之间通过消息进行通信，不会进行
直接的函数调用。目前通过core对消息进行转发，每一个service都有一个chan用于接收core转发的消
息。

## encoding

encoding is used to encode and decode socket message, here provide to encode type:  
1. a binary encode like go's json, but it encode a struct to binary, and decode from binary stream.  
2. a gob encode, it encode each value's type, can auto parse to interface{} of that type.  

### binary
binary用于进行二进制编码，主要用户客户端、服务器通信。
使用方式和golang的json有点类似。
 [useage](https://github.com/sydnash/lotou/blob/master/encoding/binary/binary_test.go)

### gob
gob用于实现服务节点间的通信，由于服务间的通信有可能在本机有可能跨节点，
本地消息通信不需要进行编码，跨节点通信的时候则需要进行编码。gob会同时将
数据的类型和值编码到package中，这样就不需要对消息进行注册，不过消息中的
结构体类型必须先注册到gob中进行，为了保证struct的ID一致，每一个节点都要
使用相同的顺序注册所有结构体。

目前对于slice、map的编码不够高效，为了可以编解码[]interface{}类型的slice，
目前会对slice每一个元素都对类型进行编码，希望可以根据elem类型进行适当的
调整。

内部使用reflect实现。
 [useage](https://github.com/sydnash/lotou/blob/master/encoding/gob/type_test.go)
	
## log

log模块用于打印调试和错误信息，目前使用的是异步打印方式，暂不支持同步打印，这种模式会在程序panic的时候有些关键信息无法打印出来，还在考虑是否换为同步打印模式。

##network

network模块实现了一个基于tcp的server和client，暂时还未加入心跳机制，client会在需要发送数据的时候自动建立tcp连接。

## core

core提供service之间通信的桥梁，所有的service都会在core进行注册，然后service之间通过core.send进行消息发送。

## topology

topology，用于支持拓扑结构的服务器，一个master和N个slave，master会作为数据交互中心，不同logic主机上的消息会通过  master进行中转。
目前可以通过slave master模块，使节点可以给本地节点的服务注册全局名字，可以通过全局名字获取服务ID，同时对发送给非
本地节点的消息进行转发。
