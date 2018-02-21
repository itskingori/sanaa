FROM ubuntu:17.10

RUN apt-get update -y && \
    apt-get install --no-install-recommends -y \
      # build process
      ca-certificates wget \
      # wkhtml
      wkhtmltopdf="0.12.3.2-3" && \
    rm -f -R /tmp/* /var/{cache,log,tmp} /var/lib/{apt,dpkg,cache,log}

# sanaa
ARG SANAA_VERSION
ARG SANAA_PACKAGE="sanaa-${SANAA_VERSION}-linux-amd64"
RUN wget --progress="dot:mega" "https://github.com/itskingori/sanaa/releases/download/${SANAA_VERSION}/${SANAA_PACKAGE}.tar.gz" && \
    wget --progress="dot:mega" "https://github.com/itskingori/sanaa/releases/download/${SANAA_VERSION}/${SANAA_PACKAGE}-shasum-256.txt" && \
    sha256sum -c "${SANAA_PACKAGE}-shasum-256.txt" && \
    tar --no-same-owner -xzf "${SANAA_PACKAGE}.tar.gz" && \
    mv "/${SANAA_PACKAGE}" "/usr/local/bin/sanaa" && \
    chmod +x "/usr/local/bin/sanaa" && \
    rm -f ${SANAA_PACKAGE}*

ENTRYPOINT ["sanaa"]
