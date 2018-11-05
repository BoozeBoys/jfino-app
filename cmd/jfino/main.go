package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BoozeBoys/jfino-app/commander"
	"github.com/BoozeBoys/jfino-app/state"
	"github.com/BoozeBoys/jfino-app/tasks"
	"github.com/tarm/serial"
)

type Config struct {
	SerialDevice string
	BaudRate     int
	Anchors      map[string]tasks.AnchorCfg
	SlackToken   string
	HttpPort     int
	Root         string
}

func (c *Config) String() string {
	return fmt.Sprintf("SerialDevice: %s\nBaudRate: %d\nAnchors: %v\nhttp port: %d\nroot: %s", c.SerialDevice, c.BaudRate, c.Anchors, c.HttpPort, c.Root)
}

func loadConfig(path string) (*Config, error) {
	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [CONFIG FILE]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	config, err := loadConfig(flag.Arg(0))
	if err != nil {
		log.Fatalf("error loading config: %s", err)
	}

	fmt.Println(config)

	s, err := serial.OpenPort(&serial.Config{Name: config.SerialDevice, Baud: config.BaudRate})
	if err != nil {
		log.Fatal("open port: ", err)
	}
	defer s.Close()

	time.Sleep(2 * time.Second)

	c := commander.New(s)

	taskList := []tasks.Task{
		tasks.NewUpdateStatus(c),
		tasks.NewEstimatePosition(config.Anchors),
		tasks.NewReadUserCommand(config.SlackToken),
		tasks.NewSendCommands(c),
		tasks.NewHttpMap(config.HttpPort, config.Root),
	}

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	state := new(state.State)

	for range t.C {
		for _, t := range taskList {
			if err := t.Perform(state); err != nil {
				log.Println(err)
			}
		}
	}
}
