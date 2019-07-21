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
RUN git clone https://github.com/postgres/postgres.git && cd ./postgres && git checkout REL_10_STABLE && cd ..
RUN mkdir /postbuild
RUN ./postgres/configure --prefix=/postbuild --with-ossp-uuid
RUN cd /postgres
RUN make
RUN make install