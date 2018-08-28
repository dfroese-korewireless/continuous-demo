FROM golang:1.10

WORKDIR /go/src/github.com/dfroese-korewireless/continuous-demo
COPY . .

ADD html/index.html /html/index.html

RUN go get -d -v ./...
RUN go build -ldflags "-X main.Version=${env.BuildNumber}" -o /demo .

ENTRYPOINT ["/demo"]