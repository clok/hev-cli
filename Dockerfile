FROM alpine:3.24.2

COPY hev /usr/local/bin/hev
RUN chmod +x /usr/local/bin/hev

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/hev" ]