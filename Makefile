GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/vault-plugin-secrets-vercel cmd/vault-plugin-secrets-vercel/main.go

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

.PHONY: build clean fmt start enable lint
