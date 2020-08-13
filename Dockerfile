FROM golang:1.11-buster as debugger

WORKDIR /go/src/app
RUN go get -u cloud.google.com/go/cmd/go-cloud-debug-agent


FROM golang:1.11-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN CGO_ENABLED=0 GO111MODULE=on go build -gcflags=all='-N -l' cmd/treactor/main.go

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=debugger /go/bin/go-cloud-debug-agent /go-cloud-debug-agent
COPY --from=build /go/src/app/main /treactor
COPY --from=build /go/src/app/elements.yaml elements.yaml
ADD source-context.json /

#CMD ["/treactor"]
CMD ["/go-cloud-debug-agent", "-sourcecontext=/source-context.json", "-appmodule=main", "-appversion=7", "--", "/treactor"]





#FROM alpine
#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
#WORKDIR /usr/local/treactor/
#COPY main treactor
#COPY elements.yaml elements.yaml
#ENTRYPOINT ["./treactor"]
