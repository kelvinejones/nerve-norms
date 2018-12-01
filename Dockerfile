FROM golang:1.11
ADD ./cmd /go/src/bellstone.ca/jitter/cmd
ADD ./lib /go/src/bellstone.ca/jitter/lib
ADD ./res /go/src/bellstone.ca/jitter/res
RUN go test bellstone.ca/jitter/cmd/... bellstone.ca/jitter/lib/...
RUN go install bellstone.ca/jitter/cmd/...
ADD ./resources /go/src/bellstone.ca/jitter/resources
WORKDIR /go/src/bellstone.ca/jitter/cmd/jitter
