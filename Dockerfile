FROM alpine:3.17.1

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN addgroup -S yaml-generator && adduser -S yaml-generator -G yaml-generator

COPY --chown=yaml-generator:yaml-generator bin/yaml-generator-app-linux-x86_64 /bin/yaml-generator-app

USER yaml-generator

EXPOSE 9000
CMD ["/bin/yaml-generator-app"]
