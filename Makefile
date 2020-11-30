# default make target.
release: fmt vet test build

build:
	go build -v ./...

fmt:
	go fmt ./...

test:
	go test -v ./...

vet:
	go vet -v ./...