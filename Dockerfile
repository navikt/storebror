FROM golang:alpine AS build
ENV D=/go/src/github.com/navikt/storebror
ADD . $D
WORKDIR $D
RUN go test ./...
RUN go build -o /storebror

FROM alpine
COPY --from=build /storebror /storebror
CMD ["/storebror"]
