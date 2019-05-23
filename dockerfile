FROM alpine:3.8

ENV CONTAINER_NAME "container_name"

ADD app.tar.gz /

ENTRYPOINT ["/demo"]
