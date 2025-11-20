.PHONY: build install clean

# Build the binary
build:
	go build -o bin/dif

# Install to user bin directory
install: build
	cp bin/dif $(HOME)/bin/dif

# Clean build artifacts
clean:
	rm -f bin/dif diffview dif

# Default target
all: install
