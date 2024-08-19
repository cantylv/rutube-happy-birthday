# Builder
FROM golang:1.22-alpine AS builder
RUN apk add --update make git curl

ARG MODULE_NAME=backend

COPY go.mod /home/${MODULE_NAME}/go.mod
COPY go.sum /home/${MODULE_NAME}/go.sum

WORKDIR /home/${MODULE_NAME}

COPY . /home/${MODULE_NAME}

RUN go run /home/${MODULE_NAME}/cmd/main/main.go