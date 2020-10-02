FROM golang:1.14.2-alpine
RUN apk add alpine-sdk

RUN mkdir -p go/datacenter
COPY . go/datacenter
WORKDIR go/datacenter

RUN go get -d -v ./...
RUN go build -o datacenter

EXPOSE 8001
CMD ["./datacenter"]


