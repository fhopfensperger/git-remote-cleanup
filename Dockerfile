FROM golang:1.20.2 as goBuilder

USER root
WORKDIR /work
COPY . .
ARG BUILD_VERSION="0.0.0"
RUN CGO_ENABLED=0 go build -a -ldflags "-X main.version=$BUILD_VERSION" -o git-remote-cleanup .

FROM alpine:3.17.3

LABEL maintainer="Florian Hopfensperger <f.hopfensperger@gmail.com>"

RUN apk add --update wget git openssl openssh ca-certificates \
    && rm /var/cache/apk/* \
    && adduser -G root -u 1000 -D -S kuser

USER 1000
WORKDIR /app

COPY --chown=1000:0 --from=goBuilder /work/git-remote-cleanup .

ENTRYPOINT ["./git-remote-cleanup"]
CMD ["--help"]