FROM golang:1.14.2-alpine
RUN apk add alpine-sdk

RUN mkdir -p go/simulator
COPY . go/simulator
WORKDIR go/simulator

RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 go build -o simulator -tags musl
CMD ["./simulator"]