#
# Makefile
# @author Ronald Doorn <rdoorn@schubergphilis.com>
#

.PHONY: update clean build build-all run package deploy test authors dist

export PATH := $(PATH):$(GOPATH)/bin

BINNAME := ops
PACKAGENAME := opscli
VERSION := $(shell cat VERSION)
LASTCOMMIT := $(shell git rev-parse --verify HEAD)
BUILD := $(shell cat tools/rpm/BUILDNR)
LDFLAGS := "-X main.version=$(VERSION) -X main.versionBuild=$(BUILD) -X main.versionSha=$(LASTCOMMIT)"
PENDINGCOMMIT := $(shell git diff-files --quiet --ignore-submodules && echo 0 || echo 1)
LOCALIP := $(shell ifconfig | grep "inet " | grep broadcast | awk {'print $$2'} | head -1 )

default: build

clean:
	@echo Cleaning up...
	@rm -f build
	@echo Done.

builddir:
	@mkdir -p ./build/osx/
	@mkdir -p ./build/linux/
	@mkdir -p ./build/packages/$(PACKAGENAME)

osx: builddir 
	@echo Building OSX...
	GOOS=darwin GOARCH=amd64 go build -v -o ./build/osx/$(BINNAME) -ldflags $(LDFLAGS) .
	@echo Done.

osx-fast: builddir
	@echo Building OSX skipping ...
	GOOS=darwin GOARCH=amd64 go build -v -o ./build/osx/$(BINNAME) -ldflags $(LDFLAGS) .
	@echo Done.

osx-race: builddir 
	@echo Building OSX...
	GOOS=darwin GOARCH=amd64 go build -race -v -o ./build/osx/$(BINNAME) -ldflags $(LDFLAGS) .
	@echo Done.

osx-static:
	@echo Building OSX...
	GOOS=darwin GOARCH=amd64 go build -v -o ./build/osx/$(BINNAME) -ldflags '-s -w --extldflags "-static”  $(LDFLAGS)' .
	@echo Done.

linux: builddir 
	@echo Building Linux...
	GOOS=linux GOARCH=amd64 go build -v -o ./build/linux/$(BINNAME) -ldflags '-s -w --extldflags "-static”  $(LDFLAGS)' .
	@echo Done.

build: osx linux

test:
	go test -v ./...
	go test -v ./... --race --short
	go vet ./...

cover: ## Shows coverage
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test ./internal/config -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

prep_package:
	gem install fpm

committed:
ifeq ($(PENDINGCOMMIT), 1)
	   $(error You have a pending commit, please commit your code before making a package $(PENDINGCOMMIT))
endif

# linux-package: builddir linux committed
linux-package: builddir linux
	cp -a ./tools/rpm/$(PACKAGENAME)/* ./build/packages/$(PACKAGENAME)/
	cp ./build/linux/$(BINNAME) ./build/packages/$(PACKAGENAME)/usr/sbin/
	fpm -s dir -t rpm -C ./build/packages/$(PACKAGENAME) --name $(PACKAGENAME) --rpm-os linux --version $(VERSION) --iteration $(BUILD) --exclude "*/.keepme"
	mv $(PACKAGENAME)-$(VERSION)*.rpm build/packages/

docker-scratch:
	if [ -a /System/Library/Keychains/SystemRootCertificates.keychain ] ; \
	then \
		security find-certificate /System/Library/Keychains/SystemRootCertificates.keychain > build/docker/ca-certificates.crt; \
	fi;
	if [ -a /etc/ssl/certs/ca-certificates.crt ] ; \
	then \
		cp /etc/ssl/certs/ca-certificates.crt build/docker/ca-certificates.crt; \
	fi;
	docker build -t ops-scratch -f build/docker/Dockerfile.scratch .

deps: ## Updates the vendored Go dependencies
	@dep ensure -v

updatedeps: ## Updates the vendored Go dependencies
	@dep ensure -update


#authors:
#	@git log --format='%aN <%aE>' | LC_ALL=C.UTF-8 sort | uniq -c | sort -nr | sed "s/^ *[0-9]* //g" > AUTHORS
#	@cat AUTHORS
#
