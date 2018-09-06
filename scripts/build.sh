#!/bin/bash

echo "Setup..."
mkdir /out /artifacts
cd /go/src/github.com/dfroese-korewireless/continuous-demo

echo "Building..."
go get -d -v ./...
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o /out/demo .
# go build -o /out/demo .

echo "Extracting..."
cp -r html/ /out/html
cp appsettings.json /out
cd /out
tar -zcvf /artifacts/app.tar.gz *