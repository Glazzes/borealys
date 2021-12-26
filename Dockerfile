FROM ubuntu:20.04

RUN mkdir /binaries
RUN apt-get update --fix-missing \
    && apt-get install -y curl xz-utils

WORKDIR /borealys
COPY . /borealys

EXPOSE 5000
ENTRYPOINT ["/bin/bash", "-c", "'./borealys/Main'"]