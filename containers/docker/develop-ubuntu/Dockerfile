FROM ubuntu:xenial

RUN \
  apt-get update && apt-get upgrade -q -y && \
  apt-get install -y --no-install-recommends golang git make gcc libc-dev ca-certificates && \
  git clone --depth 1 https://github.com/wtc/go-wtc && \
  (cd go-wtc && make gwtc) && \
  cp go-wtc/build/bin/gwtc /gwtc && \
  apt-get remove -y golang git make gcc libc-dev && apt autoremove -y && apt-get clean && \
  rm -rf /go-wtc

EXPOSE 8545
EXPOSE 10101

ENTRYPOINT ["/gwtc"]
