FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache  \
    add ca-certificates \
    libc6-compat

COPY --from=builder /app/cmd/server/main .
COPY --from=builder /app/.env .

RUN chmod +x ./main

EXPOSE 8080
CMD ["./main"]