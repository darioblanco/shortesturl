# Please keep at the top.
SHELL := /usr/bin/env bash
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

.PHONY: all benchmark build coverage format help init init-deps init-godeps install gen run run-hmr test

all: init help

benchmark: init ## execute benchmarks
	go test -bench=. -benchmem ./app/internal/http

build: init gen ## build the go app
	go build -o tmp/shortesturl ./cmd/server/main.go

coverage: init ## generate coverage reports
	go tool cover -html=coverage.out

format: init ## format the go app
	go fmt ./...

help: ## this help output
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

init: init-deps init-godeps ## initialize the repository after a fresh clone (runs implicitly)

init-deps:
	@function cmd { \
		if ! command -v "$$1" &>/dev/null ; then \
			echo "error: missing required command in PATH: $$1" >&2 ;\
			return 1 ;\
		fi \
	} ;\
	cmd go ;\
	cmd git ;\
	$(if $(PRODUCTION),true,false) || cmd docker ;

init-godeps:
	@function godep { \
		if ! command -v "$$1" &>/dev/null ; then \
			set -x ; go install "$$2" ; set +x ;\
		fi \
	} ;\
	$(if $(PRODUCTION),true,false) || godep air github.com/cosmtrek/air ;\
	$(if $(PRODUCTION),true,false) || godep goverreport github.com/mcubik/goverreport ;\
	godep swag github.com/swaggo/swag/cmd/swag ;\

install: init ## install the dependencies
	go install ./cmd/server/main.go

gen: ## generate swagger documentation
	swag init -d app/internal/http -g docs.go -o docs

run-hmr: init build ## run the go app with live-reloading enabled
	ulimit -n 65535; air

run: init build ## run the go app
	go run ./cmd/server/main.go

test: init ## execute tests
	go test -count=1 -race -timeout 10s ./... -cover -coverprofile coverage.out.tmp
	grep -v "_mock.go" coverage.out.tmp > coverage.out && rm coverage.out.tmp
	goverreport
