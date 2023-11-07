package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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
    historyContainer = container.New(layout.NewVBoxLayout())
)

type chatApp struct {
    serviceClient pb.CommunicationServiceClient
    pulledMessagesFromServer []pb.Message
    lastPulledMessageTime int64
}

func initApp() chatApp {
    client, err := connectToServer()
    if err != nil {
        log.Fatalf("Could not connect to the server %v", err)
    }

    return chatApp{
        serviceClient: client,
    }
}

//https://github.com/grpc/grpc-go/tree/master/examples/features/keepalive
//https://programmingpercy.tech/blog/streaming-data-with-grpc/
func main() {
    //TODO pass cli argument (username) when starting the client 
    log.Println("starting")
    application := initApp()
    //defer application.serviceClient.Close()

    go application.getMessages()

    chatApp := app.New()
    window := chatApp.NewWindow("Go-chat")
    window.Resize(fyne.NewSize(500, 600))
    window.SetFixedSize(true)

    input := widget.NewMultiLineEntry()
    input.SetMinRowsVisible(1)
    input.SetPlaceHolder("Message")

    //history := container.New(layout.NewVBoxLayout())
    content := container.NewScroll(historyContainer)
    middlePart := container.NewBorder(nil, nil, nil ,nil, content)

    controls := container.New(
        layout.NewGridLayout(1),
        input,
        widget.NewButton("Send", func() {
            if len(input.Text) == 0 {
                return
            }

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

            //historyContainer.Add(txt) //TODO add the listened messages
            application.sendMessage("user", input.Text)

            content.ScrollToBottom()
            input.SetText("")
        }))


    window.SetContent(container.NewBorder(nil, controls, nil, nil, middlePart))
    window.ShowAndRun()
}

func buildUiMessage(senderName string, body string) *widget.RichText {
    txt := widget.NewRichText(
        &widget.TextSegment{
            Style: widget.RichTextStyleStrong,
            Text: senderName,
        },
        &widget.TextSegment{
            Style: widget.RichTextStyleParagraph,
            Text: body,
        },
    )

    txt.Wrapping = fyne.TextWrap(fyne.TextWrapWord)
    return txt
}

func connectToServer() (pb.CommunicationServiceClient, error) {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
        return nil, err
	}

    c := pb.NewCommunicationServiceClient(conn)

    return c, nil
}

func (c *chatApp) getMessages() {
    go func() {
        for {
            request := &pb.GetMessagesRequest{
                UnixDateTimeFrom: c.lastPulledMessageTime,
            }
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()

            response, err := c.serviceClient.GetMessages(ctx, request)
            if err != nil {
                log.Printf("Could not get messages %v", err)
            }

            messages := response.GetMessages()
            for _, message := range messages {
                fmt.Println(message.GetBody())
                if c.lastPulledMessageTime < message.GetUnixDateTime() {
                    c.lastPulledMessageTime = message.GetUnixDateTime()
                }
            }
        }
    }()
    select{} // runs in Background
}

func listenMessagesStream() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
        log.Fatalf("Could not connect to the server %v", err)
	}
	defer conn.Close()
	c := pb.NewCommunicationServiceClient(conn)

    //TODO instead using a stream lets just use an endpoint that returns a list
    stream, err := c.GetMessageStream(context.Background(), &pb.EmptyRequest{})
    if err != nil {
        log.Fatalf("could not get messages: %v", err)
    }

    go func() {
        for {
            res, err := stream.Recv()
            if err == io.EOF {
                log.Println("EOF error, closing stream")
                stream.CloseSend() //@todo: What if i don't close and keep looping?

                return
            }
            //historyContainer.Add(buildUiMessage(res.GetUser().GetName(), res.GetBody()))
            if res.GetBody() != "" {
                fmt.Printf("Message in server: %v", res)
            }

        }
    }()
    select{} //non blocking
}

func (c *chatApp) sendMessage(senderName string, body string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

    sendDateTimeUnix := time.Now().Unix()
    _, err := c.serviceClient.SendMessage(ctx, &pb.Message{
        Body: body,
        UnixDateTime: sendDateTimeUnix,
    })
	if err != nil {
        log.Printf("Could not send the message: %v", err)
        return
	}

	//log.Printf("message body: %s", r.GetBody())
}

// Todo
func getLastMessageTime() int64 {
    lastMessageTime := int64(0)

    /*for _, message := range msgsFromServer {
        //todo if lastMessageTime < message.time then lastMessageTime == message.Time
    }*/
    return lastMessageTime
}

