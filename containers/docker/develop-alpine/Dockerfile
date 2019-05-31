FROM alpine:3.5

RUN \
  apk add --update go git make gcc musl-dev linux-headers ca-certificates && \
  git clone --depth 1 https://github.com/wtc/go-wtc && \
  (cd go-wtc && make gwtc) && \
  cp go-wtc/build/bin/gwtc /gwtc && \
  apk del go git make gcc musl-dev linux-headers && \
  rm -rf /go-wtc && rm -rf /var/cache/apk/*

EXPOSE 8545
EXPOSE 10101

ENTRYPOINT ["/gwtc"]
