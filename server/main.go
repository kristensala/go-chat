package main

import (
	//"google.golang.org/grpc"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/kristensala/go-chat/msgpb"
	"google.golang.org/grpc"
)

var (
    port = flag.Int("port", 50011, "Server port")
)

type server struct {
    pb.UnimplementedCommunicationServiceServer
}

func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Message, error) {
    log.Printf("Message: %v", in.GetBody())
    return &pb.Message{Body: "Hello from server: " + in.GetBody()}, nil
}


func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCommunicationServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
