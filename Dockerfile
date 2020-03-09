FROM alpine:3.10.1
ARG version="0.9.0"
LABEL maintainer="Haruaki Tamada" \
      lioss-version="${version}" \
      description="License Identifier for Open Source Software."

RUN    adduser -D lioss \
    && apk --no-cache add curl=7.66.0-r0 tar=1.32-r0 \
    && curl -s -L -O https://github.com/tamada/lioss/releases/download/v${version}/lioss-${version}_linux_amd64.tar.gz \
    && tar xfz lioss-${version}_linux_amd64.tar.gz      \
    && mv lioss-${version} /opt                         \
    && ln -s /opt/lioss-${version} /opt/lioss           \
    && ln -s /opt/lioss /usr/local/share/lioss          \
    && rm lioss-${version}_linux_amd64.tar.gz           \
    && ln -s /opt/lioss/lioss /usr/local/bin/lioss      \
    && ln -s /opt/lioss/mkliossdb /usr/local/bin/mkliossdb

ENV HOME="/home/lioss"

WORKDIR /home/lioss
USER    lioss

ENTRYPOINT [ "lioss" ]
