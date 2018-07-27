FROM golang:1.10.0-alpine3.7 as builder

ENV name=GetMeConf \
    src=github.com/YAWAL/GetMeConf

WORKDIR $GOPATH/src/$src
COPY . $GOPATH/src/$src
RUN rm -rf $GOPATH/src/$src/vendor

RUN apk add --update curl && \
    apk add git && \
    rm -rf /var/cache/apk/*

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN dep ensure -vendor-only

RUN CGO_ENABLED=0 GOOS=linux go build -o $name -a -ldflags '-extldflags "-static"' main.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/YAWAL/GetMeConf/GetMeConf /
COPY --from=builder /go/src/github.com/YAWAL/GetMeConf/.env /

CMD ["/GetMeConf"]

EXPOSE $SERVICE_PORT