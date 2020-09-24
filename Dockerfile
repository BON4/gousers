FROM golang:1.15.1
RUN mkdir /go/user-grpc
WORKDIR /go/user-grpc/cmd/server
ADD . /go/user-grpc
##COPY . /test/

ENV GO111MODULE=on

RUN go mod download
RUN go get -v
## -o создаст exe файл в текущей дериктории
##RUN go build -o main

EXPOSE 8080 80 5432 6379