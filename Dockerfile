FROM golang:1.4
MAINTAINER jerbia@gmail.com

RUN go get github.com/ChimeraCoder/anaconda
RUN go get github.com/hashicorp/errwrap

RUN apt-get update
RUN apt-get install -y socat jq

COPY . /go/src/github.com/advanderveer/docksec
WORKDIR /go/src/github.com/advanderveer/docksec
RUN go build -v

CMD /go/src/github.com/advanderveer/docksec/docksec
