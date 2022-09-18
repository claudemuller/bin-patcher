phony: all

all: clean build gui

clean:
	rm -rf bin-patcher

build:
	go build -o bin-patcher cmd/main.go

cli:
	./bin-patcher -in $(in) -out $(out) -sig $(sig) -patch $(patch)

gui:
	./bin-patcher
