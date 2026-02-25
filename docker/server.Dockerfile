FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/bastion-server .
COPY packages/core/db/migrations ./packages/core/db/migrations

EXPOSE 8080

CMD ["./bastion-server"]
