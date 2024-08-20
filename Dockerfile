# Builder
FROM golang:1.22-alpine AS builder
RUN apk add --update make git curl

ARG MODULE_NAME=backend

COPY go.mod /home/${MODULE_NAME}/go.mod
COPY go.sum /home/${MODULE_NAME}/go.sum

WORKDIR /home/${MODULE_NAME}

COPY . /home/${MODULE_NAME}

RUN go build /home/${MODULE_NAME}/cmd/main/main.go

# Service
FROM alpine:latest as production
ARG MODULE_NAME=backend
WORKDIR /root/

COPY --from=builder /home/${MODULE_NAME}/config/dev/config.yaml config/dev/config.yaml
COPY --from=builder /home/${MODULE_NAME}/main .

RUN chown root:root main

CMD ["./main"]