FROM golang:alpine AS builder
WORKDIR /app
ADD main.go .
RUN go build -o conductor main.go

FROM ghcr.io/kiwix/kiwix-tools:latest

#ADD entrypoint.sh /entrypoint.sh
#ENTRYPOINT ["/entrypoint.sh"]

COPY --from=builder /app/conductor /usr/bin/conductor

WORKDIR /data

ENV KIWIX_SERVICES=wikipedia
ENV KIWIX_LANGUAGES=en

CMD ["conductor"]
