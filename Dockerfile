FROM golang:1.9

RUN mkdir -p /go/src/redeam-rest
WORKDIR /go/src/redeam-rest

ADD . /go/src/redeam-rest
RUN mkdir -p /docker-entrypoint-initdb.d
WORKDIR /go/src/redeam-rest/cmd/server
RUN go get