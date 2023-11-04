package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	pb "github.com/kristensala/go-chat/msgpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50011", "the address to connect to")
)

//todo: ask for a username and then connect to the server
func main() {
    chatApp := app.New()
    window := chatApp.NewWindow("Go-chat")
    window.Resize(fyne.NewSize(500, 600))
    window.SetFixedSize(true)

    input := widget.NewMultiLineEntry()
    input.SetMinRowsVisible(1)
    input.SetPlaceHolder("Message")

    history := container.New(layout.NewVBoxLayout())
    content := container.NewScroll(history)
    middlePart := container.NewBorder(nil, nil, nil ,nil, content)

    controls := container.New(
        layout.NewGridLayout(1),
        input,
        widget.NewButton("Send", func() {
            if len(input.Text) == 0 {
                return
            }

            log.Println(input.Text)
            txt := widget.NewRichText(
                &widget.TextSegment{
                    Style: widget.RichTextStyleStrong,
                    Text: "sender name: ",
                },
                &widget.TextSegment{
                    Style: widget.RichTextStyleParagraph,
                    Text: input.Text,
                },
            )
            txt.Wrapping = fyne.TextWrap(fyne.TextWrapWord)

            history.Add(txt)
            content.ScrollToBottom()

            input.SetText("")
        }))


    window.SetContent(container.NewBorder(nil, controls, nil, nil, middlePart))
    window.ShowAndRun()
}

func connects() {
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

    emptyReq := &pb.NoParams{}
    stream, err := c.GetMessages(context.Background(), emptyReq)
	if err != nil {
		log.Fatalf("could not get messages: %v", err)
	}

    go func() {
        for {
            res, err := stream.Recv()
            if err != nil {
                log.Fatalf("could not receive from stream: %v", err)
            }
            fmt.Printf("Message history: %v", res)
            fmt.Println("----------")
        }
    }()

	r, err := c.SendMessage(ctx, &pb.Message{Body: "this is a message"})
	if err != nil {
		log.Fatalf("could not send a message: %v", err)
	}

	log.Printf("message body: %s", r.GetBody())

    select{}
}


//https://github.com/grpc/grpc-go/tree/master/examples/features/keepalive
//https://programmingpercy.tech/blog/streaming-data-with-grpc/
