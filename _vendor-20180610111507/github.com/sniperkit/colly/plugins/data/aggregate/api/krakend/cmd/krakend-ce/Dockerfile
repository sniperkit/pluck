FROM alpine:3.7

LABEL maintainer="dortiz@devops.faith"

RUN apk add --no-cache ca-certificates
ADD krakend-alpine /usr/bin/krakend

VOLUME [ "/etc/krakend" ]

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8080
