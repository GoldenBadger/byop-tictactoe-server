FROM golang:latest as builder
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/goldenbadger/byop-tictactoe-server
COPY . .
RUN dep ensure
RUN go build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o byop-tictactoe-server main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /go/src/github.com/goldenbadger/byop-tictactoe-server/byop-tictactoe-server .
CMD ["./byop-tictactoe-server"]
