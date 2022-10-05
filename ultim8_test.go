package ultim8

import (
	"testing"
)

func TestCommandBytes(t *testing.T) {
	cases := []struct {
		cmd    Command
		length int
		want   []byte
	}{
		{Reset, 0, []byte{0x4, 0xff, 0x0, 0x0}},
	}
	for _, c := range cases {
		got := c.cmd.Bytes(0)
		if !equalSlice(got, c.want) {
			t.Errorf("Command %s -> got: %v, want %v", c.cmd, got, c.want)
		}
	}
}

// equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equalSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
