# STEP 1 : Build web-crawler
FROM golang:1.11 AS buildstage

LABEL author="nikhilthomas1@gmail.com"

RUN go get -u github.com/golang/dep/...

WORKDIR /go/src/github.com/nikhil-thomas/web-crawler

# install dependencies before copying source code
# on subsequent builds layers till here will be taken from cache
# builds will be faster as dependencies are not fetched during each change
# cache will break only when dependencies are added/emoved/modified
ADD ./Gopkg.lock ./
ADD ./Gopkg.toml ./
RUN dep ensure -vendor-only -v

ADD . .
WORKDIR /go/src/github.com/nikhil-thomas/web-crawler/cmd/configurable-crawler
RUN CGO_ENABLED=0 GOOS=linux go build -o web-crawler

# STEP 1 : Package web-crawler
FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN update-ca-certificates

ENV GOPATH /go

WORKDIR /go/bin

COPY --from=buildstage /go/src/github.com/nikhil-thomas/web-crawler/cmd/configurable-crawler/web-crawler .

ENTRYPOINT [ "/go/bin/web-crawler"  ]
