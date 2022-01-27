FROM ubuntu:18.04 AS base
ARG TARGETARCH

RUN mkdir -p /app

COPY url /app

RUN sed -i 's/\(archive\|security\|ports\).ubuntu.com/mirrors.aliyun.com/' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y libaio-dev libaio1 unzip wget curl


CMD ["/app/url"]
