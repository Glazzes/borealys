FROM ubuntu:20.04
RUN apt update && apt install apt-utils -y && apt install curl -y
WORKDIR /borealys
COPY . /borealys
ENTRYPOINT ["bash", "-c", "./Initializer"]