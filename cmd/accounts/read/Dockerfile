FROM golang

COPY . /app

WORKDIR /app

RUN go get ./...

RUN go build -o server ./cmd/accounts/read

CMD [ "./server" ]
