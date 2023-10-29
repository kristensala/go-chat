package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/kristensala/go-chat/msgpb"
)

var (
	addr = flag.String("addr", "localhost:50011", "the address to connect to")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCommunicationServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SendMessage(ctx, &pb.Message{Body: "this is a message"})
	if err != nil {
		log.Fatalf("could not send a message: %v", err)
	}

	log.Printf("message body: %s", r.GetBody())
}
