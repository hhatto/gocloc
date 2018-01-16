.PHONY: test build

build:
	go build -v ./
	go build cmd/gocloc/main.go

test:
	go test -v
