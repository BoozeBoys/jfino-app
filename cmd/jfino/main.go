package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BoozeBoys/jfino-app/commander"
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

	c := commander.New(s)

	if err := c.Power(true); err != nil {
		log.Fatal("setting power: ", err)
	}

	if err := c.Speed(255, 255); err != nil {
		log.Fatal("setting speed: ", err)
	}

	status, err := c.Status()
	if err != nil {
		log.Fatal("read status: ", err)
	}

	for _, line := range status {
		fmt.Println(string(line))
	}
}
