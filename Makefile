BINARY_NAME = main
ENTRYPOINT = cmd/main/main.go

all: build test

build:
	go build -o ${BINARY_NAME} ${ENTRYPOINT}

tidy:
	go fmt ./...
	go mod tidy -v

easyjs:
	easyjson -no_std_marshalers -all internal/entity
