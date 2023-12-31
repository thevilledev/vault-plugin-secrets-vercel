GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

DATE 		=$(shell date '+%a %b %d %H:%m:%S %Z %Y')
REVISION 	=$(shell git rev-parse --verify --short HEAD)
VERSION 	=$(shell git describe --always --tags --exact-match 2>/dev/null || \
				echo $(REVISION))
BRANCH		=$(shell git rev-parse --abbrev-ref HEAD)
DIRTY		=false

ifneq ($(shell git status --porcelain),)
    DIRTY = true
endif

LDFLAGS =-s -w -extld ld -extldflags -static \
		  -X 'github.com/thevilledev/vault-plugin-secrets-vercel/internal/version.BuildDate=$(DATE)' \
		  -X 'github.com/thevilledev/vault-plugin-secrets-vercel/internal/version.Version=$(VERSION)' \
		  -X 'github.com/thevilledev/vault-plugin-secrets-vercel/internal/version.Commit=$(REVISION)' \
		  -X 'github.com/thevilledev/vault-plugin-secrets-vercel/internal/version.CommitDate=$(REVISION)' \
		  -X 'github.com/thevilledev/vault-plugin-secrets-vercel/internal/version.Branch=$(BRANCH)' \
		  -X 'github.com/thevilledev/vault-plugin-secrets-vercel/internal/version.Tag=$(VERSION)' \
		  -X 'github.com/thevilledev/vault-plugin-secrets-vercel/internal/version.Dirty=$(DIRTY)'
FLAGS	=-trimpath -a -ldflags "$(LDFLAGS)"

.DEFAULT_GOAL := all

all: fmt build start

build:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH="$(GOARCH)" go build $(FLAGS) -o vault/plugins/vault-plugin-secrets-vercel cmd/vault-plugin-secrets-vercel/main.go

start:
	VAULT_ADDR='http://127.0.0.1:8200' VAULT_API_ADDR='http://127.0.0.1:8200' vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins

enable:
	vault secrets enable -path=vercel-secrets vault-plugin-secrets-vercel

clean:
	rm -f ./vault/plugins/vault-plugin-secrets-vercel

fmt:
	go fmt $$(go list ./...)

lint:
	golangci-lint run

test:
	go test -v -race -parallel=8 -covermode=atomic -coverprofile=coverage.out ./...

test-acc:
	ACC_TEST=yes go test -race -parallel=4 ./...

.PHONY: build clean fmt start enable lint test test-acc
