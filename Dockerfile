FROM alpine:3.17.1

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

ADD bin/yaml-generator-app-linux-x86_64 /bin/yaml-generator-app

EXPOSE 9000
CMD ["/bin/yaml-generator-app"]