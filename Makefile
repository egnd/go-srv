#!make

MAKEFLAGS += --always-make --no-print-directory
CALL_PARAM=$(filter-out $@,$(MAKECMDGOALS))
BUILD_VERSION=dev
GOOS=linux
GOARCH=amd64

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

%:
	@:

########################################################################################################################

owner: ## Reset folder owner
	sudo chown --changes -R $$(whoami) ./
	@echo "Success"

check-conflicts: ## Find git conflicts
	@if grep -rn '^<<<\<<<< ' .; then exit 1; fi
	@if grep -rn '^===\====$$' .; then exit 1; fi
	@if grep -rn '^>>>\>>>> ' .; then exit 1; fi
	@echo "All is OK"

check-todos: ## Find TODO's
	@if grep -rn '@TO\DO:' .; then exit 1; fi
	@echo "All is OK"

check-master: ## Check for latest master in current branch
	@git remote update
	@if ! git log --pretty=format:'%H' | grep $$(git log --pretty=format:'%H' -n 1 origin/master) > /dev/null; then exit 1; fi
	@echo "All is OK"

test: ## Run unit tests
	@mkdir -p .profiles
	CGO_ENABLED=1 go test -race -cover -covermode=atomic -coverprofile=.profiles/cover.out.tmp ./...

coverage: test ## Check code coveragem
	@cat .profiles/cover.out.tmp | grep -v "mock_" > .profiles/cover.out
	go tool cover -func=.profiles/cover.out | grep "total:"
ifeq ($(DISABLE_HTML),true)
	gocover-cobertura < .profiles/cover.out > .profiles/cobertura.xml
else
	go tool cover -html=.profiles/cover.out -o .profiles/report.html
endif

lint: ## Lint source code
	@clear
	golangci-lint run --color=always --config=.golangci.yml ./...

build: ## Build app
	mkdir -p bin/$(GOOS)-$(GOARCH)
	CGO_ENABLED=0 go build -ldflags "-X 'main.version=$(BUILD_VERSION)-$(GOOS)-$(GOARCH)'" -o bin/$(GOOS)-$(GOARCH)/go-srv cmd/go-srv/*
	chmod +x bin/$(GOOS)-$(GOARCH)/go-srv
	bin/$(GOOS)-$(GOARCH)/go-srv --version

vendor:
	go mod tidy
	rm -rf vendor && go mod vendor

########################################################################################################################

_docker:
	docker build --build-arg SRC_IMAGE=${SRC_IMAGE} --file .ci/golang.Dockerfile --tag go-srv:golang .ci
	
docker-build: _docker
	docker run --rm -w /src -v "$(PWD)":"/src" go-srv:golang make build

compose: docker-build ## Run app
	GOOS=$(GOOS) GOARCH=$(GOARCH) docker-compose up --build

compose-stop: ## Stop app
	docker-compose down
