.PHONY: build run test vet clean verify

build:
	go build -o bin/runtimewall-api ./cmd/server

run: build
	./bin/runtimewall-api

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf bin/

verify: vet test build
	@echo "Build and tests passed."
