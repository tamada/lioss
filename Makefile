GO := go
NAME := lioss
VERSION := 1.0.0
DIST := $(NAME)-$(VERSION)

all: test build

setup:
	git submodule update --init

update_version:
	@for i in README.md ; do\
		sed -e 's!Version-[0-9.]-yellowgreen!Version-${VERSION}-yellowgreen!g' -e 's!tag/v[0-9.]*!tag/v${VERSION}!g' $$i > a ; mv a $$i; \
	done
	@sed 's/ARG version=".*"/ARG version="${VERSION}"/g' Dockerfile > a ; mv a Dockerfile
	@sed 's/const VERSION = .*/const VERSION = "${VERSION}"/g' cmd/lioss/main.go > a ; mv a cmd/lioss/main.go
	@sed 's/lioss version .*/lioss version ${VERSION}/g' cmd/lioss/main_test.go > a ; mv a cmd/lioss/main_test.go
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
	$(GO) build -o lioss -v cmd/lioss/main.go cmd/lioss/validator.go
	$(GO) build -o mkliossdb -v cmd/mkliossdb/main.go
	$(GO) build -o spdx2liossdb -v cmd/spdx2liossdb/main.go

createdb: build
	./spdx2liossdb -d data/NoneOSIApproved.liossgz spdx/src --without-deprecated   --without-osi-approved
	./spdx2liossdb -d data/OSIDeprecated.liossgz   spdx/src --with-deprecated      --with-osi-approved
	./spdx2liossdb -d data/OSIApproved.liossgz     spdx/src --without-deprecated   --with-osi-approved
	./spdx2liossdb -d data/Deprecated.liossgz      spdx/src --with-deprecated      --without-osi-approved

define _createDist
	mkdir -p dist/$(1)_$(2)/$(DIST)/data
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/lioss$(3) cmd/lioss/main.go cmd/lioss/validator.go
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/mkliossdb$(3) cmd/mkliossdb/main.go
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/spdx2liossdb$(3) cmd/spdx2liossdb/main.go
	cp -r README.md LICENSE dist/$(1)_$(2)/$(DIST)
	cp data/*.liossgz dist/$(1)_$(2)/$(DIST)/data
	tar cfz dist/$(DIST)_$(1)_$(2).tar.gz -C dist/$(1)_$(2) $(DIST)
endef

dist: build createdb
	@$(call _createDist,darwin,amd64,)
	@$(call _createDist,windows,amd64,.exe)
	@$(call _createDist,windows,386,.exe)
	@$(call _createDist,linux,amd64,)
	@$(call _createDist,linux,386,)

clean:
	$(GO) clean
	rm -rf $(NAME) mkliossdb

distclean: clean
	-rm -rf dist
