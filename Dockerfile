FROM golang:1.10.0-alpine3.7 as builder

WORKDIR $GOPATH/src/$svc
COPY . $GOPATH/src/$svc
RUN rm -rf $GOPATH/src/$svc/vendor

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/vX.X.X/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/***
WORKDIR /go/src/github.com/***

COPY Gopkg.toml Gopkg.lock ./
# copies the Gopkg.toml and Gopkg.lock to WORKDIR

RUN dep ensure -vendor-only



FROM alpine:latest
COPY --from=builder /go/src/github.com/GetMeConf /


CMD ["./service"]

EXPOSE $SERVICE_PORT