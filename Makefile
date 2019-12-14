###########################################################
#
# Copyright (c) 2019 codeliveroil. All rights reserved.
#
# This work is licensed under the terms of the MIT license.
# For a copy, see <https://opensource.org/licenses/MIT>.
#
###########################################################

BUILD_DIR := "build"
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
BINARY_NAME := $(notdir $(patsubst %/,%,$(dir $(MKFILE_PATH))))
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

# Exposed targets
.PHONY: all build install test clean platform dist setup dist_internal help

build: setup #!## Build binary
	go build -v -o $(BINARY_PATH)

all: test build #!## Test and build

install: #!## Install binary after running 'make build'.
	@path2=":$${PATH}:" &&                                 \
	if [ -z "${path2##*:/usr/local/bin:*}" ]; then         \
		cp ${BINARY_PATH} /usr/local/bin;                  \
	elif [ -z "${path2##*:/usr/bin:*}" ]; then             \
		cp ${BINARY_PATH} /usr/bin;                        \
	elif [ -z "${path2##*:/bin:*}" ]; then                 \
		cp ${BINARY_PATH} /bin;                            \
	else                                                   \
		echo "Standard Unix installation path not found."; \
	fi

test: clean #!## Run tests
	go test -v ./...

clean: #!## Clean build environment
	go clean ./...
	rm -rf $(BUILD_DIR)

platform: setup #!## Cross compile to a desired platform
	@echo "Select a platform:"
	@select p in $$(go tool dist list); do                                     \
		[ "$$p" == "" ] && echo "Invalid selection" && exit;                   \
		IFS='/' tokens=( $$p );                                                \
		GOOS=$${tokens[0]} GOARCH=$${tokens[1]} go build -v -o $(BINARY_PATH); \
		break;                                                                 \
	done

dist: test setup
	@make goos=darwin goarch=amd64 name=macos dist_internal
	@make goos=linux goarch=amd64 name=linux dist_internal
	@make goos=linux goarch=arm name=linux-arm dist_internal

# Internal targets
setup: clean
	mkdir $(BUILD_DIR)

dist_internal:
	go clean ./...
	@GOOS=${goos} GOARCH=${goarch} go build -v -o $(BINARY_PATH) \
	&& cd $(BUILD_DIR)                                           \
	&& zip ${BINARY_NAME}-${name} ${BINARY_NAME}                 \
	&& rm ${BINARY_NAME}

# Special targets
help:
	@grep '#!#' $(MAKEFILE_LIST) | grep -v "#NA#" \
	| sort | uniq                                 \
	| sed s/':.* \#!\#'/'#'/g                     \
	| sed -e 's/^[[:space:]]*/ /'                 \
	| column -t -s '#'
