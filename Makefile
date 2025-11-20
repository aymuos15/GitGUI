.PHONY: build install clean

# Build the binary
build:
	go build -o gg

# Install to user bin directory
install: build
	cp gg $(HOME)/bin/gg

# Clean build artifacts
clean:
	rm -f gg

# Default target
all: install
