.DEFAULT_GOAL := build

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
CONTROLLER_GEN                     := $(TOOLS_DIR)/controller-gen
CONTROLLER_GEN_VERSION             := v0.12.0
REGISTER_GEN                       := $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       := $(TOOLS_DIR)/deepcopy-gen
CODE_GEN_VERSION                   := v0.28.0
REFERENCE_DOCS                     := $(TOOLS_DIR)/genref
REFERENCE_DOCS_VERSION             := latest
TOOLS                              := $(CONTROLLER_GEN) $(REGISTER_GEN) $(DEEPCOPY_GEN) $(REFERENCE_DOCS)
ifeq ($(GOOS), darwin)
SED                                := gsed
else
SED                                := sed
endif
COMMA                              := ,

$(CONTROLLER_GEN):
	@echo Install controller-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

$(REGISTER_GEN):
	@echo Install register-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(CODE_GEN_VERSION)

$(DEEPCOPY_GEN):
	@echo Install deepcopy-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GEN_VERSION)

$(REFERENCE_DOCS):
	@echo Install genref... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/kubernetes-sigs/reference-docs/genref@$(REFERENCE_DOCS_VERSION)

.PHONY: install-tools
install-tools: $(TOOLS) ## Install tools

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

#########
# BUILD #
#########

CMD_DIR        := cmd
CLI_DIR        := $(CMD_DIR)/cli
CLI_BIN        := kyverno-json
CGO_ENABLED    ?= 0
GOOS           ?= $(shell go env GOOS)
ifdef VERSION
LD_FLAGS       := "-s -w -X $(PACKAGE)/pkg/version.BuildVersion=$(VERSION)"
else
LD_FLAGS       := "-s -w"
endif

.PHONY: fmt
fmt: ## Run go fmt
	@echo Go fmt... >&2
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo Go vet... >&2
	@go vet ./...

$(CLI_BIN): fmt vet
	@echo Build cli binary... >&2
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) go build -o ./$(CLI_BIN) -ldflags=$(LD_FLAGS) ./$(CLI_DIR)

build: $(CLI_BIN) ## Build

###########
# CODEGEN #
###########

ORG                         ?= kyverno
PACKAGE                     ?= github.com/$(ORG)/kyverno-json
GOPATH_SHIM                 := ${PWD}/.gopath
PACKAGE_SHIM                := $(GOPATH_SHIM)/src/$(PACKAGE)
INPUT_DIRS                  := $(PACKAGE)/pkg/apis/v1alpha1
CRDS_PATH                   := ${PWD}/config/crds

$(GOPATH_SHIM):
	@echo Create gopath shim... >&2
	@mkdir -p $(GOPATH_SHIM)

.INTERMEDIATE: $(PACKAGE_SHIM)
$(PACKAGE_SHIM): $(GOPATH_SHIM)
	@echo Create package shim... >&2
	@mkdir -p $(GOPATH_SHIM)/src/github.com/$(ORG) && ln -s -f ${PWD} $(PACKAGE_SHIM)

.PHONY: codegen-register
codegen-register: $(PACKAGE_SHIM) $(REGISTER_GEN) ## Generate types registrations
	@echo Generate registration... >&2
	@GOPATH=$(GOPATH_SHIM) $(REGISTER_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--input-dirs=$(INPUT_DIRS)

.PHONY: codegen-deepcopy
codegen-deepcopy: $(PACKAGE_SHIM) $(DEEPCOPY_GEN) ## Generate deep copy functions
	@echo Generate deep copy functions... >&2
	@GOPATH=$(GOPATH_SHIM) $(DEEPCOPY_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--input-dirs=$(INPUT_DIRS) \
		--output-file-base=zz_generated.deepcopy

.PHONY: codegen-crds
codegen-crds: $(CONTROLLER_GEN) ## Generate CRDs
	@echo Generate crds... >&2
	@rm -rf $(CRDS_PATH)
	@$(CONTROLLER_GEN) crd paths=./pkg/apis/... crd:crdVersions=v1 output:dir=$(CRDS_PATH)
	@echo Copy generated CRDs to embed in the CLI... >&2
	@rm -rf pkg/data/crds && mkdir -p pkg/data/crds
	@cp config/crds/* pkg/data/crds

.PHONY: codegen-api-docs-md
codegen-api-docs-md: $(REFERENCE_DOCS) ## Generate markdown API docs
	@echo Generate md api docs... >&2
	@rm -rf ./docs/user/apis/md
	@cd ./docs/user/apis/_config && $(REFERENCE_DOCS) -c config.yaml -f markdown -o ../md

.PHONY: codegen-api-docs-html
codegen-api-docs-html: $(REFERENCE_DOCS) ## Generate html API docs
	@echo Generate html api docs... >&2
	@rm -rf ./docs/user/apis/html
	@cd ./docs/user/apis/_config && $(REFERENCE_DOCS) -c config.yaml -f html -o ../html

.PHONY: codegen-api-docs
codegen-api-docs: codegen-api-docs-md codegen-api-docs-html ## Generate API docs

.PHONY: codegen-cli-docs
codegen-cli-docs: $(CLI_BIN) ## Generate CLI docs
	@echo Generate cli docs... >&2
	@rm -rf docs/user/commands && mkdir -p docs/user/commands
	@./kyverno-json docs -o docs/user/commands --autogenTag=false

.PHONY: codegen-jp-docs
codegen-jp-docs: ## Generate JP docs
	@echo Generate jp docs... >&2
	@rm -rf docs/user/jp && mkdir -p docs/user/jp
	@go run ./hack/docs/jp/main.go > docs/user/jp/functions.md

.PHONY: codegen-docs
codegen-docs: codegen-api-docs-md codegen-cli-docs codegen-jp-docs ## Generate docs

.PHONY: codegen-mkdocs
codegen-mkdocs: codegen-docs ## Generate mkdocs website
	@echo Generate mkdocs website... >&2
	@pip3 install mkdocs
	@pip3 install --upgrade pip
	@pip3 install -U mkdocs-material mkdocs-redirects mkdocs-minify-plugin mkdocs-include-markdown-plugin lunr mkdocs-rss-plugin
	@rm -rf ./website/docs/apis && mkdir -p ./website/docs/apis && cp docs/user/apis/md/* ./website/docs/apis
	@rm -rf ./website/docs/commands && mkdir -p ./website/docs/commands && cp docs/user/commands/* ./website/docs/commands
	@rm -rf ./website/docs/jp && mkdir -p ./website/docs/jp && cp docs/user/jp/* ./website/docs/jp
	@mkdocs build -f ./website/mkdocs.yaml

.PHONY: codegen-all
codegen-all: codegen-crds codegen-deepcopy codegen-register codegen-docs codegen-mkdocs ## Rebuild all generated code and docs

.PHONY: verify-codegen
verify-codegen: codegen-all ## Verify all generated code and docs are up to date
	@echo Checking codegen is up to date... >&2
	@git --no-pager diff -- .
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen-all".' >&2
	@echo 'To correct this, locally run "make codegen-all", commit the changes, and re-run tests.' >&2
	@git diff --quiet --exit-code -- .

########
# TEST #
########

.PHONY: tests
tests: $(CLI_BIN) ## Run tests
	@echo Running tests... >&2
	@go test ./...

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
