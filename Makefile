GOPATH=$(shell pwd)/vendor

go-rss2imap: *.go
	go build

deps:
	mkdir -p vendor
	go get github.com/htr/feedparser
	go get github.com/rakyll/globalconf
	go get github.com/sloonz/go-mime-message
	go get github.com/sloonz/go-qprintable
	go get code.google.com/p/go-imap/go1/imap

clean:
	rm -f go-rss2imap
	find vendor -name .git | xargs rm -rf
	find vendor -name .hg | xargs rm -rf
	go fmt *.go

pi:
	GOARCH=arm GOARM=6 go build

synology:
	GOARCH=arm GOARM=5 go build
