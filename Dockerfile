# syntax=docker/dockerfile:1

## Build Stage
FROM golang:1.22.3-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go-app ./cmd/main.go

## Runtime Stage
FROM ubuntu:20.04

# ENV ANSIBLE_VERSION 2.13.13
# ENV JINJA2_VERSION 3.1.4

# 환경 변수 설정 및 타임존 설정
ENV TZ=Asia/Seoul
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# # 필수 패키지 설치 및 청소
# RUN apt-get update && \
#     apt-get install -y \
#     gcc \
#     python3 \
#     python3-pip \
#     curl \
#     net-tools \
#     vim \
#     git \
#     openssh-server \
#     sshpass && \
#     apt-get clean && \
#     rm -rf /var/lib/apt/lists/*

# # Locale 설치
# RUN apt-get update && apt-get install -y locales && \
#     locale-gen en_US.UTF-8 && \
#     update-locale LANG=en_US.UTF-8 LC_ALL=C.UTF-8

# # 환경 변수 설정
# ENV LANG=en_US.UTF-8
# ENV LC_ALL=C.UTF-8

# Install necessary packages including locales-all
# RUN apt-get clean && \
#     rm -rf /var/lib/apt/lists/* && \
#     apt-get update --fix-missing && \
#     apt-get install -y --no-install-recommends \
#     gcc python3 python3-pip curl net-tools vim git openssh-server sshpass locales-all && \
#     apt-get clean && rm -rf /var/lib/apt/lists/*
# RUN apt-get update 
RUN apt-get clean
RUN apt-get update 
RUN apt-get upgrade -y
RUN apt-get install -y \
    gcc \
    python3 \
    python3-pip \
    curl \
    net-tools \
    vim \
    git \
    openssh-server \
    sshpass \
    locales-all
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set locale environment variables
ENV LANG=ko_KR.UTF-8
ENV LANGUAGE=ko_KR.UTF-8
ENV LC_ALL=ko_KR.UTF-8

# Python 패키지 설치
RUN pip3 install --upgrade pip && \
    pip3 install ansible jinja2

# Ansible Galaxy Collection 설치
RUN ansible-galaxy collection install community.general ansible.posix

# 작업 디렉토리 설정 및 파일 복사
WORKDIR /root
COPY --from=build /go-app /go-app
COPY config.prod.json /config.prod.json
COPY config.dev.json /config.dev.json
COPY ansible ansible

EXPOSE 8080

ENTRYPOINT ["/go-app"]
