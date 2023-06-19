#!/bin/sh
while inotifywait -q -e modify -r src; do
  make build
done