FROM alpine:3.8

ENV CONTAINER_NAME "container_name"

ADD app.tar.gz /

# WORKDIR /go/src/github.com/dfroese-korewireless/continuous-demo
# COPY . .

# ADD html/index.html /html/index.html

# RUN go get -d -v ./...
# RUN go build -o /demo .

ENTRYPOINT ["/demo"]