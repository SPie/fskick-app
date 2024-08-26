FROM golang:latest as builder

WORKDIR /fskick-api

COPY . .

WORKDIR /fskick-api/cmd/server

RUN set -x && go get -d -v . && \
    CGOENABLED=0 GOOS=linux go build -a -o server .

WORKDIR /fskick-api/cmd/cli

RUN set -x && go get -d -v . && \
    CGOENABLED=0 GOOS=linux go build -o fskick .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /fskick-api/cmd/server/server .
COPY --from=builder /fskick-api/cmd/cli/fskick .

CMD ["./server"]
