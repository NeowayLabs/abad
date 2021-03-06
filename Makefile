.PHONY: all vendor build test coverage coverage-html coverage-show

abadgopath=/go/src/github.com/NeowayLabs/abad
runabad=docker run -v `pwd`:$(abadgopath) -w $(abadgopath)
installdir?=/usr/local/bin

all: fmt build analysis test dev-test-e2e

build:
	cd cmd/abad && go build -v

install: build
	go install -v ./cmd/abad

fmt:
	gofmt -s -w .

test:
	go test -failfast -race -v ./... -timeout=30s

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
	go get golang.org/x/lint/golint
	go get honnef.co/go/tools/cmd/megacheck
	megacheck ./...
	# FIXME: right now we have to much undocumented stuff =(
	# golint ./...

vendor:
	go get github.com/madlambda/vendor
	vendor

devimgversion=0.1
devimg=neowaylabs/abadev:$(devimgversion)
devimage:
	docker build . -t $(devimg)
	
publish-devimage: devimage
	docker push $(devimg)
	
dev-shell:
	$(runabad) -ti $(devimg)
	
dev-d8:
	$(runabad) $(devimg) d8 $(abadgopath)/$(code)
	
dev-test-e2e:
	$(runabad) $(devimg) make install test-e2e
