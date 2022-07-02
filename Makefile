SHELL = /bin/bash

NAME := jx-git-operator
ORG := jenkins-x
ORG_REPO := $(ORG)/$(NAME)
RELEASE_ORG_REPO := $(ORG_REPO)
REV := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
ROOT_PACKAGE := github.com/$(ORG_REPO)
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)
GO_VERSION := 1.17.11

GO := GO111MODULE=on go
BUILD_TARGET = build
CGO_ENABLED = 0

GO_NOMOD :=GO111MODULE=off go

# set dev version unless VERSION is explicitly set via environment
#VERSION ?= $(shell echo "$$(git for-each-ref refs/tags/ --count=1 --sort=-version:refname --format='%(refname:short)' 2>/dev/null)-dev+$(REV)" | sed 's/^v//')

GOPRIVATE := github.com/jenkins-x/jx-helpers

MAIN_SRC_FILE=./main.go

BUILDFLAGS :=  -ldflags \
  " -X main.buildTime=$(BUILD_DATE) \
		-X main.gitCommit=$(REV) \
		-X main.version=$(VERSION)"

REPORTS_DIR=bin

COVER_OUT:=$(REPORTS_DIR)/cover.out
COVERFLAGS=-coverprofile=$(COVER_OUT) --covermode=count --coverpkg=./...


.PHONY: build
build:
	go build -o bin/$(NAME) main.go

release: build test

.PHONY: release
release: clean linux test

release-all: release linux win darwin

.PHONY: goreleaser
goreleaser:
	step-go-releaser --organisation=$(ORG) --revision=$(REV) --branch=$(BRANCH) --build-date=$(BUILD_DATE) --go-version=$(GO_VERSION) --root-package=$(ROOT_PACKAGE) --version=$(VERSION)

fmt:
	go fmt ./...

test:
	go test ./... --tags="integration unit"

test-coverage:
	go test --tags="integration unit" -v $(COVERFLAGS) ./...

cover:
	go tool cover -html=$(COVER_OUT)

test-report: make-reports-dir test-coverage ## Create the test report
	gocov convert $(COVER_OUT) | gocov report

test-report-html: make-reports-dir test-coverage ## Create the test report in HTML format
	gocov convert $(COVER_OUT) | gocov-html > $(REPORTS_DIR)/cover.html && open $(REPORTS_DIR)/cover.html

.PHONY: make-reports-dir
make-reports-dir:
	mkdir -p $(REPORTS_DIR)

linux: ## Build for Linux
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x build/linux/$(NAME)

.PHONY: clean
clean: ## Clean the generated artifacts
	rm -rf bin release dist

vet:
	go vet ./...

tools:
	$(GO_NOMOD) get github.com/axw/gocov/gocov
	$(GO_NOMOD) get -u gopkg.in/matm/v1/gocov-html
	$(GO_NOMOD) get golang.org/x/tools/cmd/goimports
	$(GO_NOMOD) get github.com/kisielk/errcheck
	go get honnef.co/go/tools/cmd/staticcheck
	$(GO_NOMOD) get github.com/golang/lint/golint

errors:
	errcheck -ignoretests -blank ./...

lint2:
	golint ./...

lint:
	golangci-lint run

imports:
	goimports -l -w .


all: fmt imports test lint vet

