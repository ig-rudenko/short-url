FROM golang:1.23.1-alpine AS builder

LABEL authors="ig-rudenko"

WORKDIR /app

COPY go.* /app/

RUN go mod download

COPY . /app/

# Build the Go application with CGO_ENABLED=0 to ensure a statically linked binary
# The binary will be named caching-proxy and located in /app
RUN CGO_ENABLED=0 go build -o short-url /app/cmd/app/main.go


FROM alpine

COPY --from=builder /app/short-url /app/short-url

WORKDIR /app

EXPOSE 8000

ENTRYPOINT ["./short-url"]