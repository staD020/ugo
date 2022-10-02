// Package ultim8 provides 1541Ultimate control to start programs and disks via TCP.
// It is a port of Ucodenet by TTL (https://csdb.dk/release/?id=189723).
package ultim8

import (
	"fmt"
	"io"
	"net"
	"time"
)

//------------------------------------------------------------------------------
// generic structure is:
// command lo, command hi, number of parameters lo, number of parameters hi
// followed by the parameters

var DialTimeout = 3 * time.Second

type Command uint16

const (
	DMA         Command = 0xFF01 // dma-load .prg file
	DMARun      Command = 0xFF02 // dma-load .prg file and run it
	Keyboard    Command = 0xFF03 // simulate keyboard input
	Reset       Command = 0xFF04 // reset the c64
	Wait        Command = 0xFF05 // wait n ticks
	DMAWrite    Command = 0xFF06 // write c64 memory
	REUWrite    Command = 0xFF07
	KernalWrite Command = 0xFF08
	DMAJump     Command = 0xFF09 // dma-load .prg file and jump to addr
	MountImage  Command = 0xFF0A // mount image
	RunImage    Command = 0xFF0B // mount and run image
)

func (c Command) Bytes(length int) []byte {
	return []byte{
		byte(c & 0xff), byte(c >> 8),
		byte(length & 0xff), byte(length >> 8),
	}
}

func (c Command) String() string {
	return fmt.Sprintf("0x%04x", uint16(c))
}

type manager struct {
	addr string
	c    net.Conn
}

func New(address string) (*manager, error) {
	conn, err := net.DialTimeout("tcp", address, DialTimeout)
	if err != nil {
		return nil, fmt.Errorf("net.DialTimeout %q failed: %w", address, err)
	}
	return &manager{addr: address, c: conn}, nil
}

func (m *manager) sendCommand(cmd Command, length int) error {
	if _, err := m.c.Write(cmd.Bytes(length)); err != nil {
		return fmt.Errorf("m.c.Write failed: %w", err)
	}
	return nil
}

func (m *manager) Reset() error {
	if _, err := m.c.Write(Reset.Bytes(0)); err != nil {
		return fmt.Errorf("m.c.Write failed: %w", err)
	}
	time.Sleep(time.Second)
	return nil
}

func (m *manager) RunPrg(r io.Reader) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed: %w", err)
	}
	if err = m.Reset(); err != nil {
		return fmt.Errorf("m.Reset failed: %w", err)
	}
	if err = m.sendCommand(DMARun, len(buf)); err != nil {
		return fmt.Errorf("m.sendCommand DMARun failed: %w", err)
	}
	if _, err = m.c.Write(buf); err != nil {
		return fmt.Errorf("m.c.Write failed: %w", err)
	}
	return nil
}

func (m *manager) Close() error {
	return m.c.Close()
}
