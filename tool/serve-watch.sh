#!/bin/sh

# Do an initial build.
make internal-cache-build

# And do another build every time a file changes.
while inotifywait -q -e modify -r src; do
  make internal-cache-build
done
