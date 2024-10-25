FROM golang:1.22.7 as builder

WORKDIR /fskick-api

COPY . .

RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs && \
    npm install --save-dev tailwindcss && \
    go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /fskick-api/cmd/server

RUN set -x && \
    go generate && \
    go get -d -v . && \
    CGO_ENABLED=1 GOOS=linux go build \
        -a \
        -ldflags="-X github.com/spie/fskick/internal/templates.version=$(git describe)" \
        -o server .

WORKDIR /fskick-api/cmd/cli

RUN set -x && \
    go generate && \
    go get -d -v . && \
    CGO_ENABLED=1 GOOS=linux go build -o fskick .

FROM debian:latest

WORKDIR /app

COPY --from=builder /fskick-api/cmd/server/server .
COPY --from=builder /fskick-api/cmd/server/.env .
COPY --from=builder /fskick-api/cmd/cli/fskick .

CMD ["./server"]
