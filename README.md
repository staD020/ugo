# ultim8 0.1

Ultim8 provides 1541Ultimate control to run and mount C64 programs and disks via TCP.
It is a partial port of [Ucodenet](https://csdb.dk/release/?id=189723) by TTL in pure Go by burg.

## Features

 - Resets, Mounts and runs prg/d64/d71/d81 files transparently.
 - Supports multidisk and flip disk, just hit enter at the turn disk part.
 - Force mount (no reset, no run) with the -m flag.

## Install Library

`go get github.com/staD020/ultim8`

## Install Command-line Interface

`go install github.com/staD020/ultim8@latest`
