FROM alpine:3.14.0

COPY hev /usr/local/bin/hev
RUN chmod +x /usr/local/bin/hev

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/hev" ]