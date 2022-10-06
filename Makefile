GOBUILDFLAGS=-v -trimpath
LDFLAGS=-s -w
TARGET=ugo
SRC=*.go cmd/ugo/main.go
CGO=0

all: $(TARGET)

$(TARGET): $(SRC)
	CGO_ENABLED=$(CGO) go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o ./ ./cmd/...

cross: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=darwin GOARCH=arm64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o ugo_darwin_arm64 ./cmd/ugo/

test: $(TARGET)
	go test -v -cover -race
	./$(TARGET) -a localhost:6464 testdata/evoluer.prg

install: $(TARGET)
	sudo cp $(TARGET) /usr/local/bin/

clean:
	rm -f $(TARGET) ugo_darwin_arm64
