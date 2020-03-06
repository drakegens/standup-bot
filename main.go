package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
)

type ChannelMsg struct {
	Id    int
	Name  string
	Fields json.RawMessage
}

type MemberMsg struct {
	RadQuery        string `json:"RadQuery"`
	CorrelatedQuery string `json:"CorrelatedQuery"`
}

func main() {

	token := os.Getenv("SLACK_API_KEY")
	api := slack.New(token)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	channels, err := api.GetChannels(true)
	if err != nil {
		log.Fatalf("%s: %s", "Unable to get channels", err)
	}

	var memberIds []string
//	userId := "hello"
	for _, channel := range channels {
		fmt.Println(channel.Name)
		if channel.Name == "standup"{
			memberIds = channel.Members
			//userId = channel.ID
			break
		}
	}

	fmt.Println(memberIds)
	//memberId := memberIds[0]

	for _, memberId := range memberIds {
		fmt.Println(memberId, "hey")
		_, _, channelId, _ := api.OpenIMChannel(memberId)
		rtm.SendMessage(rtm.NewOutgoingMessage("What did you do yesterday?", channelId))
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
				fmt.Printf("Received a message!!!!!")
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
