FROM alpine:3.8

RUN apk add --no-cache ca-certificates

COPY cryobot_linux_amd64 .
RUN chmod +x cryobot_linux_amd64

CMD ["./cryobot_linux_amd64"]
