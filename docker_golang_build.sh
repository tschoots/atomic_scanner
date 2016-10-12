#!/bin/bash

#_TAG=1.6.2
_TAG=1.7.1


echo "remove scanner"
rm -rf atomic_scanner

echo "remove the local image blackducksoftware/atomic"
docker rmi blackducksoftware/atomic

echo "building go executable"
docker run --rm -v $(pwd):/usr/src/myapp -w /usr/src/myapp golang:$_TAG bash -c "go get -u github.com/fsnotify/fsnotify; ./build.sh"

echo "building new image"
docker build -t blackducksoftware/atomic .
