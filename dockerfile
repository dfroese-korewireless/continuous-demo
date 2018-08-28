FROM golang:1.10 as build

WORKDIR /go/src/github.com/dfroese-korewireless/continuous-demo
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=%system.BuildNumber%" -a -installsuffix cgo -o /demo .

FROM scratch

COPY --from=build /demo /
ADD html/index.html /html/index.html

ENTRYPOINT [ "/demo" ]