WRAPPER_SUBMODULES := fasthttp gin gorestful
PWD := $(shell pwd)
DIR :=
COMPAT_COREGO_TEST := $(PWD)/compat/corego/corego_compat_test.go
CORE_GO_REF ?= master

# must ensure your go version >= 1.16
.PHONY: install
install:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install golang.org/x/tools/cmd/goimports@latest

# usage
# you muse run `make install` to install necessary tools
# make mock dir=path/to/mock
.PHONY: mock
mock:
	@for file in `find . -type d \( -path ./.git -o -path ./.github \) -prune -o -name '*.go' -print | xargs grep --files-with-matches -e '//go:generate mockgen'`; do \
		go generate $$file; \
	done


.PHONY: tidy
tidy:
	go mod tidy && go fmt ./...
	@$(foreach var,$(WRAPPER_SUBMODULES),cd $(PWD)/wrapper/$(var) && go mod tidy && go fmt ./...;)

.PHONY: test
test:
	go test ./...
	$(foreach var,$(WRAPPER_SUBMODULES),cd $(PWD)/wrapper/$(var) && go test ./...;)

.PHONY: test-compat-corego
test-compat-corego:
	@droplet_dir=$$(cd "$(PWD)" && pwd); \
	tmpdir=$$(mktemp -d); \
	trap 'rm -rf "$$tmpdir"' EXIT; \
	cp "$(COMPAT_COREGO_TEST)" "$$tmpdir/corego_compat_test.go"; \
	printf '%s\n' \
		'module corego_compat_check' \
		'' \
		'go 1.24.3' \
		'' \
		"replace github.com/shiningrush/droplet => $$droplet_dir" \
		> "$$tmpdir/go.mod"; \
	if [ -n "$(CORE_GO_DIR)" ]; then \
		core_go_dir=$$(cd "$(CORE_GO_DIR)" && pwd); \
		printf '%s\n' "replace github.com/dev-ofa/core-go => $$core_go_dir" >> "$$tmpdir/go.mod"; \
		cd "$$tmpdir" && go mod tidy && go test -tags compat_corego ./...; \
	else \
		cd "$$tmpdir" && go get github.com/dev-ofa/core-go@$(CORE_GO_REF) && go mod tidy && go test -tags compat_corego ./...; \
	fi
