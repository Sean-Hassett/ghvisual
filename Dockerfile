FROM golang:latest
RUN mkdir $GOPATH/src/ghvisual
ADD . $GOPATH/src/ghvisual/
WORKDIR $GOPATH/src/ghvisual

RUN go get -u github.com/kardianos/govendor
RUN govendor sync

RUN cd ghvisual && go build -o main .
RUN chmod +x ghvisual/main
CMD ["./ghvisual/main"]
