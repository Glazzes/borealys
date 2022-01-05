FROM ubuntu:20.04

RUN apt-get update --fix-missing \
    && DEBIAN_FRONTEND="noninteractive" apt-get install -y \
    curl xz-utils unzip python3.10 -y

WORKDIR /borealys
COPY . .
ENTRYPOINT ["./Main"]