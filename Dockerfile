FROM golang:buster as builder

ENV GOPROXY=https://goproxy.cn

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o QN .

FROM debian:buster as runner

WORKDIR /app

COPY --from=builder /app/QN .

ENTRYPOINT ["./QN"]
