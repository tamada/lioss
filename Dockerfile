FROM alpine:3.10.1
ARG version="0.9.0"
LABEL maintainer="Haruaki Tamada" \
    lioss-version="${version}" \
    description="License Identifier for Open Source Software."

RUN    adduser -D lioss \
    && apk --no-cache add --virtual .builddeps curl tar \
    && curl -s -L -O https://github.com/tamada/lioss/releases/download/v${version}/lioss-${version}_linux_amd64.tar.gz \
    && tar xfz lioss-${version}_linux_amd64.tar.gz         \
    && mv lioss-${version} /opt                            \
    && ln -s /opt/lioss-${version} /opt/lioss              \
    && ln -s /opt/lioss /usr/local/share/lioss             \
    && ln -s /opt/lioss/lioss /usr/local/bin/lioss         \
    && ln -s /opt/lioss/mkliossdb /usr/local/bin/mkliossdb \
    && rm lioss-${version}_linux_amd64.tar.gz              \
    && apk --del --purge .builddeps

ENV HOME="/home/lioss" \
    LIOSS_DBPATH="/opt/lioss/liossdb.json"

WORKDIR /home/lioss
USER    lioss

ENTRYPOINT [ "lioss" ]
