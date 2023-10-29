package main

import (
	//"google.golang.org/grpc"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/kristensala/go-chat/msgpb"
	"google.golang.org/grpc"
)

var (
    port = flag.Int("port", 50011, "Server port")
)

var messages []*pb.Message;

type server struct {
    pb.UnimplementedCommunicationServiceServer
}

func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Message, error) {
    messages = append(messages, in)
    return &pb.Message{Body: in.GetBody()}, nil
}

func (s *server) GetMessages(req *pb.NoParams, stream pb.CommunicationService_GetMessagesServer) error {
    timer := time.NewTicker(2 * time.Second)

    for {
        select {
        case <-stream.Context().Done():
            return nil
        case <-timer.C:
            for _, ms := range messages {
                err := stream.Send(ms)
                if err != nil {
                    log.Println(err.Error())
                }
            }
        }
    }
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
