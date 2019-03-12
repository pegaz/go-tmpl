OUT := go-tmpl
PKG := github.com/pegaz/go-tmpl
VERSION := $(shell git describe --always --tag --long --dirty)
PKG_LIST := $(shell go list ${PKG}/...)

all: build

build:
	go build -i -v -o ${OUT} -ldflags="-X ${PKG}/cmd.version=${VERSION}" ${PKG}

clean:
	-@rm ${OUT}

test:
	@go test -short ${PKG_LIST}
