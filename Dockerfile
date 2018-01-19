FROM golang:alpine AS build-env
ADD ./src/iot-eyecatcher-broker /gopath/src/iot-eyecatcher-broker
ENV GOPATH=/
RUN apk update && \
    apk upgrade && \
    apk add git && \
    cd /gopath/src/iot-eyecatcher-broker && \
    go get ./... && \
    go get -t ./... && \
    go generate && \
    go test ./... && \
    go build -ldflags "-linkmode external -extldflags -static" -o /iot-eyecatcher-broker .

FROM scratch
USER 100:100
EXPOSE 8080
ENV WS_LISTEN=:8080
COPY --from=build-env /iot-eyecatcher-broker /iot-eyecatcher-broker
ENTRYPOINT ["/iot-eyecatcher-broker"]
