// Package ultim8 provides 1541Ultimate control to start programs and disks via TCP.
//
// It is a partial port of Ucodenet by TTL (https://csdb.dk/release/?id=189723) in pure Go by burg.
package ultim8

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// DialTimeout contains timeout for the initial TCP connection to your 1541u.
var DialTimeout = 3 * time.Second

const (
	D64Size = 174848
	Version = "0.1"
)

// Command specifies the various commands you can send to the 1541u.
type Command uint16

// All 1541u commands.
//
// Generic structure is:
// command lo, command hi, payload length lo, payload length hi
// followed by its payload, if any.
const (
	DMA         Command = 0xff01 // dma-load .prg file
	DMARun      Command = 0xff02 // dma-load .prg file and run it
	Keyboard    Command = 0xff03 // simulate keyboard input
	Reset       Command = 0xff04 // reset the c64
	Wait        Command = 0xff05 // wait n ticks
	DMAWrite    Command = 0xff06 // write c64 memory
	REUWrite    Command = 0xff07
	KernalWrite Command = 0xff08
	DMAJump     Command = 0xff09 // dma-load .prg file and jump to addr
	MountImage  Command = 0xff0a // mount image
	RunImage    Command = 0xff0b // mount and run image
)

// Bytes returns the bytes representing this Command, 2 bytes for the command and 2 or 3 bytes for length.
func (c Command) Bytes(length int) []byte {
	buf := []byte{byte(c & 0xff), byte(c >> 8)}
	if c == MountImage || c == RunImage {
		return append(buf, byte(length&0xff), byte((length>>8)&0xff), byte((length>>16)&0xff))
	}
	return append(buf, byte(length&0xff), byte(length>>8))
}

// String returns the string representation of the command.
func (c Command) String() string {
	s := "n/a"
	switch c {
	case DMA:
		s = "DMA"
	case DMARun:
		s = "DMARun"
	case Reset:
		s = "Reset"
	case RunImage:
		s = "RunImage"
	case MountImage:
		s = "MountImage"
	}
	return fmt.Sprintf("%-10s 0x%04x", s, uint16(c))
}

// Manager is the struct containing the net.Conn to your 1541u.
type Manager struct {
	addr string
	c    net.Conn
	done chan bool
}

// New establishes a new TCP connection your 1541u and returns the connection Manager.
// It also implements the io.Closer interface, callers are expected to Close() after use.
func New(address string) (*Manager, error) {
	conn, err := net.DialTimeout("tcp", address, DialTimeout)
	if err != nil {
		return nil, fmt.Errorf("net.DialTimeout %q failed: %w", address, err)
	}
	fmt.Println("[1541U] Connection established")
	m := &Manager{addr: address, c: conn, done: make(chan bool, 1)}
	go m.backgroundReader()
	return m, nil
}

// Send sends a bytestream of the Command, it's length and payload, which may be nil.
func (m *Manager) Send(cmd Command, payload []byte) error {
	if _, err := m.c.Write(cmd.Bytes(len(payload))); err != nil {
		return fmt.Errorf("Write failed: %w", err)
	}
	fmt.Printf("[CMD] %s\n", cmd)
	if len(payload) == 0 {
		return nil
	}
	if _, err := m.c.Write(payload); err != nil {
		return fmt.Errorf("Write failed: %w", err)
	}
	return nil
}

// Reset sends the Reset Command to the 1541u and sleeps for a second.
func (m *Manager) Reset() error {
	if err := m.Send(Reset, nil); err != nil {
		return fmt.Errorf("Send Reset failed: %w", err)
	}
	time.Sleep(time.Second)
	return nil
}

// Run drains the input Reader and uploads it to the 1541u with Command cmd.
// Before upload, the Reset Command is sent.
func (m *Manager) Run(r io.Reader) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed: %w", err)
	}
	if err = m.Reset(); err != nil {
		return fmt.Errorf("Reset failed: %w", err)
	}
	cmd := DMARun
	if len(buf) >= D64Size {
		cmd = RunImage
	}
	if err = m.Send(cmd, buf); err != nil {
		return fmt.Errorf("Send %s failed: %w", cmd, err)
	}
	return nil
}

// Mount drains the Reader and uploads it to the 1541u with Command MountImage or DMA.
func (m *Manager) Mount(r io.Reader) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed: %w", err)
	}
	cmd := DMA
	if len(buf) >= D64Size {
		cmd = MountImage
	}
	if err = m.Send(cmd, buf); err != nil {
		return fmt.Errorf("Send %s failed: %w", cmd, err)
	}
	return nil
}

// Close closes the TCP connection and waits for clean disconnect.
func (m *Manager) Close() error {
	defer func() { <-m.done }()
	return m.c.Close()
}

// backgroundReader listen to responses from the 1541u and prints them to stdout.
// It signals the m.done channel when the connection is closed or on error.
func (m *Manager) backgroundReader() {
	defer func() { m.done <- true }()
LOOP:
	for {
		var buf bytes.Buffer
		n, err := io.Copy(&buf, m.c)
		switch {
		case err == nil && n > 0:
			fmt.Println("[1541U] ", buf.String())
			continue LOOP
		case errors.Is(err, net.ErrClosed):
			fmt.Println("[1541U] Connection closed")
			return
		case errors.Is(err, io.EOF):
			fmt.Println("[1541U] EOF")
			return
		case err != nil:
			log.Printf("backgroundReader io.Copy failed: %v", err)
			return
		}
		log.Print("backgroundReader error, aborting read")
		return
	}
}
