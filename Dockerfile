FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /usr/local/reactor/
COPY gopath/bin/reactor reactor
COPY elements.yaml elements.yaml
ENTRYPOINT ["./reactor"]
