FROM janeczku/alpine-kubernetes:3.3
MAINTAINER pjpires@gmail.com

EXPOSE 8080

RUN apk add --update ca-certificates

COPY bin/nats-sniffer-linux-amd64 /nats-sniffer
COPY run.sh /run.sh

CMD [ "/run.sh" ]
