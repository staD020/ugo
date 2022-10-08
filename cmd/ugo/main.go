package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/staD020/ugo"
)

func printUsage() {
	fmt.Printf("ugo %s by burg - a partial port of ucodenet to Go\n", ugo.Version)
	fmt.Println("usage: ./ugo [-h -a 192.168.2.64:64 -timeout 3] FILE [FILES]")
}

func main() {
	var (
		address        = "192.168.2.64:64"
		timeoutSeconds = 1
		mount          = false
	)
	if s := os.Getenv("UGO"); s != "" {
		address = s
	}
	flag.StringVar(&address, "a", address, "network address:port for the TCP connection to your 1541Ultimate")
	flag.IntVar(&timeoutSeconds, "timeout", timeoutSeconds, "connection timeout in seconds")
	flag.BoolVar(&mount, "m", mount, "always mount, never reset")
	flag.Parse()
	ugo.DialTimeout = time.Duration(timeoutSeconds) * time.Second
	n := flag.NArg()
	if n < 1 {
		printUsage()
		return
	}
	u, err := ugo.New(address)
	if err != nil {
		log.Fatalf("ugo.New %q failed: %v", address, err)
	}
	defer u.Close()

	if n > 1 {
		if err = processMulti(u, flag.Args(), mount); err != nil {
			log.Fatalf("processMulti %v failed: %v", flag.Args(), err)
		}
		return
	}

	path := flag.Args()[0]
	if err = process(u, path, mount); err != nil {
		log.Fatalf("process %q failed: %v", path, err)
	}
	return
}

func process(u *ugo.Manager, path string, mount bool) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os.Open %q failed: %v", path, err)
	}
	defer f.Close()
	if mount {
		if err = u.Mount(f); err != nil {
			return fmt.Errorf("u.Mount failed: %v", err)
		}
		return nil
	}
	if err = u.Run(f); err != nil {
		return fmt.Errorf("u.Run failed: %v", err)
	}
	return nil
}

func processMulti(u *ugo.Manager, files []string, mount bool) error {
	f, err := os.Open(files[0])
	if err != nil {
		return fmt.Errorf("os.Open %q failed: %v", files[0], err)
	}
	defer f.Close()
	var fn func(io.Reader) error
	a := "u.Run"
	fn = u.Run
	if mount {
		a = "u.Mount"
		fn = u.Mount
	}
	fmt.Printf("Multi mode, %s image %q\n", a, files[0])
	if err = fn(f); err != nil {
		return fmt.Errorf("%s failed: %v", a, err)
	}

	r := bufio.NewReader(os.Stdin)
	for i := 1; i < len(files); i++ {
		fmt.Printf("Press enter to mount next image %q", files[i])
		_, err = r.ReadString('\n')
		switch {
		case errors.Is(err, io.EOF):
			fmt.Println("EOF")
			return nil
		case err != nil:
			return fmt.Errorf("r.ReadString failed: %w", err)
		}
		f, err := os.Open(files[i])
		if err != nil {
			return fmt.Errorf("os.Open %q failed: %w", files[i], err)
		}
		defer f.Close()
		if err = u.Mount(f); err != nil {
			return fmt.Errorf("u.Mount %q failed: %w", files[i], err)
		}
	}
	return nil
}
