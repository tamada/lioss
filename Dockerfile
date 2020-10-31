FROM alpine:3.10.1

ARG version="1.0.0"

LABEL maintainer="Haruaki Tamada" \
    lioss-version="${version}" \
    description="License Identifier for Open Source Software."

RUN    adduser -D lioss \
    && apk --no-cache add --virtual .builddeps curl tar \
    && curl -s -L -O https://github.com/tamada/lioss/releases/download/v${version}/lioss-${version}_linux_amd64.tar.gz \
    && tar xfz lioss-${version}_linux_amd64.tar.gz             \
    && mv lioss-${version} /opt                                \
    && ln -s /opt/lioss-${version} /opt/lioss                  \
    && ln -s /opt/lioss /usr/local/share/lioss                 \
    && ln -s /opt/lioss/bin/lioss /usr/local/bin/lioss         \
    && ln -s /opt/lioss/bin/mkliossdb /usr/local/bin/mkliossdb \
    && rm lioss-${version}_linux_amd64.tar.gz                  \
    && rm -rf /opt/lioss/{README.md,LICENSE,completions}       \
    && apk del --purge .builddeps

ENV HOME="/home/lioss"  \
    LIOSS_DBPATH="/opt/lioss/data"

WORKDIR /home/lioss
USER    lioss

ENTRYPOINT [ "/opt/lioss/bin/lioss" ]
