all: bench

.PHONY: all

test:
	go test --race ./...

bench: test
	go test -bench=. -benchmem ./...
