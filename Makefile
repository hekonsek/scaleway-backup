PACKAGES := \
	github.com/hekonsek/scaleway-backup

all: format build silent-test

build:
	go build -o bin/scaleway-backup scaleway-backup.go

test:
	go test -v $(PACKAGES)

silent-test:
	go test $(PACKAGES)

format:
	go fmt $(PACKAGES)