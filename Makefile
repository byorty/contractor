GOBIN=$(GOPATH)/bin
GORUN=go run
GOBUILD=go build
GOTEST=go test

include .env

PROJECT_DIR=$(shell pwd)
EXAMPLES_DIR=$(PROJECT_DIR)/examples
CONTRACTOR=$(PROJECT_DIR)/cmd/contractor/main.go

install:
	find $(shell pwd) \( -regex '.*mock_.*' -and ! -path "*/vendor/*" \) -exec rm {} \;
	rm -rf $(PROJECT_DIR)/vendor
	GOPROXY=direct GOSUMDB=off go mod vendor

test:
	$(GOTEST) -v ./...

generate:
	go generate ./...

run-tester:
	@$(GORUN) $(CONTRACTOR) -m test \
							-s $(EXAMPLES_DIR)/swagger.yml \
							-b $(URL_BASE) \
							-t $(SPEC_TYPE) \
							-v "VAR_AUTHORIZATION:$(VAR_AUTHORIZATION)"

run-mocker:
	$(GORUN) $(CONTRACTOR) -m mock \
							-s $(EXAMPLES_DIR)/open_api_v3.yml \
							-b $(URL_BASE) \
							-t $(SPEC_TYPE)