FROM golang:latest AS builder

CMD mkdir -p /go/src/github.com/elliotcourant/noahdb
COPY ./ /go/src/github.com/elliotcourant/noahdb
WORKDIR /go/src/github.com/elliotcourant/noahdb
RUN go get -t -v ./...
RUN go build -o bin/noahdb

FROM ubuntu:16.04 as postgres-builder
RUN apt-get update
RUN apt-get --assume-yes install \
    libreadline6 \
    libreadline6-dev \
    git-all \
    build-essential \
    zlib1g-dev \
    libossp-uuid-dev \
    flex \
    bison \
    libxml2-utils \
    xsltproc
RUN git clone https://github.com/readystock/postgres.git && cd ./postgres && git checkout REL_10_STABLE && cd ..
RUN mkdir /postbuild
RUN ./postgres/configure --prefix=/postbuild --with-ossp-uuid
RUN cd /postgres
RUN make
RUN make install

FROM ubuntu:18.04 AS final
RUN mkdir /node
WORKDIR /node
RUN mkdir /node/pgdata
COPY --from=postgres-builder /postbuild /node
EXPOSE 5432
ENV LD_LIBRARY_PATH=/node/lib
ENV PGPASSWORD=""
ENV PGPORT=5432
ENV PGUSER=postgres
RUN groupadd -g 999 postgres && \
    useradd -r -u 999 -g postgres postgres
RUN chown -R postgres /node
RUN chmod -R 750 /node/pgdata
USER postgres
RUN /node/bin/initdb -D /node/pgdata
COPY ./pg_hba.conf /node/pgdata/pg_hba.conf
COPY ./postgresql.conf /node/pgdata/postgresql.conf
COPY --from=builder /go/src/github.com/elliotcourant/noahdb/bin/noahdb /node/bin/noahdb
EXPOSE 5433
COPY ./run.sh /run.sh
CMD ["/run.sh"]
