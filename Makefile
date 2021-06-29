build:
	go generate ./...
	go build

tests: build
	go test -v ./...

lint: build
	golangci-lint run

cover: build
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out

.phony: build tests lint cover
