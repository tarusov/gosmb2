export CGO_ENABLED=1
export GO111MODULE=off

# WINDOWS (MINGW)

build.windows.x64: export GOOS=windows
build.windows.x64: export GOARCH=amd64
build.windows.x64:
	go build -o ./bin/windows-x64.exe ./examples/cat/main.go
.PHONY: build.windows.x64

build.windows.x86: export GOOS=windows
build.windows.x86: export GOARCH=386
build.windows.x86:
	go build -o ./bin/windows-x86.exe ./examples/cat/main.go
.PHONY: build.windows.x86

# LINUX

build.linux.x64: export GOOS=linux
build.linux.x64: export GOARCH=amd64
build.linux.x64:
	go build -a -installsuffix cgo -ldflags '-s' -o ./bin/linux-x64 ./examples/cat/main.go
.PHONY: build.linux.x64

build.linux.x86: export GOOS=linux
build.linux.x86: export GOARCH=386
build.linux.x86:
	go build -a -installsuffix cgo -ldflags '-s' -o ./bin/linux-x86 ./examples/cat/main.go
.PHONY: build.linux.x86

# RUN

run.windows.x64: build.windows.x64
run.windows.x64:
	./bin/windows-x64.exe
.PHONY: run.windows.x64

run.windows.x86: build.windows.x86
run.windows.x86:
	./bin/windows-x86.exe
.PHONY: run.windows.x86

run.linux.x64: build.linux.x64
run.linux.x64:
	./bin/linux-x64
.PHONY: run.linux.x64

run.linux.x86: build.linux.x86
run.linux.x86:
	./bin/linux-x86
.PHONY: run.linux.x86