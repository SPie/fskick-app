FROM golang:latest as builder

WORKDIR /fskick-api

COPY . .

WORKDIR /fskick-api/cmd/api

RUN set -x && go get -d -v . && \
    CGOENABLED=0 GOOS=linux go build -a -o api .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /fskick-api/cmd/api/api .

CMD ["./api"]