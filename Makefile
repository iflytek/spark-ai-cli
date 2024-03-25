# Makefile for cross-compiling a Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=sparkai

# Build targets
.PHONY: all linux windows mac clean

all: linux windows mac

linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64

windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe

mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64

mac-m:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-darwin-arm64

clean:
	rm -f $(BINARY_NAME)-linux-amd64
	rm -f $(BINARY_NAME)-windows-amd64.exe
	rm -f $(BINARY_NAME)-darwin-amd64