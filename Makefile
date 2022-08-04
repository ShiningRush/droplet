WRAPPER_SUBMODULES := fasthttp gin gorestful
PWD := $(shell pwd)

.PHONY: tidy
tidy:
	go mod tidy && go fmt ./...
	$(foreach var,$(WRAPPER_SUBMODULES),cd $(PWD)/wrapper/$(var) && go mod tidy && go fmt ./...;)

.PHONY: test
test:
	go test ./...
	$(foreach var,$(WRAPPER_SUBMODULES),cd $(PWD)/wrapper/$(var) && go test ./...;)