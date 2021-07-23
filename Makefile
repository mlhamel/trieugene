build:
	go build -o bin/trieugene cmd/trieugene/main.go
	go build -o bin/rotondo services/rotondo/cmd/rotondo/main.go
	go build -o bin/rougecombien services/rougecombien/cmd/rougecombien/main.go
	cargo build --target-dir=bin/ --bin=trieugene-rust

all: build