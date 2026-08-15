FROM golang:1.26.6

RUN apt-get update \
    && apt-get install -y curl git \
    && apt-get install -y --no-install-recommends vim \
    && apt-get install -y libonig-dev libzip-dev unzip less \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV GOTOOLCHAIN=auto
WORKDIR /go/src

# aws-cli
RUN curl https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip -o awscliv2.zip \
    && unzip awscliv2.zip \
    && ./aws/install

RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install github.com/cweill/gotests/gotests@latest
