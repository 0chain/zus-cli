ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

SAMPLE_DIR:=$(ROOT_DIR)/sample

ZUSCLI=zus-cli

.PHONY:

# GOMODCORE           := $(GOMODBASE)/zcncore
# VERSION_FILE        := $(ROOT_DIR)/core/version/version.go
# MAJOR_VERSION       := "1.0"

include _util/printer.mk

default: help

#GO BUILD SDK
gomod-download:
	go env
	cat go.mod
	go mod download
	go mod tidy

gomod-clean:
	go clean -i -r -x -modcache  ./...

$(ZUSCLI): gomod-download
	$(eval VERSION=$(shell git describe --tags --dirty --always))
	CGO_ENABLED=1 go build -x -v -tags bn256 -ldflags "-X main.VersionStr=$(VERSION)" -o $(ZUSCLI) main.go

zus-cli-test:
	CGO_ENABLED=1 go test -tags bn256 ./...

install: $(ZUSCLI) zus-cli-test

clean: gomod-clean
	@rm -rf $(ROOT_DIR)/$(ZUSCLI)

help:
	@echo "Environment: "
	@echo "\tGOPATH=$(GOPATH)"
	@echo "\tGOROOT=$(GOROOT)"
	@echo ""
	@echo "Supported commands:"
	@echo "\tmake help              - display environment and make targets"
	@echo ""
	@echo "Install"
	@echo "\tmake install           - build, test and install zus-cli"
	@echo "\tmake zus-cli              - installs the zus-cli"
	@echo "\tmake zus-cli-test      - run zus-cli test"
	@echo ""
	@echo "Clean:"
	@echo "\tmake clean             - deletes all build output files"
	@echo "\tmake gomod-download    - download the go modules"
	@echo "\tmake gomod-clean       - clean the go modules"
