FROM alpine:3.6

RUN apk add --no-cache \
        ca-certificates \
        bash \
    && rm -f /var/cache/apk/*

ARG VERSION
ENV FOO ${VERSION}

COPY bin/streaming-service /usr/local/bin/streaming-service

CMD ["/usr/local/bin/streaming-service"]
