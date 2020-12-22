export CGO_ENABLED=1
export GO111MODULE=off

GCC_PREFIX = /usr
CC_CMD = gcc
LD_CMD = gcc
CC = $(GCC_PREFIX)/bin/$(CC_CMD)
LD = $(GCC_PREFIX)/bin/$(LD_CMD)

build.example:
	go build  -o ./bin/example ./example/main.go
.PHONY: build.example

run.example: build.example
run.example: export LD_LIBRARY_PATH=./lib
run.example:
	./bin/example
.PHONY: run.example