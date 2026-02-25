FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/bastion-cli ./bastion

ENTRYPOINT ["./bastion"]
