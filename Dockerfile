FROM golang:1.15.6
WORKDIR  $GOPATH/src/github.com/sreesa7144/url-shortener
COPY . .
EXPOSE 8080
CMD ["go","run","main.go"]