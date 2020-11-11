ARG GO_VERSION="1.15"

FROM golang:${GO_VERSION}-alpine

RUN apk add --no-cache git

WORKDIR $GOPATH/src/github.com/1efty/semtag
COPY . .

RUN go install

WORKDIR $GOPATH
ENTRYPOINT ["semtag"]
