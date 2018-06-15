FROM alpine:3.6

ADD ./bin  .

CMD ["./service"]

EXPOSE $SERVICE_PORT