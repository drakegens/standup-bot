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

//token := os.Getenv("SLACK_API_KEY")
var api = slack.New("xoxb-964613251380-964617721092-XYqK7HIDWv9leZzinW0IXhlS")
var rtm = api.NewRTM()

func getStandUpUpdateFromUser(memberId string) string {

	var update = "Yesterday: "

	_, _, channelId, err := api.OpenIMChannel(memberId)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	rtm.SendMessage(rtm.NewOutgoingMessage("What did you do yesterday?", channelId))

Loop1:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				update += ev.Text
				fmt.Printf("Message: %v\n", ev.Text)
				break Loop1

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())
				break Loop1

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
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				update += "\nToday: " + ev.Text
				fmt.Printf("Message: %v\n", ev.Text)
				break Loop2

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())
				break Loop2

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop2

			default:
				//Take no action
			}
		}
	}

	return update

}

func main() {
	go rtm.ManageConnection()

	channels, err := api.GetChannels(true)

	if err != nil {
		log.Fatalf("%s: %s", "Unable to get channels", err)
	}

	var memberIds []string
	var standupChannel slack.Channel

	for _, channel := range channels {
		fmt.Println(channel.Name)
		if channel.Name == "standup" {
			memberIds = channel.Members
			standupChannel = channel
			fmt.Println(standupChannel)
			break
		}
	}

	var standupUpdates [10]string

	for i, memberId := range memberIds {
		standupUpdates[i] = getStandUpUpdateFromUser(memberId)
	}

	_, _, err = api.PostMessage(standupChannel.ID, slack.MsgOptionText(standupUpdates[0], false))

	if err != nil {
		log.Fatalf("%s: %s", "Unable to post message", err)
	}
}
