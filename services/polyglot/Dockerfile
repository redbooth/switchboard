FROM alpine:3.6

RUN set -x \
    && apk add --no-cache ca-certificates python3-dev \
    && pip3 install --upgrade google-cloud-translate

COPY root /
ENTRYPOINT [ "/usr/bin/polyglot" ]
CMD [ ]
