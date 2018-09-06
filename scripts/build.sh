#!/bin/bash

mkdir /out /artifacts

go get -d -v ./...
go build -o /out/demo .
tar -C /out -zcvf /artifacts/app.tar.gz *