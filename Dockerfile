# Start by building the application.
FROM golang:1.13-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

#RUN go get -d -v ./...

RUN go build cmd/treactor/main.go

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/src/app/main /treactor
CMD ["/treactor"]





#FROM alpine
#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
#WORKDIR /usr/local/treactor/
#COPY main treactor
#COPY elements.yaml elements.yaml
#ENTRYPOINT ["./treactor"]
