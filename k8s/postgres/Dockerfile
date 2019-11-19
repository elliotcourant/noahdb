FROM ubuntu:19.04 as noahdb-postgres
ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update
RUN apt-get --assume-yes install \
#    libreadline6 \
    libreadline6-dev \
    git-all \
    build-essential \
    zlib1g-dev \
    libossp-uuid-dev \
    flex \
    bison \
    libxml2-utils \
    xsltproc
RUN git clone https://github.com/postgres/postgres.git && cd ./postgres && git checkout REL_12_STABLE && cd ..
RUN mkdir /postbuild
RUN ./postgres/configure --prefix=/postbuild --with-extra-version=noahdb --with-ossp-uuid
RUN cd /postgres
RUN make
RUN make install