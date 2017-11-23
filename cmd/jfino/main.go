package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BoozeBoys/jfino-app/commander"
	"github.com/BoozeBoys/jfino-app/state"
	"github.com/BoozeBoys/jfino-app/tasks"
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

	taskList := []tasks.Task{
		tasks.NewUpdateStatus(c),
		tasks.NewSendCommands(c),
	}

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	state := new(state.State)

	for range t.C {
		for _, t := range taskList {
			t.Perform(state)
		}
	}
}
