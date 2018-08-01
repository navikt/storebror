FROM golang:alpine AS build
RUN apk add --no-cache git && \
    wget https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 -O /usr/bin/dep && \
    chmod 755 /usr/bin/dep
ENV D=/go/src/github.com/navikt/storebror
ADD . $D
WORKDIR $D
RUN dep ensure
RUN go test ./...
RUN go build -o /storebror

FROM alpine
RUN apk add --no-cache git
COPY --from=build /storebror /storebror
ENTRYPOINT ["/storebror"]
