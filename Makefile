.PHONY: test test-checker test-goroot docs ci ci-tidy ci-tests ci-linter cover install gen-checker

GOPATH_DIR=`go env GOPATH`

export GO111MODULE := on

%:      # stubs to get makefile param for `test-checker` command
	@:	# see: https://stackoverflow.com/a/6273809/433041

test:
	go test -v -count=1 ./...

test-checker:
	go test -v -count=1 -run=/$(filter-out $@,$(MAKECMDGOALS)) ./...

test-goroot:
	go run cmd/gocritic/main.go check-project -enable=$(filter-out $@,$(MAKECMDGOALS)) ${GOROOT}/src

docs:
	cd ./cmd/makedocs && go run main.go

ci:
	@if [ "$(TEST_SUITE)" = "linter" ]; then make ci-linter; else make ci-tidy; make ci-tests; fi

ci-tidy:
	go mod tidy
	# If you are testing Go 1.11: https://github.com/golang/go/issues/27868#issuecomment-431413621
	go list all > /dev/null
	git diff --exit-code --quiet || (echo "Please run 'go mod tidy' to clean up the 'go.mod' and 'go.sum' files."; false)

ci-tests:
	go test -v -race -count=1 -coverprofile=coverage.out ./...

ci-linter:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.30.0
	@$(GOPATH_DIR)/bin/golangci-lint run -v
	go install github.com/quasilyte/go-consistent
	@$(GOPATH_DIR)/bin/go-consistent ./...
	go build -o gocritic ./cmd/gocritic
	./gocritic check -enableAll -disable=ioutilDeprecated ./...

cover:
	go install github.com/mattn/goveralls
	goveralls -package github.com/go-critic/go-critic/checkers -coverprofile=coverage.out -service travis-ci -repotoken ${COVERALLS_TOKEN}

gocritic:
	go build -o gocritic ./cmd/gocritic

install:
	go install ./cmd/gocritic

TEST_DIR=checkers/testdata/$(CHECKER_NAME)
ifeq ($(CHECKER_TYPE),)
CHECKER_TYPE_NAME=Expr
CHECKER_TYPE=expr ast.Expr
endif
ifeq ($(CHECKER_CATEGORY),)
CHECKER_CATEGORY=diagnostic
endif

define CHECKER_BODY
package checkers

import (
	"go/ast"

	"github.com/go-critic/go-critic/checkers/internal/astwalk"
	"github.com/go-critic/go-critic/framework/linter"
)

func init() {
	var info linter.CheckerInfo
	info.Name = "$(CHECKER_NAME)"
	info.Tags = []string{"$(CHECKER_CATEGORY)", "experimental"}
	info.Summary = "Write your summary here."
	info.Before = ``
	info.After = ``

	collection.AddChecker(&info, func(ctx *linter.CheckerContext) linter.FileWalker {
		return astwalk.WalkerFor$(CHECKER_TYPE_NAME)(&$(CHECKER_NAME)Checker{ctx: ctx})
	})
}

type $(CHECKER_NAME)Checker struct {
	astwalk.WalkHandler
	ctx *linter.CheckerContext
}

func (c *$(CHECKER_NAME)Checker) Visit$(CHECKER_TYPE_NAME)($(CHECKER_TYPE)) {
	// your code goes here
}
endef

export CHECKER_BODY

gen-checker:
	mkdir $(TEST_DIR)
	echo "package checker_test" > $(TEST_DIR)/negative_tests.go
	echo "package checker_test" > $(TEST_DIR)/positive_tests.go
	echo "$$CHECKER_BODY" > checkers/$(CHECKER_NAME)_checker.go
