package main

import (
	"context"
	"errors"
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
    messages []*pb.Message
    registeredUsers []*pb.User
}

func (s *server) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.User, error) {
    for _, user := range s.registeredUsers {
        if user.GetName() == request.GetUserName() {
            return nil, errors.New("Username already taken")
        }
    }

    newUser := &pb.User{
        Name: request.GetUserName(),
    }
    s.registeredUsers = append(s.registeredUsers, newUser)
    return newUser, nil
}

func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Message, error) {
    s.messages = append(s.messages, in)
    log.Printf("Messages received by server %v", in)
    return &pb.Message{Body: in.GetBody()}, nil
}

func (s *server) GetMessages(ctx context.Context, request *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
    filteredMessages := filterMessages(request.GetUnixDateTimeFrom(), s.messages)
    response := &pb.GetMessagesResponse{
        Messages: filteredMessages,
    }
    return response, nil
}

func (s *server) GetMessageStream(req *pb.EmptyRequest, stream pb.CommunicationService_GetMessageStreamServer) error {
    return nil
}

func filterMessages(unixDateTime int64, array []*pb.Message) []*pb.Message {
    var result []*pb.Message
    for _, message := range array {
        if unixDateTime < message.UnixDateTime {
            result = append(result, message)
        }
    }

    return result
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
