package main

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"log"
)

type ChannelMsg struct {
	Id     int
	Name   string
	Fields json.RawMessage
}

type MemberMsg struct {
	RadQuery        string `json:"RadQuery"`
	CorrelatedQuery string `json:"CorrelatedQuery"`
}

var api *slack.Client
var rtm *slack.RTM

func getStandUpUpdateFromUser(memberId string) string {

	var update = "Yesterday: "

	_, _, channelId, _ := api.OpenIMChannel(memberId)
	//go func() {
	rtm.SendMessage(rtm.NewOutgoingMessage("What did you do yesterday?", channelId))

Loop1:
	select {
	case msg := <-rtm.IncomingEvents:
		fmt.Print("Event Received: ")
		for {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				update += ev.Text
				fmt.Printf("Message: %v\n", ev.Text)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop1

			default:
				//Take no action
			}
		}
	}
	rtm.SendMessage(rtm.NewOutgoingMessage("What are you going to do today?", channelId))
Loop2:
	select {
	case msg := <-rtm.IncomingEvents:
		fmt.Print("Event Received: ")
		for {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				update += "\nToday:" + ev.Text
				fmt.Printf("Message: %v\n", ev.Text)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop2

			default:
				//Take no action
			}
		}
	}

	//	}()
	return update

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
	var standupChannelId string

	for _, channel := range channels {
		fmt.Println(channel.Name)
		if channel.Name == "standup" {
			memberIds = channel.Members
			standupChannelId = channel.ID
			fmt.Println(standupChannelId)
			break
		}
	}

	var standupUpdates [10]string

	for i, memberId := range memberIds {
		standupUpdates[i] = getStandUpUpdateFromUser(memberId)
	}

	rtm.SendMessage(rtm.NewOutgoingMessage(standupUpdates[0], standupChannelId))

	//api.GetUsers()
	//
	//api.PostMessage()
	//Loop:
	//	for {
	//		select {
	//		case msg := <-rtm.IncomingEvents:
	//			fmt.Print("Event Received: ")
	//			switch ev := msg.Data.(type) {
	//			case *slack.ConnectedEvent:
	//				fmt.Println("Connection counter:", ev.ConnectionCount)
	//
	//			case *slack.MessageEvent:
	//				fmt.Printf("Message: %v\n", ev)
	//				fmt.Printf("Received a message!!!!!")
	//				info := rtm.GetInfo()
	//				prefix := fmt.Sprintf("<@%s> ", info.User.ID)
	//
	//				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
	//					rtm.SendMessage(rtm.NewOutgoingMessage("What's up buddy!?!?", ev.Channel))
	//				}
	//
	//			case *slack.RTMError:
	//				fmt.Printf("Error: %s\n", ev.Error())
	//
	//			case *slack.InvalidAuthEvent:
	//				fmt.Printf("Invalid credentials")
	//				break Loop
	//
	//			default:
	//				//Take no action
	//			}
	//		}
	//	}
}
