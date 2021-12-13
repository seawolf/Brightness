.DEFAULT_GOAL := toggle

.PHONY: test test_coverage toggle up down clean

all: test toggle

test:
	@ go test -coverprofile=cover.out .

test_coverage: test
	@ go tool cover -html=cover.out

brightness:
	@ go build -o ./brightness brightness.go

toggle: brightness
	@ ./brightness

up: brightness
	@ ./brightness up

down: brightness
	@ ./brightness down

clean:
	@ test -f ./brightness && rm ./brightness || true
