
# ugo 0.2

Ugo provides 1541Ultimate control to run and mount C64 programs and disks via TCP.
It is a partial port of [Ucodenet](https://csdb.dk/release/?id=189723) by TTL in pure Go by burg.

## Features

 - Mounts and runs .prg and .d64 files transparently.
 - Supports multidisk and flip disk, just hit enter at the turn disk part.
 - Force mount (no run) with the -m flag.

## Install Library

`go get github.com/staD020/ugo`

## Use Library

Error handling omitted, see source for more options.

```go
package main

import (
    "os"
    "github.com/staD020/ugo"
)

func main() {
    f, _ := os.Open("file.prg")
    defer f.Close()
    u, _ := ugo.New("192.168.2.64:64")
    defer u.Close()
    _ = u.Run(f)
    return
}
```

## Install Command-line Interface

`go install github.com/staD020/ugo/cmd/ugo@latest`

usage: ./ugo [-h -a 192.168.2.64:64 -timeout 3] FILE [FILES]

You can also set your 1541u's address with environment variable UGO.

`export UGO=10.1.1.64:64`

## Options

```
  -a string
    	network address:port for the TCP connection to your 1541Ultimate (default "localhost:6464")
  -h	help
  -help
    	show help
  -m	always mount, never reset
  -timeout int
    	connection timeout in seconds (default 1)
```

