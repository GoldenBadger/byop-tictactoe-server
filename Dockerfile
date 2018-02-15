# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && dep ensure && go build -o tictactoe

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/tictactoe /app
ENTRYPOINT ./tictactoe
