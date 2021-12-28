FROM ubuntu:20.04

RUN apt-get update --fix-missing \
    && apt-get install -y curl xz-utils

WORKDIR /borealys
COPY . /borealys
RUN mkdir /binaries
EXPOSE 5000

ENTRYPOINT ["./Main"]