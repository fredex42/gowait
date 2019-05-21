FROM alpine:latest

COPY staging/watcher /usr/local/bin
COPY staging/sampleconfig.yaml /etc/gowait/sampleconfig.yaml

ENTRYPOINT /usr/local/bin/watcher -path ${WATCHPATH}
