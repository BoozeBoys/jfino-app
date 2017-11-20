package slackbot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

type Slackbot struct {
	api *slack.Client
	rtm *slack.RTM
}

type BotMsg struct {
	Channel   string
	User      string
	Msg       string
	DirectMsg bool
}

func New(token string) *Slackbot {
	sb := &Slackbot{}
	sb.api = slack.New(token)
	sb.rtm = sb.api.NewRTM()
	go sb.rtm.ManageConnection()
	return sb
}

func (sb *Slackbot) Recv() (*BotMsg, error) {

	for {
		select {
		case msg := <-sb.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				info := sb.rtm.GetInfo()
				botname := fmt.Sprintf("<@%s> ", info.User.ID)

				_, err1 := sb.api.GetChannelInfo(ev.Channel)
				_, err2 := sb.api.GetGroupInfo(ev.Channel)
				isDirectMsg := err1 != nil && err2 != nil

				if ev.User != info.User.ID &&
					(strings.HasPrefix(ev.Text, botname) || isDirectMsg) {
					return &BotMsg{ev.Channel, ev.User, strings.TrimPrefix(ev.Text, botname), isDirectMsg}, nil
				}

			case *slack.RTMError:
				return nil, errors.New("RTM error")

			case *slack.InvalidAuthEvent:
				return nil, errors.New("Invalid authentication")

			default:
				//Take no action
			}
		}
	}
}

func (sb *Slackbot) Reply(recv *BotMsg, msg string) {
	var s string
	if recv.DirectMsg {
		s = recv.Msg
	} else {
		s = fmt.Sprintf("<@%s> %s", recv.User, recv.Msg)
	}
	sb.rtm.SendMessage(sb.rtm.NewOutgoingMessage(s, recv.Channel))
}
