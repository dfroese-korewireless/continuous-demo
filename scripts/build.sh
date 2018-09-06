#!/bin/bash

mkdir /out /artifacts

cd /go/src/github.com/dfroese-korewireless/continuous-demo
go get -d -v ./...
go build -o /out/demo .
tar -C /out -zcvf /artifacts/app.tar.gz *