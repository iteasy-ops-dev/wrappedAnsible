# syntax=docker/dockerfile:1

## Build
FROM golang:1.22.3-alpine3.20 AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /go-app ./cmd/main.go

FROM ubuntu:20.04

# ENV ANSIBLE_VERSION 2.13.13
# ENV JINJA2_VERSION 3.1.4

ENV TZ=Asia/Seoul
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update; \
    apt-get install -y gcc python3; \
    apt-get install -y python3-pip; \
    apt-get install -y curl net-tools vim git; \
    apt-get install -y openssh-server sshpass; \
    apt-get clean all
RUN pip3 install --upgrade pip
RUN pip3 install ansible 
RUN pip3 install jinja2

RUN ansible-galaxy collection install community.general

WORKDIR /root
COPY --from=build /go-app /go-app
COPY config.json /config.json
COPY ansible ansible
COPY server.crt /server.crt
COPY server.csr /server.csr
COPY server.key /server.key
EXPOSE 8080

ENTRYPOINT ["/go-app"]