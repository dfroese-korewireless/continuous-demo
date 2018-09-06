FROM alpine:3.8

ADD app.tar.gz /home

# WORKDIR /go/src/github.com/dfroese-korewireless/continuous-demo
# COPY . .

# ADD html/index.html /html/index.html

# RUN go get -d -v ./...
# RUN go build -o /demo .

RUN ls /home

ENTRYPOINT ["/home/demo"]