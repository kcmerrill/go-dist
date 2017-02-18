FROM golang:1.6
MAINTAINER kc merrill <kcmerrill@gmail.com>
RUN go get github.com/mitchellh/gox
COPY . /go/src/github.com/kcmerrill/go-dist
WORKDIR /go/src/github.com/kcmerrill/go-dist
RUN  go build -ldflags "-X main.Commit=`git rev-parse HEAD` -X main.Version=0.1.`git rev-list --count HEAD`" -o /usr/local/bin/go-dist
EXPOSE 80
WORKDIR /tmp
ENTRYPOINT ["go-dist"]
