phony: all

all: clean build run

clean:
	rm -rf bin-patcher

build:
	go build -o bin-patcher cmd/main.go

run:
	./bin-patcher -in $(in) -out $(out) -sig $(sig) -patch $(patch)
