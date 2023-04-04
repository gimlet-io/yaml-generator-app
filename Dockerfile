FROM --platform=$BUILDPLATFORM alpine:3.17.1

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir /var/lib/yaml-generator-app
WORKDIR /yaml-generator-app

ARG TARGETPLATFORM
ARG BUILDPLATFORM

COPY bin/${TARGETPLATFORM}/yaml-generator-app yaml-generator-app

EXPOSE 9000
CMD ["/yaml-generator-app/yaml-generator-app"]