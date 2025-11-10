-include Makefile.mk

GOFLAGS=-mod=vendor

PKG := `go list ${GOFLAGS} -f {{.Dir}} ./...`

ifeq ($(RACE),1)
	GOFLAGS+=-race
endif

LINT_VERSION := v2.4.0

MAIN := ${NAME}/cmd/${NAME}


.PHONY: *

init:
	@cp -n Makefile.mk.dist Makefile.mk
	@cp -n cfg/local.toml.dist cfg/local.toml

show-env:
	@echo "NAME=$(NAME)"
	@echo "GOFLAGS=$(GOFLAGS)"

tools:
	@go install github.com/vmkteam/mfd-generator@latest
	@go install github.com/vmkteam/pgmigrator@latest
	@go install github.com/vmkteam/colgen/cmd/colgen@latest
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${LINT_VERSION}

fmt:
	@golangci-lint fmt

lint:
	@golangci-lint version
	@golangci-lint config verify
	@golangci-lint run

build:
	@CGO_ENABLED=0 go build $(GOFLAGS) -o ${NAME} $(MAIN)

run:
	@echo "Compiling"
	@go run $(GOFLAGS) $(MAIN) -config=cfg/local.toml -dev


test-short:
	@go test $(GOFLAGS) -v -test.short -test.run="Test[^D][^B]" -coverprofile=coverage.txt -covermode count $(PKG)

mod:
	@go mod tidy
	@go mod vendor
	@git add vendor




--check-ns:
ifeq ($(NS),"NONE")
	$(error "You need to set NS variable before run this command. For example: NS=common make $(MAKECMDGOALS) or: make $(MAKECMDGOALS) NS=common")
endif
