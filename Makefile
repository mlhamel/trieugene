build:
	go build -o bin/trieugene cmd/trieugene/main.go
	go build -o bin/rotondo services/rotondo/cmd/rotondo/main.go
	go build -o bin/rougecombien services/rougecombien/cmd/rougecombien/main.go

all: build