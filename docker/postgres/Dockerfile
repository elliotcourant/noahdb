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

FROM debian:latest
RUN mkdir /postgres
RUN mkdir /db
COPY --from=postgres-builder /postbuild /postgres
EXPOSE 5432
ENV LD_LIBRARY_PATH=/postgres/lib
RUN groupadd -g 999 postgres && \
    useradd -r -u 999 -g postgres postgres
RUN chown -R postgres /db
RUN chmod -R 750 /db
USER postgres
RUN /postgres/bin/initdb -D /db
COPY ./pg_hba.conf /db/pg_hba.conf
COPY ./postgresql.conf /db/postgresql.conf
CMD ["/postgres/bin/postgres", "-D", "/db"]