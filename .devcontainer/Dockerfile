FROM golang:1.20.11

RUN apt-get update \
    && apt-get install -y curl git \
    && apt-get install -y --no-install-recommends vim \
    && apt-get install -y libonig-dev libzip-dev unzip less \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
# alpine
# FROM golang:1.20.2-alpine3.17
#
# RUN apk update \
#     # && apk add --no-cache alpine-sdk build-base \
#     && apk add libstdc++ \
#     && apk add --no-cache bash vim curl git unzip less \
#     && rm -rf /var/cache/apk/*

ENV GO111MODULE on
WORKDIR /go/src

# aws-cli
RUN curl https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip -o awscliv2.zip \
    && unzip awscliv2.zip \
    && ./aws/install

# install go tools（自動補完等に必要なツールをコンテナにインストール）
# https://qiita.com/Gatsby/items/106236ac87dbfdddcb65
# インポート可能なパッケージのリストを表示
RUN go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest
# JSON表現を抽出するためのシンプルなユーティリティ
RUN go install github.com/ramya-rao-a/go-outline@master
# Linter
RUN go install honnef.co/go/tools/cmd/staticcheck@2023.1.2
# パッケージシンボルの JSON 表記を抽出するためのユーティリティ
RUN go install github.com/acroca/go-symbols@master
# 構造体フィールドのタグを更新、追加、削除する
RUN go install github.com/fatih/gomodifytags@latest
# implスタブメソッド生成
RUN go install github.com/josharian/impl@latest
# web interface
RUN go install github.com/haya14busa/goplay/cmd/goplay@master
# デバッガ
RUN go install github.com/go-delve/delve/cmd/dlv@master
# LanguageServer(補完等)
RUN go install golang.org/x/tools/gopls@latest
