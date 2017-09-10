FROM golang:latest
MAINTAINER TomoProg <helloworld0306.xxx@gmail.com>
RUN apt-get update && apt-get install -y \
	vim-tiny
RUN go get "github.com/google/go-github/github"
RUN go get "golang.org/x/oauth2"
RUN go get "github.com/nlopes/slack"
WORKDIR /go/app
