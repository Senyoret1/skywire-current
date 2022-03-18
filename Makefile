
.PHONY : check lint install-linters dep test
.PHONY : build clean install format  bin
.PHONY : host-apps bin
.PHONY : docker-image docker-clean docker-network
.PHONY : docker-apps docker-bin docker-volume
.PHONY : docker-run docker-stop

VERSION := $(shell git describe)
RFC_3339 := "+%Y-%m-%dT%H:%M:%SZ"
COMMIT := $(shell git rev-list -1 HEAD)

PROJECT_BASE := github.com/skycoin/skywire
DMSG_BASE := github.com/skycoin/dmsg
ifeq ($(OS),Windows_NT)
	SHELL := pwsh
	OPTS?=powershell -Command setx GO111MODULE on;
	DATE := $(shell powershell -Command date -u ${RFC_3339})
	.DEFAULT_GOAL := help-windows
else
	SHELL := /bin/bash
	OPTS?=GO111MODULE=on
	DATE := $(shell date -u $(RFC_3339))
	.DEFAULT_GOAL := help
endif

STATIC_OPTS?= $(OPTS) CC=musl-gcc
MANAGER_UI_DIR = static/skywire-manager-src
GO_BUILDER_VERSION=v1.17
MANAGER_UI_BUILT_DIR=cmd/skywire-visor/static

TEST_OPTS:=-cover -timeout=5m -mod=vendor

GOARCH:=$(shell go env GOARCH)

ifneq (,$(findstring 64,$(GOARCH)))
    TEST_OPTS:=$(TEST_OPTS) -race
endif

BUILDINFO_PATH := $(DMSG_BASE)/buildinfo

BUILDINFO_VERSION := -X $(BUILDINFO_PATH).version=$(VERSION)
BUILDINFO_DATE := -X $(BUILDINFO_PATH).date=$(DATE)
BUILDINFO_COMMIT := -X $(BUILDINFO_PATH).commit=$(COMMIT)
BUILDTAGINFO := -X $(PROJECT_BASE)/pkg/visor.BuildTag=$(BUILDTAG)

BUILDINFO?=$(BUILDINFO_VERSION) $(BUILDINFO_DATE) $(BUILDINFO_COMMIT) $(BUILDTAGINFO)

BUILD_OPTS?="-ldflags=$(BUILDINFO)" -mod=vendor $(RACE_FLAG)
BUILD_OPTS_DEPLOY?="-ldflags=$(BUILDINFO) -w -s"

check: lint test ## Run linters and tests

check-windows: lint-windows test-windows ## Run linters and tests on appveyor windows image

build: host-apps bin ## Install dependencies, build apps and binaries. `go build` with ${OPTS}

build-windows: host-apps-windows bin-windows ## Install dependencies, build apps and binaries. `go build` with ${OPTS}

build-systray: host-apps-systray bin-systray ## Install dependencies, build apps and binaries `go build` with ${OPTS}, with CGO and systray

build-systray-windows: host-apps-systray-windows bin-systray-windows ## Builds systray binary in windows

build-static: host-apps-static bin-static ## Build apps and binaries. `go build` with ${OPTS}

installer: mac-installer ## Builds MacOS installer for skywire-visor

install-generate: ## Installs required execs for go generate.
	${OPTS} go install github.com/mjibson/esc
	${OPTS} go install github.com/vektra/mockery/cmd/mockery

generate: ## Generate mocks and config README's
	go generate ./...

clean: ## Clean project: remove created binaries and apps
	-rm -rf ./apps
	-rm -f ./skywire-visor ./skywire-cli ./setup-node

clean-windows: ## Clean project: remove created binaries and apps
	powershell -Command Remove-Item -Path ./apps -Force -Recurse
	powershell -Command Remove-Item -Path .\skywire-visor.exe,.\skywire-cli.exe,.\setup-node.exe -Force

install: ## Install `skywire-visor`, `skywire-cli`, `setup-node`
	${OPTS} go install ${BUILD_OPTS} ./cmd/skywire-visor ./cmd/skywire-cli ./cmd/setup-node

install-windows: ## Install `skywire-visor`, `skywire-cli`, `setup-node`
	powershell 'Get-ChildItem .\cmd | % { ${OPTS} go install ${BUILD_OPTS} ./ $$_.FullName }'

install-static: ## Install `skywire-visor`, `skywire-cli`, `setup-node`
	${STATIC_OPTS} go install -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' ./cmd/skywire-visor ./cmd/skywire-cli ./cmd/setup-node

lint: ## Run linters. Use make install-linters first
	${OPTS} golangci-lint run -c .golangci.yml ./...

lint-windows: ## Run linters. Use make install-linters-windows first
	powershell 'golangci-lint run -c .golangci.yml ./...'

lint-appveyor-windows: ## Run linters for appveyor only on windows
	C:\Users\appveyor\go\bin\golangci-lint run -c .golangci.yml ./...

test: ## Run tests
	-go clean -testcache &>/dev/null
	${OPTS} go test ${TEST_OPTS} ./internal/...
	${OPTS} go test ${TEST_OPTS} ./pkg/...

test-windows: ## Run tests on windows
	@go clean -testcache
	${OPTS} go test ${TEST_OPTS} ./internal/...
	${OPTS} go test ${TEST_OPTS} ./pkg/...

install-linters: ## Install linters
	- VERSION=latest ./ci_scripts/install-golangci-lint.sh
	${OPTS} go install golang.org/x/tools/cmd/goimports@latest
	${OPTS} go install github.com/incu6us/goimports-reviser/v2@latest

install-linters-windows: ## Install linters
	${OPTS} go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	${OPTS} go install golang.org/x/tools/cmd/goimports@latest

tidy: ## Tidies and vendors dependencies.
	${OPTS} go mod tidy -v

format: tidy ## Formats the code. Must have goimports and goimports-reviser installed (use make install-linters).
	${OPTS} goimports -w -local ${PROJECT_BASE} ./pkg
	${OPTS} goimports -w -local ${PROJECT_BASE} ./cmd
	${OPTS} goimports -w -local ${PROJECT_BASE} ./internal
	find . -type f -name '*.go' -not -path "./.git/*" -not -path "./vendor/*"  -exec goimports-reviser -project-name ${PROJECT_BASE} -file-path {} \;

format-windows: tidy ## Formats the code. Must have goimports and goimports-reviser installed (use make install-linters).
	powershell 'Get-ChildItem -Directory | where Name -NotMatch vendor | % { Get-ChildItem $$_ -Recurse -Include *.go } | % {goimports -w -local ${PROJECT_BASE} $$_ }'

dep: tidy ## Sorts dependencies
	${OPTS} go mod vendor -v

snapshot:
	goreleaser --snapshot --skip-publish --rm-dist

snapshot-clean: ## Cleans snapshot / release
	rm -rf ./dist

host-apps: ## Build app
	mkdir -p ./apps
	${OPTS} go build ${BUILD_OPTS} -o ./apps/ ./cmd/apps/skychat
	${OPTS} go build ${BUILD_OPTS} -o ./apps/ ./cmd/apps/skysocks
	${OPTS} go build ${BUILD_OPTS} -o ./apps/ ./cmd/apps/skysocks-client
	${OPTS} go build ${BUILD_OPTS} -o ./apps/ ./cmd/apps/vpn-server
	${OPTS} go build ${BUILD_OPTS} -o ./apps/ ./cmd/apps/vpn-client

host-apps-windows:
	powershell -Command new-item .\apps -itemtype directory -force
	powershell 'Get-ChildItem .\cmd\apps | % { ${OPTS} go build ${BUILD_OPTS} -o ./apps $$_.FullName }'

host-apps-systray: ## Build app
	${OPTS} go build ${BUILD_OPTS} -o ./apps/ ./cmd/apps/skychat
	${OPTS} go build ${BUILD_OPTS} -o ./apps/ ./cmd/apps/skysocks
	${OPTS} go build ${BUILD_OPTS} -o ./apps/  ./cmd/apps/skysocks-client
	${OPTS} go build ${BUILD_OPTS} -tags systray -o ./apps/ ./cmd/apps/vpn-server
	${OPTS} go build ${BUILD_OPTS} -tags systray -o ./apps/ ./cmd/apps/vpn-client

host-apps-systray-windows:
	powershell -Command new-item .\apps -itemtype directory -force
	powershell 'go build ${BUILD_OPTS} -o .\apps\skychat.exe .\cmd\apps\skychat'
	powershell 'go build ${BUILD_OPTS} -o .\apps\skysocks.exe .\cmd\apps\skysocks'
	powershell 'go build ${BUILD_OPTS} -o .\apps\skysocks-client.exe .\cmd\apps\skysocks-client'
	powershell 'go build ${BUILD_OPTS} -tags systray -o .\apps\vpn-server.exe .\cmd\apps\vpn-server'
	powershell 'go build ${BUILD_OPTS} -tags systray -o .\apps\vpn-client.exe .\cmd\apps\vpn-client'

# Static Apps
host-apps-static: ## Build app
	mkdir -p ./apps
	${STATIC_OPTS} go build -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' -o ./apps/ ./cmd/apps/skychat
	${STATIC_OPTS} go build -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' -o ./apps/ ./cmd/apps/skysocks
	${STATIC_OPTS} go build -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' -o ./apps/ ./cmd/apps/skysocks-client
	${STATIC_OPTS} go build -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' -o ./apps/ ./cmd/apps/vpn-server
	${STATIC_OPTS} go build -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' -o ./apps/ ./cmd/apps/vpn-client

# Bin
bin: ## Build `skywire-visor`, `skywire-cli`
	${OPTS} go build ${BUILD_OPTS} -o ./ ./cmd/skywire-visor
	${OPTS} go build ${BUILD_OPTS} -o ./ ./cmd/skywire-cli
	${OPTS} go build ${BUILD_OPTS} -o ./ ./cmd/setup-node

bin-windows: ## Build `skywire-visor`, `skywire-cli`
	powershell 'Get-ChildItem .\cmd | % { ${OPTS} go build ${BUILD_OPTS} -o ./ $$_.FullName }'

bin-systray-windows: ## Build `skywire-visor` and `skywire-cli` with systray support
	powershell 'Get-ChildItem .\cmd | % { ${OPTS} go build ${BUILD_OPTS} -tags systray -o ./ $$_.FullName }'

bin-systray: ## Build `skywire-visor`, `skywire-cli`
	${OPTS} go build ${BUILD_OPTS} -tags systray -o ./ ./cmd/skywire-visor
	${OPTS} go build ${BUILD_OPTS} -tags systray -o ./ ./cmd/skywire-cli
	${OPTS} go build ${BUILD_OPTS} -o ./ ./cmd/setup-node

# Static Bin
bin-static: ## Build `skywire-visor`, `skywire-cli`
	${STATIC_OPTS} go build -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' -o ./skywire-visor ./cmd/skywire-visor
	${STATIC_OPTS} go build -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' -o ./skywire-cli  ./cmd/skywire-cli
	${STATIC_OPTS} go build -trimpath --ldflags '-linkmode external -extldflags "-static" -buildid=' -o ./setup-node ./cmd/setup-node

build-deploy: ## Build for deployment Docker images
	${OPTS} go build -tags netgo ${BUILD_OPTS_DEPLOY} -o /release/skywire-visor ./cmd/skywire-visor
	${OPTS} go build ${BUILD_OPTS_DEPLOY} -o /release/skywire-cli ./cmd/skywire-cli
	${OPTS} go build ${BUILD_OPTS_DEPLOY} -o /release/apps/skychat ./cmd/apps/skychat
	${OPTS} go build ${BUILD_OPTS_DEPLOY} -o /release/apps/skysocks ./cmd/apps/skysocks
	${OPTS} go build ${BUILD_OPTS_DEPLOY} -o /release/apps/skysocks-client ./cmd/apps/skysocks-client

github-release:
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags | cut -c 2-))
	sed '/^## ${GITHUB_TAG}$$/,/^## .*/!d;//d;/^$$/d' ./CHANGELOG.md > releaseChangelog.md
	goreleaser --rm-dist --release-notes releaseChangelog.md


build-docker: ## Build docker image
	./ci_scripts/docker-push.sh -t latest -b

# Manager UI
install-deps-ui:  ## Install the UI dependencies
	cd $(MANAGER_UI_DIR) && npm ci

run: ## Run skywire visor with skywire-config.json, and start a browser if running a hypervisor
	./skywire-visor -c ./skywire-config.json

## Run skywire from source, without compiling binaries - requires skywire cloned
run-source:
	test -d apps && rm -r apps || true
	ln -s scripts/_apps apps
	chmod +x apps/*
	go run ./cmd/skywire-cli/skywire-cli.go config gen -ibro ./skywire-config.json
	go run ./cmd/skywire-visor/skywire-visor.go -c ./skywire-config.json || true

lint-ui:  ## Lint the UI code
	cd $(MANAGER_UI_DIR) && npm run lint

build-ui: install-deps-ui  ## Builds the UI
	cd $(MANAGER_UI_DIR) && npm run build
	mkdir -p ${PWD}/bin
	rm -rf ${MANAGER_UI_BUILT_DIR}
	mkdir ${MANAGER_UI_BUILT_DIR}
	cp -r ${MANAGER_UI_DIR}/dist/. ${MANAGER_UI_BUILT_DIR}

build-ui-windows: install-deps-ui ## Builds the UI on windows
	cd $(MANAGER_UI_DIR) && npm run build
	powershell 'Remove-Item -Recurse -Force -Path ${MANAGER_UI_BUILT_DIR}'
	powershell 'New-Item -Path ${MANAGER_UI_BUILT_DIR} -ItemType Directory'
	powershell 'Copy-Item -Recurse ${MANAGER_UI_DIR}\dist\* ${MANAGER_UI_BUILT_DIR}'

deb-install-prequisites: ## Create unsigned application
	sudo chmod +x ./scripts/deb_installer/prequisites.sh
	./scripts/deb_installer/prequisites.sh

deb-package: deb-install-prequisites ## Create unsigned application
	./scripts/deb_installer/package_deb.sh

deb-package-help: ## Show installer creation help
	./scripts/deb_installer/package_deb.sh -h

mac-installer: ## Create signed and notarized application, run make mac-installer-help for more
	./scripts/mac_installer/create_installer.sh -s -n

mac-installer-help: ## Show installer creation help
	./scripts/mac_installer/create_installer.sh -h

win-installer:
	@powershell '.\scripts\win_installer\script.ps1'

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

help-windows: ## Display help for windows
	@powershell 'Select-String -Pattern "windows[a-zA-Z_-]*:.*## .*$$" $(MAKEFILE_LIST) | % { $$_.Line -split ":.*?## " -Join "`t:`t" } '
