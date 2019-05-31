.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/wtc/go-wtc/cmd/swarm
	go build -o bin/gwtc     github.com/wtc/go-wtc/cmd/gwtc
	go build -o bin/bootnode github.com/wtc/go-wtc/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/wtc/go-wtc/swarm/...
