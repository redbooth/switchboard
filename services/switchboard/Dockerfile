FROM alpine:3.6

RUN set -x \
    && apk add --no-cache ca-certificates

COPY root /
ENTRYPOINT [ "/usr/bin/switchboard" ]
CMD [ ]
