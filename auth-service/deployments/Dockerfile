FROM golang:1.23.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go build -o main ./cmd/server

FROM alpine:3
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

CMD ["/bin/sh", "-c", "goose -dir migrations postgres \"$DB_DSN\""]
CMD ["/app/main"]