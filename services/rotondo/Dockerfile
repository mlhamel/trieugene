FROM golang:1.17.6-alpine3.15 AS base

WORKDIR /go/src/rotondo
COPY . .
ADD go.mod .
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["rotondo", "dev"]