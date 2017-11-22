package testutils

import (
	"bytes"
)

type SerialMock struct {
	reply    *bytes.Buffer
	commands []byte
}

func NewSerialMock(reply [][]byte) *SerialMock {
	buf := new(bytes.Buffer)

	for _, line := range reply {
		buf.Write(append(line, []byte("\r\n")...))
	}

	return &SerialMock{reply: buf}
}

func (sm *SerialMock) Read(data []byte) (int, error) {
	return sm.reply.Read(data)
}

func (sm *SerialMock) Write(data []byte) (int, error) {
	sm.commands = append(sm.commands, data...)

	return len(data), nil
}

func (sm *SerialMock) Commands() [][]byte {
	commands := bytes.Split(sm.commands, []byte("\r\n"))
	return commands[:len(commands)-1]
}

func (sm *SerialMock) LastCommand() []byte {
	commands := sm.Commands()
	return commands[len(commands)-1]
}
