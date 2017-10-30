FROM golang:1.9

WORKDIR /go/src/app
ADD . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
