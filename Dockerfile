FROM ubuntu:22.04

WORKDIR /usr/compiler-api

ENV GO_ENV="production"
ENV PORT ${PORT}

ENV DEBIAN_FRONTEND noninteractive

ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION 18.12.1

SHELL ["/bin/bash", "-c"]

RUN tee /etc/apt/sources.list.d/mono-official-stable.list

RUN apt update

RUN apt-get install -y golang-go

RUN apt-get install -y zip unzip curl

RUN curl --silent -o- https://raw.githubusercontent.com/creationix/nvm/v0.31.2/install.sh | bash
RUN . "$NVM_DIR/nvm.sh" && nvm install $NODE_VERSION
RUN . "$NVM_DIR/nvm.sh" && nvm use v$NODE_VERSION
RUN . "$NVM_DIR/nvm.sh" && nvm alias default v$NODE_VERSION
ENV PATH="/root/.nvm/versions/node/v$NODE_VERSION/bin/:$PATH"

RUN apt-get install -y python3

RUN apt-get install -y default-jre

RUN apt-get install -y default-jdk

RUN apt-get install -y rustc

RUN apt-get install -y build-essential

RUN apt-get install -y mono-mcs

RUN apt-get install -y python3-pip

RUN pip3 install numpy

RUN curl -s "https://get.sdkman.io" | bash

RUN source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install kotlin

ENV PATH=/root/.sdkman/candidates/kotlin/current/bin:$PATH

RUN npm i -g typescript ts-node

COPY go.mod .
COPY go.sum .

RUN ["go", "mod", "download"]

COPY . .

RUN ["go", "build", "-buildvcs=false", "-o", "/build"]

EXPOSE ${PORT}

CMD ["/build"]
