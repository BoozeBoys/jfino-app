package commander

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

var Error = errors.New("jfino error")
var Timeout = errors.New("operation timed out")

type Commander struct {
	rw *bufio.ReadWriter
}

func New(rw io.ReadWriter) *Commander {
	return &Commander{
		rw: bufio.NewReadWriter(
			bufio.NewReader(rw),
			bufio.NewWriter(rw),
		),
	}
}

func (c *Commander) WaitHelo() error {
	line, isPrefix, err := c.rw.ReadLine()
	if err != nil {
		return err
	}

	if isPrefix || !bytes.Equal(line, []byte("HELO")) {
		return Error
	}

	return nil
}

func (c *Commander) Power(on bool) error {
	var value int
	if on {
		value = 1
	}

	command := []byte(fmt.Sprintf("POWER %d", value))

	if _, err := c.sendCommand(command); err != nil {
		return err
	}

	return nil
}

func (c *Commander) Speed(speed1, speed2 int) error {
	command := []byte(fmt.Sprintf("SPEED %d %d", speed1, speed2))

	if _, err := c.sendCommand(command); err != nil {
		return err
	}

	return nil
}

func (c *Commander) Status() (payload [][]byte, err error) {
	return c.sendCommand([]byte("STATUS"))
}

func (c *Commander) sendCommand(command []byte) (payload [][]byte, err error) {
	command = append(command, []byte("\r\n")...)

	if _, err = c.rw.Write(command); err != nil {
		return
	}

	if err = c.rw.Flush(); err != nil {
		return
	}

	for {
		var reply []byte
		var isPrefix bool

		reply, isPrefix, err = c.rw.ReadLine()
		if bytes.Equal(reply, []byte("OK")) {
			return
		}

		if isPrefix || bytes.Equal(reply, []byte("ERR")) {
			err = Error
			return
		}

		r := make([]byte, len(reply))
		copy(r, reply)
		payload = append(payload, r)
	}
}
