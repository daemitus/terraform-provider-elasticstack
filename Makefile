.DEFAULT_GOAL = help
SHELL := /bin/bash

VERSION ?= 0.11.5-pre

NAME = elasticstack
BINARY = terraform-provider-${NAME}
MARCH = "$$(go env GOOS)_$$(go env GOARCH)"

ACCTEST_PARALLELISM ?= 10
ACCTEST_TIMEOUT = 120m
ACCTEST_COUNT = 1
TEST ?= ./...
SWAGGER_VERSION ?= 8.7

ELASTICSEARCH_NAME ?= terraform-elasticstack-es
ELASTICSEARCH_ENDPOINTS ?= http://$(ELASTICSEARCH_NAME):9200
ELASTICSEARCH_USERNAME ?= elastic
ELASTICSEARCH_PASSWORD ?= password
ELASTICSEARCH_MEM ?= 1024m

KIBANA_NAME ?= terraform-elasticstack-kb
KIBANA_ENDPOINT ?= http://$(KIBANA_NAME):5601
KIBANA_SYSTEM_USERNAME ?= kibana_system
KIBANA_SYSTEM_PASSWORD ?= password
KIBANA_API_KEY_NAME ?= kibana-api-key

SOURCE_LOCATION ?= $(shell pwd)

export GOBIN = $(shell pwd)/bin


$(GOBIN): ## create bin/ in the current directory
	mkdir -p $(GOBIN)

## Downloads all the Golang dependencies.
vendor:
	@ go mod download

.PHONY: build-ci
build-ci: ## build the terraform provider
	go build -o ${BINARY}

.PHONY: build
build: lint build-ci ## build the terraform provider

.PHONY: testacc
testacc: ## Run acceptance tests
	TF_ACC=1 go test -v ./... -count $(ACCTEST_COUNT) -parallel $(ACCTEST_PARALLELISM) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT)

.PHONY: test
test: ## Run unit tests
	go test -v $(TEST) $(TESTARGS) -timeout=5m -parallel=4

# To run specific test (e.g. TestAccResourceActionConnector) execute `make docker-testacc TESTARGS='-run ^TestAccResourceActionConnector$$'`
# To enable tracing (or debugging), execute `make docker-testacc TF_LOG=TRACE`
.PHONY: docker-testacc
docker-testacc: docker-elasticsearch docker-kibana ## Run acceptance tests in the docker container
	@ docker run --rm \
		-e ELASTICSEARCH_ENDPOINTS="$(ELASTICSEARCH_ENDPOINTS)" \
		-e KIBANA_ENDPOINT="$(KIBANA_ENDPOINT)" \
		-e ELASTICSEARCH_USERNAME="$(ELASTICSEARCH_USERNAME)" \
		-e ELASTICSEARCH_PASSWORD="$(ELASTICSEARCH_PASSWORD)" \
		-e TF_LOG="$(TF_LOG)" \
		--network $(ELASTICSEARCH_NETWORK) \
		-w "/provider" \
		-v "$(SOURCE_LOCATION):/provider" \
		golang:$(GOVERSION) make testacc TESTARGS="$(TESTARGS)"


.PHONY: docs-generate
docs-generate: tools ## Generate documentation for the provider
	@ $(GOBIN)/tfplugindocs

.PHONY: gen
gen: docs-generate ## Generate the code and documentation
	@ go generate ./...


.PHONY: tools
tools: $(GOBIN) ## Install useful tools for linting, docs generation and development
	@ cd tools && go install github.com/client9/misspell/cmd/misspell
	@ cd tools && go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	@ cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	@ cd tools && go install github.com/goreleaser/goreleaser
	@ cd tools && go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen

.PHONY: misspell
misspell:
	@ $(GOBIN)/misspell -error -source go ./internal/
	@ $(GOBIN)/misspell -error -source text ./templates/


.PHONY: golangci-lint
golangci-lint:
	@ $(GOBIN)/golangci-lint run --max-same-issues=0 --timeout=300s $(GOLANGCIFLAGS) ./internal/...


.PHONY: lint
lint: setup misspell golangci-lint check-fmt check-docs ## Run lints to check the spelling and common go patterns

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	terraform fmt --recursive

.PHONY:check-fmt
check-fmt: fmt ## Check if code is formatted
	@if [ "`git status --porcelain `" ]; then \
	  echo "Unformatted files were detected. Please run 'make fmt' to format code, and commit the changes" && echo `git status --porcelain docs/` && exit 1; \
	fi

.PHONY: check-docs
check-docs: docs-generate  ## Check uncommitted changes on docs
	@if [ "`git status --porcelain docs/`" ]; then \
	  echo "Uncommitted changes were detected in the docs folder. Please run 'make docs-generate' to autogenerate the docs, and commit the changes" && echo `git status --porcelain docs/` && exit 1; \
	fi


.PHONY: setup
setup: tools ## Setup the dev environment


.PHONY: release-snapshot
release-snapshot: tools ## Make local-only test release to see if it works using "release" command
	@ $(GOBIN)/goreleaser release --snapshot --clean


.PHONY: release-no-publish
release-no-publish: tools check-sign-release ## Make a release without publishing artifacts
	@ $(GOBIN)/goreleaser release --skip=publish,announce,validate  --parallelism=2


.PHONY: release
release: tools check-sign-release check-publish-release ## Build, sign, and upload your release
	@ $(GOBIN)/goreleaser release --clean  --parallelism=4


.PHONY: check-sign-release
check-sign-release:
ifndef GPG_FINGERPRINT_SECRET
	$(error GPG_FINGERPRINT_SECRET is undefined, but required for signing the release)
endif


.PHONY: check-publish-release
check-publish-release:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined, but required to make build and upload the released artifacts)
endif


.PHONY: release-notes
release-notes: ## greps UNRELEASED notes from the CHANGELOG
	@ awk '/## \[Unreleased\]/{flag=1;next}/## \[.*\] - /{flag=0}flag' CHANGELOG.md


.PHONY: help
help: ## this help
	@ awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m\t%s\n", $$1, $$2 }' $(MAKEFILE_LIST) | column -s$$'\t' -t

.PHONY: generate-alerting-client
generate-alerting-client: ## generate Kibana alerting client
	@ docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli:v7.0.1 generate \
		-i /local/generated/alerting/bundled.yaml \
		--skip-validate-spec \
		--git-repo-id terraform-provider-elasticstack \
		--git-user-id elastic \
		-p isGoSubmodule=true \
		-p packageName=alerting \
		-p generateInterfaces=true \
		-g go \
		-o /local/generated/alerting
	@ rm -rf generated/alerting/go.mod generated/alerting/go.sum generated/alerting/test
	@ go fmt ./generated/alerting/...

.PHONY: generate-data-views-client
generate-data-views-client: ## generate Kibana data-views client
	@ docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli:v7.0.1 generate \
		-i /local/generated/data_views/bundled.yaml \
		--skip-validate-spec \
		--git-repo-id terraform-provider-elasticstack \
		--git-user-id elastic \
		-p isGoSubmodule=true \
		-p packageName=data_views \
		-p generateInterfaces=true \
		-g go \
		-o /local/generated/data_views
	@ rm -rf generated/data_views/go.mod generated/data_views/go.sum generated/data_views/test
	@ go fmt ./generated/data_views/...

.PHONY: generate-connectors-client
generate-connectors-client: tools ## generate Kibana connectors client
	@ cd tools && go generate
	@ go fmt ./generated/connectors/...

.PHONY: generate-slo-client
generate-slo-client: tools ## generate Kibana slo client
	@ rm -rf generated/slo
	@ docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli:v7.0.1 generate \
		-i /local/generated/slo-spec.yml \
		--git-repo-id terraform-provider-elasticstack \
		--git-user-id elastic \
		-p isGoSubmodule=true \
		-p packageName=slo \
		-p generateInterfaces=true \
		-p useOneOfDiscriminatorLookup=true \
		-g go \
		-o /local/generated/slo \
		 --type-mappings=float32=float64
	@ rm -rf generated/slo/go.mod generated/slo/go.sum generated/slo/test
	@ go fmt ./generated/...

.PHONY: generate-clients
generate-clients: generate-alerting-client generate-slo-client generate-data-views-client generate-connectors-client ## generate all clients
