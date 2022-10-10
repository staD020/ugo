package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/staD020/ugo"
)

func main() {
	var (
		address        = "192.168.2.64:64"
		help           = false
		timeoutSeconds = 1
		mount          = false
	)
	if s := os.Getenv("UGO"); s != "" {
		address = s
	}
	flag.StringVar(&address, "a", address, "network address:port for the TCP connection to your 1541Ultimate")
	flag.IntVar(&timeoutSeconds, "timeout", timeoutSeconds, "connection timeout in seconds")
	flag.BoolVar(&mount, "m", mount, "always mount, never reset")
	flag.BoolVar(&help, "h", help, "help")
	flag.BoolVar(&help, "help", help, "show help")
	flag.Parse()
	ugo.DialTimeout = time.Duration(timeoutSeconds) * time.Second
	if help {
		flag.CommandLine.SetOutput(os.Stdout)
		printHelp()
		return
	}
	n := flag.NArg()
	if n < 1 {
		fmt.Printf("ugo %s by burg - a partial port of ucodenet to Go\n", ugo.Version)
		printUsage()
		return
	}

	u, err := ugo.New(address)
	if err != nil {
		log.Fatalf("ugo.New %q failed: %v", address, err)
	}
	defer u.Close()

	files, err := expandWildcards(flag.Args())
	if err != nil {
		log.Fatalf("expandWildcards %v failed: %v", flag.Args(), err)
	}
	if len(files) > 1 {
		if err = processMulti(u, files, mount); err != nil {
			log.Fatalf("processMulti %v failed: %v", files, err)
		}
		return
	}
	if err = process(u, files[0], mount); err != nil {
		log.Fatalf("process %q failed: %v", files[0], err)
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
		fmt.Printf("Press enter to mount next image %q\n", files[i])
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

func expandWildcards(filenames []string) (result []string, err error) {
	for _, filename := range filenames {
		if !strings.ContainsAny(filename, "?*") {
			result = append(result, filename)
			continue
		}
		dir := filepath.Dir(filename)
		ff, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("os.ReadDir %q failed: %w", dir, err)
		}
		name := filepath.Base(filename)
		for _, f := range ff {
			if f.IsDir() {
				continue
			}
			ok, err := filepath.Match(name, f.Name())
			if err != nil {
				return nil, fmt.Errorf("filepath.Match %q failed: %w", filename, err)
			}
			if ok {
				result = append(result, filepath.Join(dir, f.Name()))
			}
		}
	}
	return result, nil
}

func printUsage() {
	fmt.Println("usage: ./ugo [-h -a 192.168.2.64:64 -timeout 3] FILE [FILES]")
}

func printHelp() {
	fmt.Println()
	fmt.Println("# ugo", ugo.Version)
	fmt.Println()
	fmt.Println("Ugo provides 1541Ultimate control to run and mount C64 programs and disks via TCP.")
	fmt.Println("It is a partial port of [Ucodenet](https://csdb.dk/release/?id=189723) by TTL in pure Go by burg.")
	fmt.Println()
	fmt.Println("## Features")
	fmt.Println()
	fmt.Println(" - Resets, Mounts and runs prg/d64/d71/d81 files transparently.")
	fmt.Println(" - Supports multidisk and flip disk, just hit enter at the turn disk part.")
	fmt.Println(" - Force mount (no reset, no run) with the -m flag.")
	fmt.Println()
	fmt.Println("## Install Library")
	fmt.Println()
	fmt.Println("`go get github.com/staD020/ugo`")
	fmt.Println()
	fmt.Println("## Use Library")
	fmt.Println()
	fmt.Println("Error handling omitted, see source for more options.")
	fmt.Println()
	fmt.Println("```go")
	fmt.Println("package main")
	fmt.Println()
	fmt.Println("import (")
	fmt.Println("    \"os\"")
	fmt.Println("    \"github.com/staD020/ugo\"")
	fmt.Println(")")
	fmt.Println()
	fmt.Println("func main() {")
	fmt.Println("    f, _ := os.Open(\"file.prg\")")
	fmt.Println("    defer f.Close()")
	fmt.Println("    u, _ := ugo.New(\"192.168.2.64:64\")")
	fmt.Println("    defer u.Close()")
	fmt.Println("    _ = u.Run(f)")
	fmt.Println("    return")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()
	fmt.Println("## Install Command-line Interface")
	fmt.Println()
	fmt.Println("`go install github.com/staD020/ugo/cmd/ugo@latest`")
	fmt.Println()
	printUsage()
	fmt.Println()
	fmt.Println("You can also set your 1541u's address with environment variable UGO.")
	fmt.Println()
	fmt.Println("`export UGO=10.1.1.64:64`")
	fmt.Println()
	fmt.Println("## Options")
	fmt.Println()
	fmt.Println("```")
	flag.PrintDefaults()
	fmt.Println("```")
	fmt.Println()
}
