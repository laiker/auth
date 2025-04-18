FROM golang:1.23-alpine AS builder

COPY . /github.com/laiker/auth/
WORKDIR /github.com/laiker/auth/
RUN go mod download
RUN go build -v -o ./bin/auth ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/laiker/auth/.env .
COPY --from=builder /github.com/laiker/auth/service.key .
COPY --from=builder /github.com/laiker/auth/service.pem .
COPY --from=builder /github.com/laiker/auth/bin/auth .

CMD ["./auth"]