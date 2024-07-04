# A simple image that contains all of the tools needed to build and serve this website.

FROM alpine:latest

RUN apk add git go imagemagick imagemagick-jpeg imagemagick-webp inotify-tools make python3

WORKDIR /workdir
