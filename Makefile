PROJ_DIR=$(shell pwd)


run:
	go run cmd/main.go


### Linting ###

check-golangci-lint:
	$(eval GOLANGCI_LINT_VERSION=$(shell curl -s https://api.github.com/repos/golangci/golangci-lint/releases/latest | jq -r '.name'))
	@./bin/golangci-lint --version | grep -qF "$(GOLANGCI_LINT_VERSION:1)" || { \
		echo "golangci-lint not found or version mismatch. Installing..."; \
		$(MAKE) install-golangci-lint; \
	}

install-golangci-lint: ## install golang lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

lint-check: check-golangci-lint
	$(PROJ_DIR)/bin/golangci-lint run --config $(PROJ_DIR)/.golangci.yaml


### Unit Tests ###

ut:
	go test -parallel=1 -race ./...


### SQL ###

DB_CODE_DIR=$(PROJ_DIR)/internal/database/querier
DB_SPEC_DIR=$(PROJ_DIR)/schemas/sql

gen-sql: gen-clean-sql
	echo $(PROJ_DIR)
	docker run --rm \
		-v $(DB_CODE_DIR):/src/querier \
		-v $(DB_SPEC_DIR):/src/sql \
		-w /src/sql \
		sqlc/sqlc generate

gen-clean-sql:
	@cd $(DB_SPEC_DIR) && rm -f *gen.go


### OpenAPI ###

API_CODE_DIR=$(PROJ_DIR)/internal/api
API_SPEC_DIR=$(PROJ_DIR)/schemas/api

check-oapi-codegen:
	$(eval OAPI_GEN_VERSION=$(shell curl -s https://api.github.com/repos/oapi-codegen/oapi-codegen/releases/latest | jq -r '.tag_name'))
	@$(PROJ_DIR)/bin/oapi-codegen --version | grep -qF "$(OAPI_GEN_VERSION)" || { \
		echo "oapi-codegen not found or version mismatch (expected $(OAPI_GEN_VERSION)). Installing..."; \
		$(MAKE) install-oapi-codegen; \
	}

install-oapi-codegen:
	@GOBIN=$(PROJ_DIR)/bin go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

oapi-codegen: check-oapi-codegen
	$(PROJ_DIR)/bin/oapi-codegen \
		--config=$(API_SPEC_DIR)/oapi-service.yaml \
		--package=v1 \
		$(API_SPEC_DIR)/v1.yaml \
		> $(API_CODE_DIR)/v1/spec.gen.go


### Build and Deploy ###

codegen:
	make oapi-codegen
	make gen-sql
	go mod tidy

codegen-check:
	git restore :/ && git clean -d -f
	make codegen
	git diff --exit-code

build-docker: codegen
	docker build -f ./Dockerfile . -t go-template
