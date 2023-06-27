FROM golang:1.20

RUN apt update && apt upgrade -y &&\
    apt install -y git\
    make openssh-client

WORKDIR /go/src/app

COPY . ./
RUN go mod tidy \
    && go mod verify

EXPOSE 2112

RUN go build -o  /main

CMD ["/main"]
