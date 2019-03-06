FROM golang:latest

RUN mkdir -p /go/src/redeam-rest
WORKDIR /go/src/redeam-rest

ADD . /go/src/redeam-rest
RUN mkdir -p /docker-entrypoint-initdb.d
WORKDIR /go/src/redeam-rest/cmd/server
RUN go get
EXPOSE 8080
EXPOSE 5432
EXPOSE 9090