export GOPATH:=$(shell pwd)/vendor
export PATH:=$(PATH):$(GOPATH)/bin
BINARY=rss2imap

$(BINARY): *.go
	go build -o $(BINARY)

deps:
	mkdir -p vendor
	go get github.com/htr/feedparser
	go get github.com/rakyll/globalconf
	go get github.com/sloonz/go-mime-message
	go get github.com/sloonz/go-qprintable
	go get code.google.com/p/go-imap/go1/imap
	find vendor -name .git | xargs rm -rf
	find vendor -name .hg | xargs rm -rf

clean:
	rm -f $(BINARY)
	go fmt *.go

pi:
	GOOS=linux GOARCH=arm GOARM=6 go build -o $(BINARY)

synology:
	GOOS=linux GOARCH=arm GOARM=5 go build -o $(BINARY)
