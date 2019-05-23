#!/bin/bash

echo "Setup..."
mkdir -p /out/debug /out/release /artifacts
cd /go/src/github.com/dfroese-korewireless/continuous-demo

echo "Building..."
go get -d -v ./...
go build -o /out/debug/demo-debug .
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o /out/release/demo .

echo "Extracting..."
cp -r html/ /out/release/html
cp -r html/ /out/debug/html
cp appsettings.json /out/release
cp appsettings.json /out/debug
cd /out/release
tar -zcf /go/src/github.com/dfroese-korewireless/continuous-demo/app.tar.gz *
cd /out/debug
tar -zcf /go/src/github.com/dfroese-korewireless/continous-demo/app-debug-$BUILD_NUMBER.tar.gz *
