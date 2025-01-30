FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

WORKDIR /app/cmd

# TODO: dodać obsługę innych systemów
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main /app/main

WORKDIR /app

CMD ["./main"]
