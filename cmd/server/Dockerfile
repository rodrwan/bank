FROM golang:alpine AS builder

COPY . /app

WORKDIR /app

RUN go get ./...

RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]
