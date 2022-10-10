GOBUILDFLAGS=-v -trimpath
LDFLAGS=-s -w
TARGET=ugo
CROSS_TARGETS=ugo_linux_amd64 ugo_darwin_arm64 ugo_darwin_amd64 ugo_win_amd64.exe ugo_win_x86.exe
SRC=*.go cmd/ugo/main.go
CGO=0

all: $(TARGET)

$(TARGET): $(SRC)
	CGO_ENABLED=$(CGO) go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o ./ ./cmd/...

cross: $(CROSS_TARGETS)

dist: $(CROSS_TARGETS)
	mkdir ugo-0.1
	cp ugo_*6* ugo-0.1/
	cp README.md ugo-0.1/
	zip -r ugo-0.1.zip ugo-0.1

test: $(TARGET)
	go test -v -cover -race
	./$(TARGET) -a localhost:6464 testdata/evoluer.prg

install: $(TARGET)
	sudo cp $(TARGET) /usr/local/bin/

ugo_linux_amd64: $(SRC)
	CGO_ENABLED=$(CGO) go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o $@ ./cmd/ugo/

ugo_darwin_arm64: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=darwin GOARCH=arm64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o $@ ./cmd/ugo/

ugo_darwin_amd64: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=darwin GOARCH=amd64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o $@ ./cmd/ugo/

ugo_win_amd64.exe: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=windows GOARCH=amd64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o $@ ./cmd/ugo/

ugo_win_x86.exe: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=windows GOARCH=386 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o $@ ./cmd/ugo/

clean:
	rm -f $(TARGET) $(CROSS_TARGETS)
