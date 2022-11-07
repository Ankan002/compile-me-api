FROM ubuntu:22.04

WORKDIR /usr/compiler-api

ENV GO_ENV ${GO_ENV}
ENV PORT ${PORT}

ENV DEBIAN_FRONTEND noninteractive

RUN apt update

RUN apt-get install -y golang-go

RUN apt-get install -y nodejs

RUN apt-get install -y npm

RUN apt-get install -y python3

RUN apt-get install -y default-jre

RUN apt-get install -y default-jdk

RUN apt-get install -y rustc

RUN npm i -g typescript ts-node

COPY go.mod .
COPY go.sum .

RUN ["go", "mod", "download"]

COPY . .

RUN ["go", "build", "-o", "/build"]

EXPOSE ${PORT}

CMD ["/build"]
