FROM ubuntu:20.04
SHELL ["/bin/bash", "-c"]

RUN apt-get update --fix-missing \
    && apt-get install -y curl xz-utils

RUN mkdir /binaries

WORKDIR /borealys
COPY . /borealys

RUN ./Initializer
ENTRYPOINT ["/bin/bash", "-c"]