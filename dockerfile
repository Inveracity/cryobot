FROM alpine:3.8

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tini

COPY config.yml .
COPY cryobot_linux_amd64 .
RUN chmod +x cryobot_linux_amd64

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["./cryobot_linux_amd64"]
