.PHONY: all build test coverage coverage-html coverage-show

abadgopath=/go/src/github.com/NeowayLabs/abad
runabad=docker run -v `pwd`:$(abadgopath) -w $(abadgopath)

all: build test

build:
	go build -o ./cmd/abad/abad -v ./cmd/abad 

test:
	go test -race -v ./... -timeout=30s

coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	
coverage-html: coverage
	go tool cover -html=coverage.txt -o coverage.html
	@echo "coverage file: coverage.html"

coverage-show: coverage-html
	xdg-open coverage.html
	
analysis:
	go get honnef.co/go/tools/cmd/megacheck
	megacheck ./...

vendor:
	go get github.com/madlambda/vendor
	vendor
	
devimg=neowaylabs/abadtest
devimage:
	docker build . -t $(devimg)
	
devshell: devimage
	$(runabad) -ti $(devimg)