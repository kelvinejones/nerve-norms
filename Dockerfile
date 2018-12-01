FROM golang:1.11

ADD ./vendor /go/src/jitter/vendor
ADD ./cmd /go/src/jitter/cmd
ADD ./lib /go/src/jitter/lib
ADD ./res /go/src/jitter/res

RUN go test jitter/cmd/... jitter/lib/...
RUN go install jitter/cmd/...
ADD ./resources /go/src/jitter/resources

WORKDIR /go/src/jitter
CMD /go/bin/jitter
