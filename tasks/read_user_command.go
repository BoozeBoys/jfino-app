package tasks

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/BoozeBoys/jfino-app/loc"

	"github.com/BoozeBoys/jfino-app/slackbot"
	"github.com/BoozeBoys/jfino-app/state"
)

type ReadUserCommand struct {
	bot *slackbot.Slackbot
}

func NewReadUserCommand(slackToken string) *ReadUserCommand {
	bot := slackbot.New(slackToken)
	go bot.Start()

	return &ReadUserCommand{bot: bot}
}

func (t *ReadUserCommand) Perform(s *state.State) error {
	msg, err := t.bot.GetLastMessage()
	if err != nil {
		return nil
	}

	tokens := strings.Split(msg.Msg, " ")
	if len(tokens) == 0 {
		t.bot.Reply(msg, "?")
		return nil
	}

	switch strings.ToLower(tokens[0]) {
	case "accenditi":
		s.Power = true
		t.bot.Reply(msg, "OK :smile:")
	case "spengiti":
		s.Power = false
		t.bot.Reply(msg, "OK :smile:")
	case "fermati":
		fallthrough
	case "stop":
		s.Motors[0].Speed = 0
		s.Motors[1].Speed = 0
		t.bot.Reply(msg, "OK :smile:")
	case "vai":
		if len(tokens[1:]) != 1 {
			t.bot.Reply(msg, "forse volevi dire _vai piano_ o _vai sodo_")
			return nil
		}
		switch tokens[1] {
		case "piano":
			s.Motors[0].Speed = 50
			s.Motors[1].Speed = 50
			t.bot.Reply(msg, "OK :smile:")
		case "sodo":
			s.Motors[0].Speed = 255
			s.Motors[1].Speed = 255
			t.bot.Reply(msg, "OK :smile:")
		default:
			t.bot.Reply(msg, "forse volevi dire _vai piano_ o _vai sodo_")
		}

	case "power":
		if len(tokens[1:]) != 1 {
			t.bot.Reply(msg, "forse volevi dire _power 1_ o _power 0_")
			return nil
		}
		s.Power = tokens[1] == "1"
		t.bot.Reply(msg, "OK :smile:")

	case "speed":
		if len(tokens[1:]) != 2 {
			t.bot.Reply(msg, "forse volevi dire _speed 255 255_")
			return nil
		}

		speed1, err := strconv.Atoi(tokens[1])
		if err != nil {
			t.bot.Reply(msg, "il valore del motore 1 deve essere tra -255 e 255")
			return nil
		}

		speed2, err := strconv.Atoi(tokens[2])
		if err != nil {
			t.bot.Reply(msg, "il valore del motore 2 deve essere tra -255 e 255")
			return nil
		}

		s.Motors[0].Speed = speed1
		s.Motors[1].Speed = speed2

		t.bot.Reply(msg, "OK :smile:")
	case "status":
		data, err := json.Marshal(s)
		if err != nil {
			t.bot.Reply(msg, "non ci riesco :scream:")

		}
		t.bot.Reply(msg, string(data))
	case "naviga":
		if len(tokens[1:]) != 2 {
			t.bot.Reply(msg, "forse volevi dire _naviga x y_")
			return nil
		}

		x, err := strconv.ParseFloat(tokens[1], 64)
		if err != nil {
			t.bot.Reply(msg, "forse volevi dire _naviga x y_")
			return err
		}

		y, err := strconv.ParseFloat(tokens[2], 64)
		if err != nil {
			t.bot.Reply(msg, "forse volevi dire _naviga x y_")
			return err
		}

		s.DestinationPoint[0] = loc.Meters(x)
		s.DestinationPoint[1] = loc.Meters(y)
		s.DestinationPoint[2] = 0
		s.Autopilot = true
		t.bot.Reply(msg, "OK, viaggio! :smile:")
	default:
		t.bot.Reply(msg, "non ho capito! :scream:")
	}

	return nil
}
