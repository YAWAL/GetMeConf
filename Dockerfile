FROM alpine:latest

ADD ./bin  .

CMD ["./service"]

EXPOSE $SERVICE_PORT