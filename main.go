package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

type ChannelMsg struct {
	Id    int
	Name  string
	Fields json.RawMessage
}

func main() {

	token := os.Getenv("SLACK_API_KEY")
	fmt.Println(token)
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	channels, err := api.GetChannels(true)
	if err != nil {
		log.Fatalf("%s: %s", "Unable to get channels", err)
	}

	var memberIds []string
	channelId := "hello"
	for _, channel := range channels {
		fmt.Println(channel.Name)
		if channel.Name == "standup"{
			memberIds = channel.Members
			channelId = channel.ID
			break
		}
	}

	fmt.Println(memberIds)

	for _, memberId := range memberIds {
		fmt.Println(memberId, "hey")
		rtm.SendMessage(rtm.NewOutgoingMessage("Hey", channelId))
	}

	//api.GetUsers()
	//
	//api.PostMessage()
Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					rtm.SendMessage(rtm.NewOutgoingMessage("What's up buddy!?!?", ev.Channel))
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}
