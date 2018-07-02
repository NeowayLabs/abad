.PHONY: all build test coverage coverage-html coverage-show

abadgopath=/go/src/github.com/NeowayLabs/abad
runabad=docker run -v `pwd`:$(abadgopath) -w $(abadgopath)
installdir?=/usr/local/bin

all: build test analysis

build:
	go build -o ./cmd/abad/abad -v ./cmd/abad

install: build
	cp ./cmd/abad/abad $(installdir)

test:
	go test -race -v ./... -timeout=30s

test-e2e:
	go test -v ./tests/e2e -tags e2e

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
	
devimg=neowaylabs/abadev
devimage:
	docker build . -t $(devimg)
	
dev-shell: devimage
	$(runabad) -ti $(devimg)
	
dev-test-e2e: devimage
	$(runabad) -ti $(devimg) make install && /bin/sh