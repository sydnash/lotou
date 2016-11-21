# lotou

a golang game server framework

I want write it as a game server.


## encoding

encoding is used to encode and decode socket message, here provide to encode type:
1. a binary encode like go's json, but it encode a struct to binary, and decode from binary stream.
2. a gob encode, it encode each value's type, can auto parse to interface{} of that type.

### binary

#### useage:
	see binary/binary_test.go

### gob

#### useage:
	see gob/type_test.go
	
