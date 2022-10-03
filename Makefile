GOBUILDFLAGS=-v -trimpath
LDFLAGS=-s -w
TARGET=ultim8
SRC=*.go cmd/ultim8/main.go

all: $(TARGET)

$(TARGET): $(SRC)
	go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o ./ ./cmd/...

cross: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=darwin GOARCH=arm64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o ultim8_darwin_arm64 ./cmd/ultim8/

test: all
	./$(TARGET) testdata/evoluer.prg

install:
	sudo cp $(TARGET) /usr/local/bin/

clean:
	rm -f $(TARGET) ultim8_darwin_arm64
