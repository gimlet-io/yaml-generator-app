FROM golang:alpine as builder 

RUN apk add --no-cache git curl make ca-certificates gcc libtool musl-dev

RUN mkdir /app

ADD cmd /app
ADD go.mod /app
ADD go.sum /app
ADD Makefile /app

WORKDIR /app

RUN make build-yaml-generator-app

FROM alpine

WORKDIR /

COPY --from=builder /app/build/yaml-generator-app .

EXPOSE 9000

ENTRYPOINT ["./yaml-generator-app"]
