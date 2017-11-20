package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BoozeBoys/jfino-app/commander"
	"github.com/BoozeBoys/jfino-app/slackbot"
	"github.com/tarm/serial"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [FILE]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	s, err := serial.OpenPort(&serial.Config{Name: flag.Arg(0), Baud: 115200})
	if err != nil {
		log.Fatal("open port: ", err)
	}
	defer s.Close()

	time.Sleep(2 * time.Second)

	c := commander.New(s)
	bot := slackbot.New(os.Getenv("SLACK_TOKEN"))

	for {
		msg, err := bot.Recv()
		if err != nil {
			log.Fatal("recv msg: ", err)
		}
		args := strings.Split(msg.Msg, " ")
		switch args[0] {
		case "power":
			if err := c.Power(args[1] == "1"); err != nil {
				bot.Reply(msg, "Il comando ha riportato un errore, perdinci.")
				break
			}
			bot.Reply(msg, "Vabbene")
		case "speed":
			speed1, err := strconv.Atoi(args[1])
			if err != nil {
				bot.Reply(msg, "Speed1 non valida, coglione.")
				break
			}
			speed2, err := strconv.Atoi(args[2])
			if err != nil {
				bot.Reply(msg, "Speed2 non valida, deficiente.")
				break
			}
			if err := c.Speed(speed1, speed2); err != nil {
				bot.Reply(msg, "Il comando ha riportato un errore, perdinci.")
				break
			}
			bot.Reply(msg, "Sì padrone.")
		case "status":
			status, err := c.Status()
			if err != nil {
				bot.Reply(msg, "Il comando ha riportato un errore, perdinci.")
				break
			}
			res := bytes.Join(status, []byte{'\n'})
			bot.Reply(msg, string(res))
		default:
			bot.Reply(msg, "Non ho capito un cazzo, dé")
		}

	}
}
