BINARY_NAME=flightpathtracker
.DEFAULT_GOAL := run

build:
	go build -o ./target/${BINARY_NAME} ./src/main/main.go

run: build
	./target/${BINARY_NAME}

clean:
	go clean
	rm -f ./target/${BINARY_NAME}