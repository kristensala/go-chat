# Go-chat

## How to use proto files
* Run in the root of the project to generate message protobufs

`protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    msgpb/message.proto`
