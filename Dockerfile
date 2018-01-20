FROM golang AS build-env
ADD . /gopath/src/github.com/tarent/iot-eyecatcher-broker
RUN cd /gopath/src/github.com/tarent/iot-eyecatcher-broker && \
    export GOPATH=/gopath && \
    export PATH=$PATH:/gopath/bin && \
    go get ./... && \
    go get -t ./... && \
    go generate && \
    go test ./... && \
    go build -ldflags "-linkmode external -extldflags -static" -o /iot-eyecatcher-broker .

FROM scratch
USER 100:100
EXPOSE 8080
ENV WS_LISTEN=:8080
ENV USERNAME=tarent
ENV PASSWORD=changeme
COPY --from=build-env /iot-eyecatcher-broker /iot-eyecatcher-broker
ENTRYPOINT ["/iot-eyecatcher-broker"]
