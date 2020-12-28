FROM golang:alpine

COPY . /app

WORKDIR /app

EXPOSE 3000

RUN go get ./...
