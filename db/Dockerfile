FROM golang:1.11
RUN go get -v bitbucket.org/liamstask/goose/cmd/goose
WORKDIR /
ADD . /db
CMD /db/wrapper.sh
