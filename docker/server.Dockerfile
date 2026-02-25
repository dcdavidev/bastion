FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY bastion-server .
COPY internal/db/migrations ./internal/db/migrations

EXPOSE 8080

CMD ["./bastion-server"]
