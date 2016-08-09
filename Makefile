ROOT_GOPATH=$(shell pwd)
VENDOR_GOPATH=$(shell pwd)/vendor
GOPATH=$(VENDOR_GOPATH)
APP_NAME=goslash-sample-server

install:
	mkdir -p $(VENDOR_GOPATH)
	go get golang.org/x/tools/cmd/cover
	go get github.com/stretchr/testify
	go get golang.org/x/net/context
	go get -u github.com/unrolled/render
	go get -u github.com/kyokomi/goslash/goslash
	go get -u github.com/kyokomi/goslash/plugins/time
	go get -u github.com/kyokomi/goslash/plugins/echo
	go get -u github.com/kyokomi/goslash/plugins/suddendeath
	go get -u github.com/kyokomi/goslash/plugins/lgtm
	go get -u github.com/kyokomi/goslash/plugins/akari
	go get -u github.com/PuerkitoBio/goquery

clean:
	rm $(APP_NAME)
	rm -rf $(VENDOR_GOPATH)

fmt:
	gofmt -w ./src/server.go

run:
	go run ./src/server.go -dev

build:
	go build -o $(APP_NAME) ./src/server.go

test:
	GOPATH=$(VENDOR_GOPATH):$(ROOT_GOPATH) go test -v ./src/...
