FROM golang:latest AS builder

CMD mkdir -p /go/src/github.com/elliotcourant/noahdb
COPY ./ /go/src/github.com/elliotcourant/noahdb
WORKDIR /go/src/github.com/elliotcourant/noahdb
RUN go get -t -v ./...
RUN go build -o bin/noahdb

FROM ubuntu:18.04 AS final
WORKDIR /root/
COPY --from=builder /go/src/github.com/elliotcourant/noahdb/bin/noahdb ./noahdb
CMD ["./noahdb", "start"]
