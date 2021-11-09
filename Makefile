.DEFAULT_GOAL := run

.PHONY: test test_coverage run clean

all: test run

test:
	@ go test -coverprofile=cover.out .

test_coverage: test
	@ go tool cover -html=cover.out

brightness:
	@ go build -o ./brightness brightness.go

run: brightness
	@ ./brightness

clean:
	@ test -f ./brightness && rm ./brightness || true
