FROM alpine:3.13.4

COPY hev /usr/local/bin/hev
RUN chmod +x /usr/local/bin/hev

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/hev" ]