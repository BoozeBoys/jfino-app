package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	token := os.Getenv("SLACK_TOKEN")
	logger.Println("SLACK_TOKEN:", token)
	api := slack.New(token)

	slack.SetLogger(logger)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				logger.Println("Connected, counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				logger.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				botname := fmt.Sprintf("<@%s> ", info.User.ID)

				_, err1 := api.GetChannelInfo(ev.Channel)
				_, err2 := api.GetGroupInfo(ev.Channel)
				isDirectMsg := err1 != nil && err2 != nil

				if ev.User != info.User.ID &&
					(strings.HasPrefix(ev.Text, botname) || isDirectMsg) {
					rtm.SendMessage(rtm.NewOutgoingMessage("TEST!", ev.Channel))
				}

			case *slack.RTMError:
				logger.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				logger.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}
