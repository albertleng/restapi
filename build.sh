#!/bin/bash

echo "Go build"
cd "$GOPATH"/src/github.com/albertleng/restapi || exit
go build -o restapi
echo "Go build done"
./restapi
