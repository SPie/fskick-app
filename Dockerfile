FROM golang:1.22.7 as builder

WORKDIR /fskick-api

COPY . .

RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs && \
    npm install --save-dev tailwindcss && \
    go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /fskick-api/cmd/server

RUN set -x && \
    go get -d -v . && \
    go generate && \
    CGOENABLED=0 GOOS=linux go build -a -o server .

WORKDIR /fskick-api/cmd/cli

RUN set -x && \
    go get -d -v . && \
    go generate && \
    CGOENABLED=0 GOOS=linux go build -o fskick .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /fskick-api/cmd/server/server .
COPY --from=builder /fskick-api/cmd/server/.env .
COPY --from=builder /fskick-api/cmd/cli/fskick .

CMD ["./server"]
