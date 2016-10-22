FROM golang:1.6
MAINTAINER kc merrill <kcmerrill@gmail.com>
RUN go get github.com/mitchellh/gox
RUN go get github.com/kcmerrill/go-dist
EXPOSE 80
ENTRYPOINT "go-gist"
