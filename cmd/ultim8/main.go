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

	"github.com/staD020/ultim8"
)

func printUsage() {
	fmt.Printf("ultim8 %s by burg - a partial port of ucodenet to Go\n", ultim8.Version)
	fmt.Println("usage: ./ultim8 [-h -a 127.0.0.1:64 -timeout 3] FILE [FILES]")
}

func main() {
	var (
		address        = "192.168.2.64:64"
		timeoutSeconds = 1
		mount          = false
	)
	if s := os.Getenv("ULTIM8"); s != "" {
		address = s
	}
	flag.StringVar(&address, "a", address, "network address:port for the TCP connection to your 1541Ultimate")
	flag.IntVar(&timeoutSeconds, "timeout", timeoutSeconds, "connection timeout in seconds")
	flag.BoolVar(&mount, "m", mount, "always mount, never reset")
	flag.Parse()
	ultim8.DialTimeout = time.Duration(timeoutSeconds) * time.Second
	n := flag.NArg()
	if n < 1 {
		printUsage()
		return
	}
	u, err := ultim8.New(address)
	if err != nil {
		log.Fatalf("ultim8.New %q failed: %v", address, err)
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

func process(u *ultim8.Manager, path string, mount bool) error {
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

func processMulti(u *ultim8.Manager, files []string, mount bool) error {
	f, err := os.Open(files[0])
	if err != nil {
		return fmt.Errorf("os.Open %q failed: %v", files[0], err)
	}
	defer f.Close()
	if mount {
		fmt.Printf("Multi mode, mounting image %q\n", files[0])
		if err = u.Mount(f); err != nil {
			return fmt.Errorf("u.Mount failed: %v", err)
		}
	} else {
		fmt.Printf("Multi mode, starting image %q\n", files[0])
		if err = u.Run(f); err != nil {
			return fmt.Errorf("u.Run failed: %v", err)
		}
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
