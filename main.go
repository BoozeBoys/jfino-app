package main

import (
	"fmt"
	"jfino-app/slackbot"
	"os"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")
	sb := slackbot.New(token)

	for {
		recv, err := sb.Recv()
		if err != nil {
			panic(err)
		}

		fmt.Println(recv.Msg)
		sb.Reply(recv, recv.Msg)
	}
}
