BINARY := whoisusing

.PHONY: all build clean fmt fmt-check test race vet check install run

all: build

build:
	go build -o $(BINARY) .

clean:
	rm -f $(BINARY)
	rm -rf dist

fmt:
	go fmt ./...

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:"; gofmt -l .; exit 1)

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

check: fmt-check test race vet

install:
	go install .

run:
	go run . $(ARGS)
