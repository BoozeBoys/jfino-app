package slackbot

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/nlopes/slack"
)

var ErrEmptyQueue = errors.New("no commands left")

type Slackbot struct {
	api *slack.Client
	rtm *slack.RTM

	m        sync.Mutex
	messages []BotMsg
}

type BotMsg struct {
	Channel   string
	User      string
	Msg       string
	DirectMsg bool
}

func New(token string) *Slackbot {
	sb := new(Slackbot)
	sb.api = slack.New(token)
	sb.rtm = sb.api.NewRTM()
	go sb.rtm.ManageConnection()
	return sb
}

func (sb *Slackbot) Start() {
	for {
		msg := <-sb.rtm.IncomingEvents

		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			info := sb.rtm.GetInfo()
			botname := fmt.Sprintf("<@%s> ", info.User.ID)

			_, err1 := sb.api.GetChannelInfo(ev.Channel)
			_, err2 := sb.api.GetGroupInfo(ev.Channel)
			isDirectMsg := err1 != nil && err2 != nil

			if ev.User != info.User.ID &&
				(strings.HasPrefix(ev.Text, botname) || isDirectMsg) {
				sb.m.Lock()
				sb.messages = append(sb.messages, BotMsg{ev.Channel, ev.User, strings.TrimPrefix(ev.Text, botname), isDirectMsg})
				sb.m.Unlock()
			}
		}
	}
}

func (sb *Slackbot) GetLastMessage() (BotMsg, error) {
	sb.m.Lock()
	defer sb.m.Unlock()

	var msg BotMsg

	if len(sb.messages) == 0 {
		return msg, ErrEmptyQueue
	}

	msg, sb.messages = sb.messages[0], sb.messages[1:]

	return msg, nil
}

func (sb *Slackbot) Reply(recv BotMsg, msg string) {
	var s string
	if recv.DirectMsg {
		s = msg
	} else {
		s = fmt.Sprintf("<@%s> %s", recv.User, msg)
	}
	sb.rtm.SendMessage(sb.rtm.NewOutgoingMessage(s, recv.Channel))
}
