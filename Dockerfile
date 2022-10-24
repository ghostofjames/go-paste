FROM golang:latest

WORKDIR $GOPATH/src/go-paste

COPY . . 

RUN go build -o main .

EXPOSE 8000

CMD ./main