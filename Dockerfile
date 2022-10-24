FROM ubuntu:22.04

WORKDIR /usr/compiler-api

ARG GO_ENV
ARG PORT

ENV DEBIAN_FRONTEND noninteractive

RUN apt update

RUN apt-get install -y golang-go

RUN apt-get install -y nodejs

RUN apt-get install -y npm

RUN apt-get install -y python3

RUN npm i -g typescript ts-node

COPY go.mod .
COPY go.sum .

RUN ["go", "mod", "download"]

COPY . .

RUN ["go", "build", "-o", "/build"]

CMD ["/build"]
