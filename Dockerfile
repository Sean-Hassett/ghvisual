FROM golang:latest
RUN mkdir $GOPATH/src/ghvisual
ADD . $GOPATH/src/ghvisual/
WORKDIR $GOPATH/src/ghvisual

RUN go get -u github.com/kardianos/govendor
RUN govendor sync

RUN go build -o main .
RUN chmod +x main
CMD ["/main"]
