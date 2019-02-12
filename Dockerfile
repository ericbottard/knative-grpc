FROM golang:1.11.5

WORKDIR /go/src/app
COPY . .
ENV GO111MODULE=on
RUN go build -o run-server ./server

CMD ["./run-server"]