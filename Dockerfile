FROM ubuntu:20.04

RUN apt-get update --fix-missing \
    && DEBIAN_FRONTEND="noninteractive" apt-get install \
    curl xz-utils unzip golang python3.10 -y

WORKDIR /borealys
COPY . .
ENTRYPOINT ["./Main"]