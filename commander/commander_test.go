package commander_test

import (
	"bytes"
	"testing"

	"github.com/BoozeBoys/jfino-app/commander"
)

type SerialMock struct {
	reply   *bytes.Buffer
	command []byte
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
	sm.command = append(sm.command, data...)

	return len(data), nil
}

func (sm *SerialMock) Command() []byte {
	return sm.command
}

func TestPower(t *testing.T) {
	cases := []struct {
		reply [][]byte
		on    bool
		err   error
		want  []byte
	}{
		{[][]byte{[]byte("OK")}, true, nil, []byte("POWER 1")},
		{[][]byte{[]byte("OK")}, false, nil, []byte("POWER 0")},
		{[][]byte{[]byte("ERR")}, true, commander.Error, []byte("POWER 1")},
		{[][]byte{[]byte("ERR")}, false, commander.Error, []byte("POWER 0")},
	}

	for _, tc := range cases {
		rw := NewSerialMock(tc.reply)
		c := commander.New(rw)

		err := c.Power(tc.on)

		if err != tc.err {
			t.Errorf("want: %v got: %v\n", tc.err, err)
		}

		got := rw.Command()
		want := append(tc.want, []byte("\r\n")...)

		if !bytes.Equal(want, got) {
			t.Errorf("want: %s, got: %s\n", want, got)
		}
	}
}

func TestSpeed(t *testing.T) {
	cases := []struct {
		reply  [][]byte
		speed1 int
		speed2 int
		err    error
		want   []byte
	}{
		{[][]byte{[]byte("OK")}, 255, 255, nil, []byte("SPEED 255 255")},
		{[][]byte{[]byte("OK")}, -255, -255, nil, []byte("SPEED -255 -255")},
		{[][]byte{[]byte("OK")}, 0, 0, nil, []byte("SPEED 0 0")},
		{[][]byte{[]byte("ERR")}, 255, 255, commander.Error, []byte("SPEED 255 255")},
		{[][]byte{[]byte("ERR")}, -255, -255, commander.Error, []byte("SPEED -255 -255")},
		{[][]byte{[]byte("ERR")}, 0, 0, commander.Error, []byte("SPEED 0 0")},
	}

	for _, tc := range cases {
		rw := NewSerialMock(tc.reply)
		c := commander.New(rw)

		err := c.Speed(tc.speed1, tc.speed2)

		if err != tc.err {
			t.Errorf("want: %v got: %v\n", tc.err, err)
		}

		got := rw.Command()
		want := append(tc.want, []byte("\r\n")...)

		if !bytes.Equal(want, got) {
			t.Errorf("want: %s, got: %s\n", want, got)
		}
	}
}

func TestStatus(t *testing.T) {
	cases := []struct {
		reply [][]byte
		err   error
		want  []byte
	}{
		{[][]byte{
			[]byte("POWER 1"),
			[]byte("SPEED 0 255"),
			[]byte("SPEED 1 255"),
			[]byte("CURRENT 0 1023"),
			[]byte("CURRENT 1 1023"),
			[]byte("OK"),
		}, nil, []byte("STATUS")},
		{[][]byte{[]byte("ERR")}, commander.Error, []byte("STATUS")},
	}

	for _, tc := range cases {
		rw := NewSerialMock(tc.reply)
		c := commander.New(rw)

		reply, err := c.Status()

		if err != tc.err {
			t.Errorf("want: %v got: %v\n", tc.err, err)
		}

		if len(tc.reply) != len(reply)+1 {
			t.Errorf("want: %d got: %d\n", len(tc.reply), len(reply))
		}

		for n, line := range reply {
			if !bytes.Equal(tc.reply[n], line) {
				t.Errorf("want: %s, got: %s\n", tc.reply[n], line)
			}
		}

		got := rw.Command()
		want := append(tc.want, []byte("\r\n")...)

		if !bytes.Equal(want, got) {
			t.Errorf("want: %s, got: %s\n", want, got)
		}
	}
}
