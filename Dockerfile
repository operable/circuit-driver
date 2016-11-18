FROM alpine:3.4

MAINTAINER Kevin Smith <kevin@operable.io>

RUN mkdir -p /operable/circuit/bin

COPY _build/circuit-driver /operable/circuit/bin

RUN chmod +x /operable/circuit/bin && \
    apk -U add binutils && \
    strip --strip-all --verbose /operable/circuit/bin/circuit-driver && \
    apk del binutils

VOLUME /operable/circuit
