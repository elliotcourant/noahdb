FROM golang:latest AS builder

CMD mkdir -p /go/src/github.com/elliotcourant/noahdb
COPY ./ /go/src/github.com/elliotcourant/noahdb
WORKDIR /go/src/github.com/elliotcourant/noahdb
RUN go build -o bin/noahdb



FROM ubuntu:18.04 AS final
RUN mkdir /node
WORKDIR /node
RUN mkdir -p /data/db

COPY --from=noahdb/postgres:local /postbuild /node
ENV LD_LIBRARY_PATH=/node/lib
ENV PGPASSWORD=""
ENV PGPORT=5432
ENV PGUSER=postgres
RUN groupadd -g 999 postgres && \
    useradd -r -u 999 -g postgres postgres
RUN chown -R postgres /node
RUN chown -R postgres /data/db
RUN chmod -R 750 /data/db
USER postgres

RUN /node/bin/initdb -D /data/db
VOLUME /data/db
COPY ./pg_hba.conf /data/db/pg_hba.conf
COPY ./postgresql.conf /data/db/postgresql.conf
COPY --from=builder /go/src/github.com/elliotcourant/noahdb/bin/noahdb /node/bin/noahdb

EXPOSE 5433 5432
COPY ./run.sh /run.sh
CMD ["/run.sh"]
