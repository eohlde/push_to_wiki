FROM golang:latest

WORKDIR /go
COPY wiki_push.go .

RUN go build wiki_push.go

RUN mkdir /code

RUN mv ./wiki_push /code/

WORKDIR /code

ENTRYPOINT [ "/code/wiki_push" ]