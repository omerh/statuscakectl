FROM golang:1.13.7 AS Builder

COPY . /go/src
WORKDIR /go/src
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags '-w -extldflags "-static"'

FROM alpine:3.9.5
COPY --from=Builder /go/src/statuscakectl /usr/bin/

CMD ["statuscakectl"]

