FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/bastion-server .
COPY packages/db/migrations ./packages/db/migrations

EXPOSE 8080

CMD ["./bastion-server"]
