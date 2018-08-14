package commander_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/BoozeBoys/jfino-app/commander"
	"github.com/BoozeBoys/jfino-app/testutils"
)

func TestHelo(t *testing.T) {
	cases := []struct {
		reply [][]byte
		err   error
	}{
		{[][]byte{[]byte("HELO")}, nil},
		{[][]byte{[]byte("FOO")}, commander.Error},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("reply=%v error=%v", tc.reply, tc.err != nil), func(tt *testing.T) {
			tt.Parallel()

			rw := testutils.NewSerialMock(tc.reply)
			c := commander.New(rw)

			err := c.WaitHelo()

			if err != tc.err {
				tt.Errorf("want: %v got: %v\n", tc.err, err)
			}
		})
	}
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
		{[][]byte{}, true, io.ErrShortWrite, []byte("POWER 1")},
		{[][]byte{}, false, io.ErrShortWrite, []byte("POWER 0")},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("on=%v error=%v", tc.on, tc.err != nil), func(tt *testing.T) {
			tt.Parallel()

			rw := testutils.NewSerialMock(tc.reply)
			c := commander.New(rw)

			err := c.Power(tc.on)

			if err != tc.err {
				tt.Errorf("want: %v got: %v\n", tc.err, err)
			}

			got := rw.LastCommand()
			want := tc.want

			if !bytes.Equal(want, got) {
				tt.Errorf("want: %s, got: %s\n", want, got)
			}
		})
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
		{[][]byte{}, 255, 255, io.ErrShortWrite, []byte("SPEED 255 255")},
		{[][]byte{}, 255, 255, io.ErrShortWrite, []byte("SPEED 255 255")},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("speed1=%v speed2=%v err=%v", tc.speed1, tc.speed2, tc.err != nil), func(tt *testing.T) {
			tt.Parallel()

			rw := testutils.NewSerialMock(tc.reply)
			c := commander.New(rw)

			err := c.Speed(tc.speed1, tc.speed2)

			if err != tc.err {
				tt.Errorf("want: %v got: %v\n", tc.err, err)
			}

			got := rw.LastCommand()
			want := tc.want

			if !bytes.Equal(want, got) {
				tt.Errorf("want: %s, got: %s\n", want, got)
			}
		})
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
		{[][]byte{}, io.ErrShortWrite, []byte("STATUS")},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("err=%v", tc.err != nil), func(tt *testing.T) {
			tt.Parallel()

			rw := testutils.NewSerialMock(tc.reply)
			c := commander.New(rw)

			reply, err := c.Status()

			if err != tc.err {
				tt.Errorf("want: %v got: %v\n", tc.err, err)
			}

			for n, line := range reply {
				if !bytes.Equal(tc.reply[n], line) {
					tt.Errorf("want: %s, got: %s\n", tc.reply[n], line)
				}
			}

			got := rw.LastCommand()
			want := tc.want

			if !bytes.Equal(want, got) {
				tt.Errorf("want: %s, got: %s\n", want, got)
			}
		})
	}
}
