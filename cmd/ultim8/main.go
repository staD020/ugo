package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/staD020/ultim8"
)

func main() {
	var (
		address        string
		timeoutSeconds int
	)
	flag.StringVar(&address, "a", "", "address")
	flag.StringVar(&address, "address", "192.168.2.64:64", "network address:port for the TCP connection to your 1541Ultimate")
	flag.IntVar(&timeoutSeconds, "timeout", 1, "connection timeout")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Incorrect number of arguments: %d", flag.NArg())
	}
	path := flag.Args()[0]
	ultim8.DialTimeout = time.Duration(timeoutSeconds) * time.Second

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("os.Open %q failed: %v", path, err)
	}
	defer f.Close()

	u, err := ultim8.New(address)
	if err != nil {
		log.Fatalf("New %q failed: %v", address, err)
	}
	defer u.Close()
	if err = u.RunPrg(f); err != nil {
		log.Fatalf("u.RunPrg failed: %v", err)
	}
	return
}
