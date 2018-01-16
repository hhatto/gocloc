.PHONY: test build

build:
	mkdir -p bin
	go build -o ./bin/gocloc cmd/gocloc/main.go

update-package:
	go get -u github.com/hhatto/gocloc

test:
	go test -v
