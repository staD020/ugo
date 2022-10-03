package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/staD020/ultim8"
)

func printUsage() {
	fmt.Println("usage: ./ultim8 [-h -a 127.0.0.1:64 -timeout 3] FILE")
}

func main() {
	var (
		address        = "192.168.2.64:64"
		timeoutSeconds int
	)
	if s := os.Getenv("ULTIM8"); s != "" {
		address = s
	}
	flag.StringVar(&address, "a", address, "network address:port for the TCP connection to your 1541Ultimate")
	flag.IntVar(&timeoutSeconds, "timeout", 1, "connection timeout")
	flag.Parse()
	n := flag.NArg()
	if n < 1 {
		printUsage()
		return
	}
	ultim8.DialTimeout = time.Duration(timeoutSeconds) * time.Second
	path := flag.Args()[0]

	u, err := ultim8.New(address)
	if err != nil {
		log.Fatalf("ultim8.New %q failed: %v", address, err)
	}
	defer u.Close()
	if err = process(u, path); err != nil {
		log.Fatalf("process %q failed: %v", path, err)
	}
	return
}

func process(u *ultim8.Manager, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os.Open %q failed: %v", path, err)
	}
	defer f.Close()
	if err = u.Run(f); err != nil {
		return fmt.Errorf("u.Run failed: %v", err)
	}
	return nil
}
