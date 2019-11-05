FROM golang:alpine

COPY . /usr/src/vutung2311-golang-test

WORKDIR /usr/src/vutung2311-golang-test

RUN go build -o /usr/local/bin/http internal/cmd/http/main.go

CMD /usr/local/bin/http.go