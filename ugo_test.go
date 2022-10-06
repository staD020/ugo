package ugo

import (
	"bytes"
	"testing"
)

func TestCommandBytes(t *testing.T) {
	cases := []struct {
		cmd    Command
		length int
		want   []byte
	}{
		{Reset, 0, []byte{0x4, 0xff, 0x0, 0x0}},
		{DMARun, 0, []byte{0x2, 0xff, 0x0, 0x0}},
		{RunImage, 0, []byte{0xb, 0xff, 0x0, 0x0, 0x0}},
		{MountImage, 0, []byte{0xa, 0xff, 0x0, 0x0, 0x0}},

		{DMARun, 256, []byte{0x2, 0xff, 0x0, 0x1}},
		{DMARun, 51308, []byte{0x2, 0xff, 0x6c, 0xc8}},
		{RunImage, 116844, []byte{0xb, 0xff, 0x6c, 0xc8, 0x1}},
		{MountImage, 174848, []byte{0xa, 0xff, 0x0, 0xab, 0x2}},
	}
	for _, c := range cases {
		got := c.cmd.Bytes(c.length)
		if !bytes.Equal(got, c.want) {
			t.Errorf("Command %s len %d -> got: %v, want %v", c.cmd, c.length, got, c.want)
		}
	}
}
