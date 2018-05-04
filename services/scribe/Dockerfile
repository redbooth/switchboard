FROM alpine:3.6

RUN set -x \
    && apk add --no-cache --virtual .build-deps g++ linux-headers make musl-dev \
    && apk add --no-cache ca-certificates libgcc python3-dev \
    && pip3 install --upgrade google-cloud-speech \
    && apk del .build-deps

COPY root /
ENTRYPOINT [ "/usr/bin/scribe" ]
CMD [ ]
