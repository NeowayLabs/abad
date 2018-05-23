.PHONY: all build test coverage coverage-html coverage-show

all: build test

build:
	go build -o ./cmd/abad/abad -v ./cmd/abad 

test:
	go test -race -v ./...

coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	
coverage-html: coverage
	go tool cover -html=coverage.txt -o coverage.html
	@echo "coverage file: coverage.html"

coverage-show: coverage-html
	xdg-open coverage.html

vendor:
	go get github.com/madlambda/vendor
	vendor