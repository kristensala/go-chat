# Go-chat
![image](/img/2023-11-10_00-41.png)

## How to use proto files
* Run in the root of the project to generate message protobufs

`protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    msgpb/message.proto`

* then Run
`go get -u all`
