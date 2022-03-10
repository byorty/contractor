GOBIN=$(GOPATH)/bin
GORUN=go run
GOBUILD=go build
GOTEST=go test

include .env

PROJECT_DIR=$(shell pwd)
EXAMPLES_DIR=$(PROJECT_DIR)/specs
CONTRACTOR=$(PROJECT_DIR)/cmd/contractor/main.go

install:
	find $(shell pwd) \( -regex '.*mock_.*' -and ! -path "*/vendor/*" \) -exec rm {} \;
	rm -rf $(PROJECT_DIR)/vendor
	GOPROXY=direct GOSUMDB=off go mod vendor

build:
	go install ./...
	echo $(GOPATH)
	echo $(GOBIN)
	ls -la $(GOBIN)

test:
	$(GOTEST) -v ./...

generate:
	go generate ./...

run-tester:
	$(GORUN) $(CONTRACTOR) -m test \
							-s $(EXAMPLES_DIR)/oa2.yml \
							-b $(URL_BASE) \
							-f $(SPEC_TYPE) \
							-v "VAR_AUTHORIZATION:$(VAR_AUTHORIZATION)"

run-mocker:
	$(GORUN) $(CONTRACTOR) -m mock \
							-s $(EXAMPLES_DIR)/oa2.yml \
							-b $(URL_BASE) \
							-f $(SPEC_TYPE)

send-request:
	curl -H "Content-Type: application/json" http://localhost:8181/v1/news/11401

check-mocker:
	$(GORUN) $(CONTRACTOR) -m test \
							-s $(EXAMPLES_DIR)/oa2.yml \
							-b "http://localhost:8181" \
							-t "get_news_by_id" \
							-f $(SPEC_TYPE)