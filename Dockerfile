FROM golang:latest as builder

WORKDIR /fskick-api

COPY . .

WORKDIR /fskick-api/cmd/api

RUN set -x && go get -d -v . && \
    CGOENABLED=0 GOOS=linux go build -a -o api .

WORKDIR /fskick-api/cmd/cli

RUN set -x && go get -d -v . && \
    CGOENABLED=0 GOOS=linux go build -o fskick .

WORKDIR /fskick-api/cmd/migrations

RUN set -x && go get -d -v . && \
    CGOENABLED=0 GOOS=linux go build -a -o migrations .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /fskick-api/cmd/api/api .
COPY --from=builder /fskick-api/cmd/cli/fskick .
COPY --from=builder /fskick-api/cmd/migrations/migrations .

CMD ["./api"]