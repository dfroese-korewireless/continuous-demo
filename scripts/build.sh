#!/bin/bash

echo "Setup..."
mkdir /out /artifacts
cd /go/src/github.com/dfroese-korewireless/continuous-demo

echo "Building..."
go get -d -v ./...
go build -o /out/demo .

echo "Extracting..."
cp -r html/ /out/html
cd /out
tar -zcvf /artifacts/app.tar.gz *