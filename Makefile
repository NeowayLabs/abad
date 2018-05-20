.PHONY: all build test

all: build test

build:
	go build -o ./cmd/abad/abad -v ./cmd/abad 

test:
	go test -race -v ./... 