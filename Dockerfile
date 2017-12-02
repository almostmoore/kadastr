FROM alpine:3.5

RUN apk add --no-cache ca-certificates && update-ca-certificates

ADD kadastr /
RUN chmod +x /kadastr