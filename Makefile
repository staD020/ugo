GOBUILDFLAGS=-v -trimpath
LDFLAGS=-s -w
TARGET=ultim8
SRC=*.go cmd/ultim8/main.go
CGO=0

all: $(TARGET)

$(TARGET): $(SRC)
	CGO_ENABLED=$(CGO) go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o ./ ./cmd/...

cross: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=darwin GOARCH=arm64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o ultim8_darwin_arm64 ./cmd/ultim8/

test: $(TARGET)
	go test -v -cover -race
	./$(TARGET) -a localhost:6464 testdata/evoluer.prg

install:
	sudo cp $(TARGET) /usr/local/bin/

clean:
	rm -f $(TARGET) ultim8_darwin_arm64
