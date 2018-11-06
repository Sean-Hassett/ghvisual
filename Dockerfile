FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN cd ghvisual && go build -o main . && cd ..
CMD ["/app/main"]

