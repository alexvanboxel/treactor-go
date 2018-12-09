FROM alpine
COPY reactor reactor
COPY elements.yaml elements.yaml
CMD ./reactor
