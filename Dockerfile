FROM alpine:edge
MAINTAINER Tobias Gesellchen <tobias@gesellix.de> (@gesellix)

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/gesellix/reverse-proxy

ENV ADD_PACKAGES go@community
ENV DEL_PACKAGES go

COPY . $APPPATH

RUN echo '@community http://dl-4.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories \
    && apk upgrade --update --available \
    && apk add $ADD_PACKAGES \
    && cd $APPPATH && go get -d && go build -o /bin/reverse-proxy \
    && apk del --purge $DEL_PACKAGES \
    && rm -rf /var/cache/apk/* && rm -rf $GOPATH

# enforce go to prefer /etc/hosts
ENV GODEBUG netdns=go+1

ENTRYPOINT [ "/bin/reverse-proxy" ]
