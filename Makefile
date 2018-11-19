#-------------------------------------------------------------------------------
# VARIABLES
#-------------------------------------------------------------------------------

NAME      = cbdyncluster
DIST_NAME ?= $(NAME)
MAIN_PATH = github.com/couchbaselabs/cbdyncluster

VERSION = $(shell git describe --always --tags)
GOFLAGS = -ldflags "-X ${MAIN_PATH}/cmd.Version=${VERSION}"

GOOS 	?= darwin
GOARCH  ?= amd64

#-------------------------------------------------------------------------------
# SPECIAL
#-------------------------------------------------------------------------------

.DEFAULT: default
.PHONY: default
default: clean build

.PHONY: all
all: clean build-all

.PHONY: release
release: clean build-all

#-------------------------------------------------------------------------------
# EXECUTABLE
#-------------------------------------------------------------------------------

.PHONY: clean
clean:
	go clean .
	rm -rf bin

.PHONY: run
run:
	./bin/$(NAME).$(GOOS)_$(GOARCH)

.PHONY: install
install:
	cp bin/$(NAME).$(GOOS)_$(GOARCH) ~/bin/$(DIST_NAME)

#-------------------------------------------------------------------------------
# BUILD
#-------------------------------------------------------------------------------

.PHONY: build
build:
	$(GOENV) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -o bin/$(NAME).$(GOOS)_$(GOARCH)
	chmod +x bin/$(NAME).$(GOOS)_$(GOARCH)

.PHONY:	build-all
build-all: build-mac build-linux

.PHONY: build-mac
build-mac: GOOS=darwin
build-mac:
	@$(MAKE) build GOOS=$(GOOS) GOARCH=$(GOARCH)

.PHONY: build-linux
build-linux: GOOS=linux
build-linux:
	@$(MAKE) build GOOS=$(GOOS) GOARCH=$(GOARCH)
