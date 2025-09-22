.PHONY: all generate test bench

all: generate test bench

generate:
	go generate ./examples

test:
	go test -v ./examples

bench:
	go test -v -bench=. -run=^$$ ./examples
