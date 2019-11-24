FROM golang:1.12-alpine3.9 as build

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update \
    && apk --no-cache add git
COPY src /go/src
WORKDIR /go
ENV GO111MODULE on
ENV GOPROXY "https://goproxy.io"
RUN cd /go/src/api && go build


FROM alpine:3.9

WORKDIR /www
ENV GOPATH /www
COPY --from=build /go/src/api/api .
COPY .env.default /www/.env
ARG COMMIT
ENV COMMIT ${COMMIT}

ENTRYPOINT ["/www/api"]