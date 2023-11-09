package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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
    pulledMessagesFromServer []*pb.Message
    lastPulledMessageTime int64
    username string
}

func initApp() chatApp {
    usernameArg := os.Args[1]
    log.Println(usernameArg)

    client, err := connectToServer()
    if err != nil {
        log.Fatalf("Could not connect to the server %v", err)
    }

    return chatApp{
        serviceClient: client,
        username: *&usernameArg,
    }
}

//https://github.com/grpc/grpc-go/tree/master/examples/features/keepalive
//https://programmingpercy.tech/blog/streaming-data-with-grpc/
func main() {
    chatApp := app.New()
    window := chatApp.NewWindow("Go-chat")
    window.Resize(fyne.NewSize(500, 600))
    window.SetFixedSize(true)

    input := widget.NewMultiLineEntry()
    input.SetMinRowsVisible(1)
    input.SetPlaceHolder("Message")

    content := container.NewScroll(historyContainer)
    middlePart := container.NewBorder(nil, nil, nil ,nil, content)

    application := initApp()
    go application.getMessages(content)

    controls := container.New(
        layout.NewGridLayout(1),
        input,
        widget.NewButton("Send", func() {
            if len(input.Text) == 0 {
                return
            }

            application.sendMessage(application.username, input.Text)

            input.SetText("")
        }))


    window.SetContent(container.NewBorder(nil, controls, nil, nil, middlePart))
    window.ShowAndRun()
}

func (c *chatApp) buildUiMessage(senderName string, body string) *widget.RichText {
    textColor := theme.ColorNameError
    if c.isMessageSenderMe(senderName) {
        textColor = theme.ColorNamePrimary
    }
    
    style := widget.RichTextStyle{
		ColorName: textColor,
		Inline:    true,
		SizeName:  theme.SizeNameText,
		TextStyle: fyne.TextStyle{Bold: true},
    }

    txt := widget.NewRichText(
        &widget.TextSegment{
            Style: style,
            Text: senderName + " ",
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

func (c *chatApp) getMessages(contentContainer *container.Scroll) {
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
                c.pulledMessagesFromServer = append(c.pulledMessagesFromServer, message)
                if c.lastPulledMessageTime < message.GetUnixDateTime() {
                    c.lastPulledMessageTime = message.GetUnixDateTime()
                }

                uiMessage := c.buildUiMessage(message.GetUser().GetName(), message.GetBody())

                var wg sync.WaitGroup
                wg.Add(1)
                go func() {
                    historyContainer.Add(uiMessage)
                    wg.Done()
                }()
                wg.Wait()

                contentContainer.ScrollToBottom()
            }
        }
    }()
    select{}
}

func (c *chatApp) sendMessage(senderName string, body string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

    sendDateTimeUnix := time.Now().Unix()
    _, err := c.serviceClient.SendMessage(ctx, &pb.Message{
        Body: body,
        UnixDateTime: sendDateTimeUnix,
        User: &pb.User{
            Name: senderName,
        },
    })
	if err != nil {
        log.Printf("Could not send the message: %v", err)
        return
	}

	//log.Printf("message body: %s", r.GetBody())
}

func (c *chatApp) isMessageSenderMe(senderName string) bool {
    return senderName == c.username
}

