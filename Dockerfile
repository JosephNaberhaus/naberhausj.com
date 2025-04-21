# A simple image that contains all of the tools needed to build and serve this website.

FROM alpine:latest

RUN apk add git go imagemagick imagemagick-jpeg imagemagick-webp inotify-tools make python3

# Hack to fix dubious ownership error: https://github.com/actions/checkout/issues/1048
RUN git config --global --add safe.directory '*'

WORKDIR /workdir
