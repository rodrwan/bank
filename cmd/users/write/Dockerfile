FROM golang

COPY . /app

WORKDIR /app

RUN go get ./...

RUN go build -o server ./cmd/users/write

CMD [ "./server" ]
