GO := go
NAME := lioss
VERSION := 0.9.0
DIST := $(NAME)-$(VERSION)

all: test build

setup:
	git submodule update --init

update_version:
	@for i in README.md site/content/_index.md; do\
		sed -e 's!Version-[0-9.]-yellowgreen!Version-${VERSION}-yellowgreen!g' -e 's!tag/v[0-9.]*!tag/v${VERSION}!g' $$i > a ; mv a $$i; \
	done
	@sed 's/const VERSION = .*/const VERSION = "${VERSION}"/g' cmd/lioss/main.go > a ; mv a cmd/lioss/main.go
	@echo "Replace version to \"${VERSION}\""

start:
	make -C site start

stop:
	make -C site stop

www:
	make -C site build

test: setup update_version
	$(GO) test -covermode=count -coverprofile=coverage.out $$(go list ./...)

build: test
	$(GO) build -o lioss -v cmd/lioss/main.go
	$(GO) build -o mkliossdb -v cmd/mkliossdb/main.go

define _createDist
	mkdir -p dist/$(1)_$(2)/$(DIST)
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/lioss cmd/lioss/main.go
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/mkliossdb cmd/mkliossdb/main.go
	cp -r README.md LICENSE dist/$(1)_$(2)/$(DIST)
	cp testdata/liossdb.json dist/$(1)_$(2)/$(DIST)
	tar cfz dist/$(DIST)_$(1)_$(2).tar.gz -C dist/$(1)_$(2) $(DIST)
endef

dist:
	@$(call _createDist,darwin,386)
	@$(call _createDist,darwin,amd64)
	@$(call _createDist,windows,amd64)
	@$(call _createDist,windows,386)
	@$(call _createDist,linux,amd64)
	@$(call _createDist,linux,386)

clean:
	$(GO) clean
	rm -rf $(NAME) mkliossdb

distclean: clean
	-rm -rf dist
