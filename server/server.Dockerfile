FROM golang

ADD . /go/src/github.com/koki/short/server

RUN go install github.com/koki/short/server

ENTRYPOINT ["/go/bin/server"]
CMD []
