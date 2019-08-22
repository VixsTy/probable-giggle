# Arguments
ARG BUILD_DATE
ARG VERSION
ARG VCS_REF

## -------------------------------------------------------------------------------------------------

FROM golang:1.12 as tools
RUN set -eux; \
    apt-get update -y && \
    apt-get install -y apt-utils upx
WORKDIR /src
# Force go modules
ENV GO111MODULE=on
COPY tools tools
RUN cd tools && go run github.com/magefile/mage

## -------------------------------------------------------------------------------------------------

FROM tools AS deps
COPY grimoire/ magefile.go go.mod ./
COPY grimoire grimoire
RUN set -eux; \
    go run github.com/magefile/mage go:deps

FROM deps as source
COPY . .

## -------------------------------------------------------------------------------------------------

FROM source AS build
RUN go run github.com/magefile/mage
# Compress binaries
RUN set -eux; \
    upx -9 bin/* && \
    chmod +x bin/*

## -------------------------------------------------------------------------------------------------

FROM gcr.io/distroless/base:latest
WORKDIR /
# Arguments
ARG BUILD_DATE
ARG VERSION
ARG VCS_REF
# Metadata
LABEL \
    org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="ProbableGiggle" \
    org.label-schema.description="Typical golang project" \
    org.label-schema.url="https://github.com/VixsTy/probable-giggle" \
    org.label-schema.vcs-url="https://github.com/VixsTy/probable-giggle.git" \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vendor="Kevin LARQUEMIN" \
    org.label-schema.version=$VERSION \
    org.label-schema.schema-version="1.0" \
    org.zenithar.licence="MIT"
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/bin/* /app/bin
ENTRYPOINT [ "/app/bin" ]
CMD ["--help"]



