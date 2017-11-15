FROM golang:1.8

RUN cat /etc/os-release

RUN apt-get update

RUN wget http://http.us.debian.org/debian/pool/non-free/x/xml2rfc/xml2rfc_2.4.8-1_all.deb

#Install binaries from debian needed for golang, python, xml2rfc, and mmark
RUN apt-get -y install python python-lxml

RUN dpkg -i xml2rfc_2.4.8-1_all.deb
RUN apt-get install -f

RUN go get github.com/BurntSushi/toml
RUN go get github.com/miekg/mmark/...
RUN go get -u github.com/jteeuwen/go-bindata/...

RUN mkdir -p /go/src/github.com/koki/short 
WORKDIR /go/src/github.com/koki/short
COPY .* /go/src/github.com/koki/short/

ENTRYPOINT ["scripts/ci.sh"]
