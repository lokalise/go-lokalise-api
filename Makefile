.DEFAULT_GOAL := build

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./...

GODEP := /usr/local/bin/dep

$(GODEP):
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

.PHONY: dependencies
dependencies: $(GODEP)
	dep ensure
