FROM golang:1.17 AS builder
ENV CGO_ENABLED 0
ADD . /app
WORKDIR /app
RUN go build -ldflags "-s -w" -v -o http-code-checker-discord-hook main.go

FROM alpine:3
RUN apk update && \
    apk add openssl && \
    rm -rf /var/cache/apk/* \
    && mkdir /app

WORKDIR /app

ADD Dockerfile /Dockerfile

COPY --from=builder /app/http-code-checker-discord-hook /app/http-code-checker-discord-hook

RUN chown nobody /app/http-code-checker-discord-hook \
    && chmod 500 /app/http-code-checker-discord-hook

USER nobody

ENTRYPOINT ["/app/http-code-checker-discord-hook"]
