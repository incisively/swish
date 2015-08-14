FROM golang:1.4.2
MAINTAINER github.com/skinnyfit

RUN apt-get update && apt-get install -qy socat net-tools

# test & install
COPY . /go/src/github.com/incisively/swish/
RUN cd /go/src/github.com/incisively/swish && go build && ./test.sh
RUN go install github.com/incisively/swish

EXPOSE 8999
CMD ["swish"]
