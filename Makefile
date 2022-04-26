GOBIN=$(GOPATH)/bin
GORUN=go run
GOBUILD=go build
GOTEST=go test

include .env
export

PROJECT_DIR=$(shell pwd)
EXAMPLES_DIR=$(PROJECT_DIR)/specs
CONTRACTOR=$(PROJECT_DIR)/cmd/contractor/main.go

clean:
	find $(PROJECT_DIR) \( -path "*/mocks" -and ! -path "*/vendor/*" \) -exec rm -rf {} +
	rm -rf $(PROJECT_DIR)/vendor

install:
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

update: clean install generate

generate-graylog-client:
	rm -rf $(PROJECT_DIR)/tester/graylog/client
	mkdir -p $(PROJECT_DIR)/tester/graylog/client
	swagger generate client \
		-m client/models \
		-f $(EXAMPLES_DIR)/graylog.json \
		-t $(PROJECT_DIR)/tester/graylog

run-tester:
	$(GORUN) $(CONTRACTOR) -m test \
							-c $(PROJECT_DIR)/config.yml \
							-s $(EXAMPLES_DIR)/oa2.yml \
							-u $(URL_BASE) \
							-f $(SPEC_TYPE) \
							-v "VAR_AUTHORIZATION:$(VAR_AUTHORIZATION)"

run-tester2:
	$(GORUN) $(CONTRACTOR) -m test2 \
							-c $(PROJECT_DIR)/config.yml

run-mocker:
	$(GORUN) $(CONTRACTOR) -m mock \
							-c $(PROJECT_DIR)/config.yml \
							-s $(EXAMPLES_DIR)/oa2.yml \
							-u $(URL_BASE) \
							-f $(SPEC_TYPE) \
							-v "VAR_AUTHORIZATION:$(VAR_AUTHORIZATION)"

test-graylog:
	curl -i -u 1geh46b37k3mib77lj8lt401mkah24d1alomosish521tsfqsuaf:token -H 'Accept: application/json' -X GET 'https://graylog.setpartnerstv.ru:443/api/search/universal/relative?query=app%3Aautopayment-processor&range=300&limit=100&batch_size=500&fields=correlation_id,message' > r.json
#sys.qa
#5xVa6tsM2dvstFU
#1geh46b37k3mib77lj8lt401mkah24d1alomosish521tsfqsuaf

test-c2p:
	curl -i -X POST -H 'Content-Type: application/json' -H 'Authorization: eea00e6ff92930a4b876c8e6c975c62d' 'https://vp.voicepassword.ru/api/voice-password/send/' -d '{"number":"+79255308012","flashcall": {"code": 1234}}'
