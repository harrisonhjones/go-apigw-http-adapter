# default make target.
.PHONY: release
release: fmt vet test build

.PHONY: build
build:
	go build -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: vet
vet:
	go vet -v ./...
