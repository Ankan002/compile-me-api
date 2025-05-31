FROM ubuntu:22.04

WORKDIR /usr/compiler-api

ENV GO_ENV="production"
ENV PORT ${PORT}
ENV INVOCATION ${INVOCATION}
ENV CGO_ENABLED=0
ENV GOOS=linux

ARG BUILD_ARCH=${BUILD_ARCH}
ENV GO_ARCH=$BUILD_ARCH

ENV DEBIAN_FRONTEND noninteractive

ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION 18.12.1

SHELL ["/bin/bash", "-c"]

RUN echo "BUILD_ARCH is: $BUILD_ARCH"
RUN echo "GO_ARCH is: $GO_ARCH"

RUN tee /etc/apt/sources.list.d/mono-official-stable.list

RUN mono --version

RUN apt update
RUN apt-get install -y zip unzip curl wget

# RUN apt-get install -y golang-go

RUN wget https://go.dev/dl/go1.24.3.linux-$BUILD_ARCH.tar.gz
RUN tar -C /usr/local -xzf go1.24.3.linux-$BUILD_ARCH.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"

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

RUN pip3 --no-cache-dir install numpy

RUN curl -s "https://get.sdkman.io" | bash

RUN source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install kotlin

ENV PATH=/root/.sdkman/candidates/kotlin/current/bin:$PATH

RUN npm i -g typescript ts-node
RUN echo "Building with Go version:" && go version

COPY go.mod .
COPY go.sum .

RUN ["go", "mod", "download"]

COPY . .

# RUN ["go", "build", "-buildvcs=false", "-o", "/build"]
# -tags musl --ldflags "-extldflags -static"

RUN go build -buildvcs=false -o bootstrap .
RUN mkdir -p /var/runtime && mv bootstrap /var/runtime/bootstrap && chmod +x /var/runtime/bootstrap
WORKDIR /var/runtime

EXPOSE 5000-10000

CMD ["bootstrap"]
#

# FROM public.ecr.aws/amazonlinux/amazonlinux:2023

# WORKDIR /usr/compiler-api

# ENV GO_ENV="production"
# ENV PORT=${PORT}
# ENV INVOCATION=${INVOCATION}
# ENV CGO_ENABLED=0
# ENV GOOS=linux
# ARG BUILD_ARCH=arm64
# ENV GO_ARCH=$BUILD_ARCH
# ENV NODE_VERSION=18.12.1
# ENV NVM_DIR=/root/.nvm
# ENV SDKMAN_DIR="/root/.sdkman"

# # Extend PATH for Go, Node, Kotlin
# ENV PATH="/usr/local/go/bin:$PATH:/root/.nvm/versions/node/v$NODE_VERSION/bin:$PATH:$SDKMAN_DIR/candidates/kotlin/current/bin"

# SHELL ["/bin/bash", "-c"]

# RUN echo "BUILD_ARCH is: $BUILD_ARCH"
# RUN echo "GO_ARCH is: $GO_ARCH"

# # RUN dnf install -y curl

# # -----------------------
# # Install essential dev tools and languages
# # -----------------------
# RUN dnf groupinstall -y "Development Tools" && \
#     dnf install -y \
#     gcc gcc-c++ glibc-devel \
#     rust cargo \
#     git cmake make automake autoconf libtool libtoolize \
#     gettext glibc-devel glib2-devel zlib-devel libcurl-devel \
#     wget unzip tar gzip \
#     java-17-amazon-corretto python3 python3-pip

# # -----------------------
# # Install Go
# # -----------------------
# RUN wget https://go.dev/dl/go1.24.3.linux-${BUILD_ARCH}.tar.gz && \
#     tar -C /usr/local -xzf go1.24.3.linux-${BUILD_ARCH}.tar.gz && \
#     rm go1.24.3.linux-${BUILD_ARCH}.tar.gz

# # -----------------------
# # Install Node via NVM
# # -----------------------
# RUN curl --silent -o- https://raw.githubusercontent.com/creationix/nvm/v0.31.2/install.sh | bash && \
#     . "$NVM_DIR/nvm.sh" && \
#     nvm install $NODE_VERSION && \
#     nvm use $NODE_VERSION && \
#     nvm alias default $NODE_VERSION

# # -----------------------
# # Build Mono from source (required for AL2023 ARM)
# # -----------------------
# RUN dnf install -y libtool automake autoconf gettext make gcc gcc-c++ && git clone --depth 1 --branch main https://github.com/mono/mono.git && \
#     cd mono && \
#     ./autogen.sh --prefix=/usr/local && \
#     make -j$(nproc) && \
#     make install && \
#     cd .. && rm -rf mono
# RUN mono --version

# # -----------------------
# # Python dependencies
# # -----------------------
# RUN pip3 install --no-cache-dir numpy

# # -----------------------
# # Kotlin via SDKMAN
# # -----------------------
# RUN curl -s "https://get.sdkman.io" | bash && \
#     bash -c "source $SDKMAN_DIR/bin/sdkman-init.sh && sdk install kotlin"

# # -----------------------
# # Global npm tools
# # -----------------------
# RUN npm install -g typescript ts-node

# # -----------------------
# # Go setup
# # -----------------------
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# COPY . .

# RUN go build -buildvcs=false -o bootstrap .
# RUN mkdir -p /var/runtime && mv bootstrap /var/runtime/bootstrap && chmod +x /var/runtime/bootstrap

# WORKDIR /var/runtime

# EXPOSE 5000-10000

# CMD ["./bootstrap"]
