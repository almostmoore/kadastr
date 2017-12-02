FROM golang:1.9.0

COPY . $GOPATH/src/github.com/iamsalnikov/kadastr

RUN go get -u github.com/govend/govend &&\
    cd $GOPATH/src/github.com/iamsalnikov/kadastr &&\
    govend -v &&\
    go build &&\
    go install github.com/iamsalnikov/kadastr