FROM ubuntu:20.04

RUN apt-get update --fix-missing \
    && apt-get install -y curl xz-utils

RUN DEBIAN_FRONTEND="noninteractive" apt-get install python3.10 -y

WORKDIR /borealys
COPY . .
ENTRYPOINT ["./Main"]