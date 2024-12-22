FROM golang:1.20.3-alpine AS builder

COPY . /github.com/laiker/auth/
WORKDIR /github.com/laiker/auth/

RUN ls -la
RUN pwd
RUN go mod download
RUN go build -v -o ./bin/auth ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/laiker/auth/bin/auth .

CMD ["./auth"]