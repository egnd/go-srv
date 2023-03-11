#!/usr/bin/env sh

GOLANGCI_VER=v1.51.2
MOCKERY_VER=v2.22.1
COBERTURA_VER=v1.2.0

export GOPROXY=https://proxy.golang.org,direct

# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_VER}
wget https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh -O - | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_VER}
golangci-lint --version

go install github.com/vektra/mockery/v2@${MOCKERY_VER}
go install github.com/boumenot/gocover-cobertura@${COBERTURA_VER}
