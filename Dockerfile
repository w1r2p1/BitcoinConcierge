FROM golang:latest

MAINTAINER Anh Nguyen <anh.nguyen221193@gmail.com>

RUN curl https://glide.sh/get | sh

RUN mkdir /go/src/BitcoinBot
WORKDIR /go/src/BitcoinBot
COPY . .

RUN glide install
RUN go build -o main .
CMD ["/go/src/BitcoinBot/main"]
EXPOSE 8080

