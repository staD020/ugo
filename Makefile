GOBUILDFLAGS=-v -trimpath
LDFLAGS=-s -w
TARGET=ultim8
SRC=*.go cmd/ultim8/main.go

all: $(TARGET)

$(TARGET): $(SRC)
	go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o ./ ./cmd/...

test: all
	./$(TARGET) testdata/evoluer.prg

install:
	sudo cp $(TARGET) /usr/local/bin/

clean:
	rm -f $(TARGET)
