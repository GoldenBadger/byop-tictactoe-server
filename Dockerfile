FROM golang:latest as builder
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/goldenbadger/byop-tictactoe-server
COPY . .
RUN dep ensure
RUN go build .

FROM alpine:latest
WORKDIR /root
COPY --from=builder /go/src/github.com/goldenbadger/byop-tictactoe-server/byop-tictactoe-server .
CMD ["./byop-tictactoe-server"]
