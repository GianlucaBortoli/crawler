.PHONY: build
build:
	go build -race -o ./bin/crawler

.PHONY: test
test:
	go test -race ./...

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: lint
lint:
	golangci-lint run --tests=false
