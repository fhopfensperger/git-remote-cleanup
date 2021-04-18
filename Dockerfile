FROM golang:1.16 as goBuilder

USER root
WORKDIR /work
COPY . .
ARG BUILD_VERSION="0.0.0"
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=$BUILD_VERSION" -o git-remote-cleanup .
RUN find .

FROM alpine:3.13.5

LABEL maintainer="Florian Hopfensperger <f.hopfensperger@gmail.com>"

RUN apk add --update ca-certificates \
    && apk add --update -t wget git openssl \
    && rm /var/cache/apk/* \
    && adduser -G root -u 1000 -D -S kuser

USER 1000
WORKDIR /app

COPY --chown=1000:0 --from=goBuilder /work/git-remote-cleanup .

ENTRYPOINT ["./git-remote-cleanup"]
CMD ["--help"]