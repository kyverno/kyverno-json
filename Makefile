.DEFAULT_GOAL := build

##########
# CONFIG #
##########

ORG                                ?= kyverno
PACKAGE                            ?= github.com/$(ORG)/kyverno-json
KIND_IMAGE                         ?= kindest/node:v1.28.0
KIND_NAME                          ?= kind
GIT_SHA                            := $(shell git rev-parse HEAD)
GOOS                               ?= $(shell go env GOOS)
GOARCH                             ?= $(shell go env GOARCH)
REGISTRY                           ?= ghcr.io
REPO                               ?= kyverno-json
LOCAL_PLATFORM                     := linux/$(GOARCH)
KO_REGISTRY                        := ko.local
KO_PLATFORMS                       := all
KO_TAGS                            := $(GIT_SHA)
KO_CACHE                           ?= /tmp/ko-cache
CLI_BIN                            := kyverno-json

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
CLIENT_GEN                         := $(TOOLS_DIR)/client-gen
LISTER_GEN                         := $(TOOLS_DIR)/lister-gen
INFORMER_GEN                       := $(TOOLS_DIR)/informer-gen
REGISTER_GEN                       := $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       := $(TOOLS_DIR)/deepcopy-gen
CODE_GEN_VERSION                   := v0.28.0
CONTROLLER_GEN                     := $(TOOLS_DIR)/controller-gen
CONTROLLER_GEN_VERSION             := v0.12.0
REFERENCE_DOCS                     := $(TOOLS_DIR)/genref
REFERENCE_DOCS_VERSION             := latest
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.20.0
HELM                               := $(TOOLS_DIR)/helm
HELM_VERSION                       := v3.10.1
KO                                 := $(TOOLS_DIR)/ko
KO_VERSION                         := v0.14.1
TOOLS                              := $(CLIENT_GEN) $(LISTER_GEN) $(INFORMER_GEN) $(REGISTER_GEN) $(DEEPCOPY_GEN) $(CONTROLLER_GEN) $(REFERENCE_DOCS) $(KIND) $(HELM) $(KO)
PIP                                ?= "pip"
ifeq ($(GOOS), darwin)
SED                                := gsed
else
SED                                := sed
endif
COMMA                              := ,

$(CLIENT_GEN):
	@echo Install client-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/client-gen@$(CODE_GEN_VERSION)

$(LISTER_GEN):
	@echo Install lister-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/lister-gen@$(CODE_GEN_VERSION)

$(INFORMER_GEN):
	@echo Install informer-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/informer-gen@$(CODE_GEN_VERSION)

$(REGISTER_GEN):
	@echo Install register-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(CODE_GEN_VERSION)

$(DEEPCOPY_GEN):
	@echo Install deepcopy-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GEN_VERSION)

$(CONTROLLER_GEN):
	@echo Install controller-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

$(REFERENCE_DOCS):
	@echo Install genref... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/kubernetes-sigs/reference-docs/genref@$(REFERENCE_DOCS_VERSION)

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)

$(HELM):
	@echo Install helm... >&2
	@GOBIN=$(TOOLS_DIR) go install helm.sh/helm/v3/cmd/helm@$(HELM_VERSION)

$(KO):
	@echo Install ko... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/google/ko@$(KO_VERSION)

.PHONY: install-tools
install-tools: $(TOOLS) ## Install tools

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

#########
# BUILD #
#########

CGO_ENABLED    ?= 0
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

$(CLI_BIN): fmt vet build-wasm codegen-crds codegen-deepcopy codegen-register codegen-client codegen-playground
	@echo Build cli binary... >&2
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) go build -o ./$(CLI_BIN) -ldflags=$(LD_FLAGS) .

.PHONY: build
build: $(CLI_BIN) ## Build

.PHONY: build-wasm
build-wasm: fmt vet codegen-crds codegen-deepcopy codegen-register codegen-client ## Build the wasm binary
	@echo Build wasm module... >&2
	@GOOS=js GOARCH=wasm go build -o ./website/playground/assets/main.wasm -ldflags=$(LD_FLAGS) ./cmd/wasm/main.go

.PHONY: serve-playground
serve-playground: $(CLI_BIN) ## Serve playground
	@echo Serve playground... >&2
	@./$(CLI_BIN) playground

.PHONY: ko-build
ko-build: $(KO) ## Build image (with ko)
	@echo Build image with ko... >&2
	@LDFLAGS=$(LD_FLAGS) KOCACHE=$(KO_CACHE) KO_DOCKER_REPO=$(KO_REGISTRY) \
		$(KO) build ./$(CLI_DIR) --preserve-import-paths --tags=$(KO_TAGS) --platform=$(LOCAL_PLATFORM)

###########
# CODEGEN #
###########

GOPATH_SHIM                 := ${PWD}/.gopath
PACKAGE_SHIM                := $(GOPATH_SHIM)/src/$(PACKAGE)
INPUT_DIRS                  := $(PACKAGE)/pkg/apis/v1alpha1
CRDS_PATH                   := ${PWD}/.crds
INPUT_DIRS                  := $(PACKAGE)/pkg/apis/v1alpha1
OUT_PACKAGE                 := $(PACKAGE)/pkg/client
CLIENTSET_PACKAGE           := $(OUT_PACKAGE)/clientset
LISTERS_PACKAGE             := $(OUT_PACKAGE)/listers
INFORMERS_PACKAGE           := $(OUT_PACKAGE)/informers

$(GOPATH_SHIM):
	@echo Create gopath shim... >&2
	@mkdir -p $(GOPATH_SHIM)

.INTERMEDIATE: $(PACKAGE_SHIM)
$(PACKAGE_SHIM): $(GOPATH_SHIM)
	@echo Create package shim... >&2
	@mkdir -p $(GOPATH_SHIM)/src/github.com/$(ORG) && ln -s -f ${PWD} $(PACKAGE_SHIM)

.PHONY: codegen-client
codegen-client: $(PACKAGE_SHIM) $(CLIENT_GEN) $(LISTER_GEN) $(INFORMER_GEN) ## Generate client, informers and listers
	@echo Generate clientset... >&2
	@GOPATH=$(GOPATH_SHIM) $(CLIENT_GEN) \
		--go-header-file ./.hack/boilerplate.go.txt \
		--clientset-name versioned \
		--output-package $(CLIENTSET_PACKAGE) \
		--input-base "" \
		--input $(INPUT_DIRS)
	@echo Generate listers... >&2
	@GOPATH=$(GOPATH_SHIM) $(LISTER_GEN) \
		--go-header-file ./.hack/boilerplate.go.txt \
		--output-package $(LISTERS_PACKAGE) \
		--input-dirs $(INPUT_DIRS)
	@echo Generate informers... >&2
	@GOPATH=$(GOPATH_SHIM) $(INFORMER_GEN) \
		--go-header-file ./.hack/boilerplate.go.txt \
		--output-package $(INFORMERS_PACKAGE) \
		--input-dirs $(INPUT_DIRS) \
		--versioned-clientset-package $(CLIENTSET_PACKAGE)/versioned \
		--listers-package $(LISTERS_PACKAGE)

.PHONY: codegen-register
codegen-register: $(PACKAGE_SHIM) $(REGISTER_GEN) ## Generate types registrations
	@echo Generate registration... >&2
	@GOPATH=$(GOPATH_SHIM) $(REGISTER_GEN) \
		--go-header-file=./.hack/boilerplate.go.txt \
		--input-dirs=$(INPUT_DIRS)

.PHONY: codegen-deepcopy
codegen-deepcopy: $(PACKAGE_SHIM) $(DEEPCOPY_GEN) ## Generate deep copy functions
	@echo Generate deep copy functions... >&2
	@GOPATH=$(GOPATH_SHIM) $(DEEPCOPY_GEN) \
		--go-header-file=./.hack/boilerplate.go.txt \
		--input-dirs=$(INPUT_DIRS) \
		--output-file-base=zz_generated.deepcopy

.PHONY: codegen-crds
codegen-crds: $(CONTROLLER_GEN) ## Generate CRDs
	@echo Generate crds... >&2
	@rm -rf $(CRDS_PATH)
	@$(CONTROLLER_GEN) crd paths=./pkg/apis/... crd:crdVersions=v1 output:dir=$(CRDS_PATH)
	@echo Copy generated CRDs to embed in the CLI... >&2
	@rm -rf pkg/data/crds && mkdir -p pkg/data/crds
	@cp $(CRDS_PATH)/* pkg/data/crds

.PHONY: codegen-api-docs
codegen-api-docs: $(REFERENCE_DOCS) ## Generate API docs
	@echo Generate md api docs... >&2
	@rm -rf ./website/docs/apis
	@cd ./website/apis && $(REFERENCE_DOCS) -c config.yaml -f markdown -o ../docs/apis

.PHONY: codegen-cli-docs
codegen-cli-docs: $(CLI_BIN) ## Generate CLI docs
	@echo Generate cli docs... >&2
	@rm -rf docs/user/commands && mkdir -p docs/user/commands
	@./kyverno-json docs -o docs/user/commands --autogenTag=false

.PHONY: codegen-jp-docs
codegen-jp-docs: ## Generate JP docs
	@echo Generate jp docs... >&2
	@rm -rf docs/user/jp && mkdir -p docs/user/jp
	@go run ./.hack/docs/jp/main.go > docs/user/jp/functions.md

.PHONY: codegen-catalog
codegen-catalog: ## Generate policy catalog
	@echo Generate policy catalog... >&2
	@go run ./.hack/docs/catalog/main.go

.PHONY: codegen-docs
codegen-docs: codegen-api-docs codegen-cli-docs codegen-jp-docs codegen-catalog ## Generate docs

.PHONY: codegen-mkdocs
codegen-mkdocs: codegen-docs ## Generate mkdocs website
	@echo Generate mkdocs website... >&2
	@pip install mkdocs
	@pip install --upgrade pip
	@pip install -U mkdocs-material mkdocs-redirects mkdocs-minify-plugin mkdocs-include-markdown-plugin lunr mkdocs-rss-plugin
	@rm -rf ./website/docs/commands && mkdir -p ./website/docs/commands && cp docs/user/commands/* ./website/docs/commands
	@rm -rf ./website/docs/jp && mkdir -p ./website/docs/jp && cp docs/user/jp/* ./website/docs/jp
	@mkdocs build -f ./website/mkdocs.yaml

.PHONY: codegen-schemas-openapi
codegen-schemas-openapi: $(KIND) $(HELM) ## Generate openapi schemas (v2 and v3)
	@echo Generate openapi schema... >&2
	@rm -rf ./.schemas
	@mkdir -p ./.schemas/openapi/v2
	@mkdir -p ./.schemas/openapi/v3/apis/json.kyverno.io
	@$(KIND) create cluster --name schema --image $(KIND_IMAGE)
	@kubectl create -f $(CRDS_PATH)
	@sleep 15
	@kubectl get --raw /openapi/v2 > ./.schemas/openapi/v2/schema.json
	@kubectl get --raw /openapi/v3/apis/json.kyverno.io/v1alpha1 > ./.schemas/openapi/v3/apis/json.kyverno.io/v1alpha1.json
	@$(KIND) delete cluster --name schema

.PHONY: codegen-schemas-json
codegen-schemas-json: codegen-schemas-openapi ## Generate json schemas
	@$(PIP) install openapi2jsonschema
	@rm -rf ./.schemas/json
	@openapi2jsonschema ./.schemas/openapi/v2/schema.json --kubernetes --stand-alone --expanded -o ./.schemas/json

.PHONY: codegen-schemas
codegen-schemas: codegen-schemas-openapi codegen-schemas-json ## Generate openapi and json schemas

.PHONY: codegen-playground
codegen-playground: build-wasm ## Generate playground
	@echo Generate playground... >&2
	@rm -rf ./pkg/server/ui/dist && mkdir -p ./pkg/server/ui/dist && cp -r ./website/playground/* ./pkg/server/ui/dist
	@rm -rf ./website/docs/_playground && mkdir -p ./website/docs/_playground && cp -r ./website/playground/* ./website/docs/_playground

.PHONY: codegen-helm-crds
codegen-helm-crds: codegen-crds ## Generate helm CRDs
	@echo Generate helm crds... >&2
	@cat $(CRDS_PATH)/* \
		| $(SED) -e '1i{{- if .Values.crds.install }}' \
		| $(SED) -e '$$a{{- end }}' \
		| $(SED) -e '/^  annotations:/a \ \ \ \ {{- end }}' \
 		| $(SED) -e '/^  annotations:/a \ \ \ \ {{- toYaml . | nindent 4 }}' \
		| $(SED) -e '/^  annotations:/a \ \ \ \ {{- with .Values.crds.annotations }}' \
 		| $(SED) -e '/^  annotations:/i \ \ labels:' \
		| $(SED) -e '/^  labels:/a \ \ \ \ {{- end }}' \
 		| $(SED) -e '/^  labels:/a \ \ \ \ {{- toYaml . | nindent 4 }}' \
		| $(SED) -e '/^  labels:/a \ \ \ \ {{- with .Values.crds.labels }}' \
		| $(SED) -e '/^  labels:/a \ \ \ \ {{- include "kyverno-json.labels" . | nindent 4 }}' \
 		> ./charts/kyverno-json/templates/crds.yaml

.PHONY: codegen-helm-docs
codegen-helm-docs: ## Generate helm docs
	@echo Generate helm docs... >&2
	@docker run -v ${PWD}/charts:/work -w /work jnorwood/helm-docs:v1.11.0 -s file

.PHONY: codegen
codegen: codegen-crds codegen-deepcopy codegen-register codegen-client codegen-docs codegen-mkdocs codegen-schemas codegen-playground codegen-helm-crds codegen-helm-docs ## Rebuild all generated code and docs

.PHONY: verify-codegen
verify-codegen: codegen ## Verify all generated code and docs are up to date
	@echo Checking codegen is up to date... >&2
	@git --no-pager diff -- .
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen".' >&2
	@echo 'To correct this, locally run "make codegen", commit the changes, and re-run tests.' >&2
	@git diff --quiet --exit-code -- .

########
# TEST #
########

.PHONY: tests
tests: $(CLI_BIN) ## Run tests
	@echo Running tests... >&2
	@go test ./...

########
# KIND #
########

.PHONY: kind-create
kind-create: $(KIND) ## Create kind cluster
	@echo Create kind cluster... >&2
	@$(KIND) create cluster --name $(KIND_NAME) --image $(KIND_IMAGE)

.PHONY: kind-delete
kind-delete: $(KIND) ## Delete kind cluster
	@echo Delete kind cluster... >&2
	@$(KIND) delete cluster --name $(KIND_NAME)

.PHONY: kind-load
kind-load: $(KIND) ko-build ## Build image and load in kind cluster
	@echo Load image... >&2
	@$(KIND) load docker-image --name $(KIND_NAME) $(KO_REGISTRY)/$(PACKAGE)/$(CLI_DIR):$(GIT_SHA)

.PHONY: kind-install
kind-install: $(HELM) kind-load ## Build image, load it in kind cluster and deploy helm chart
	@echo Install chart... >&2
	@$(HELM) upgrade --install kyverno-json --namespace kyverno-json --create-namespace --wait ./charts/kyverno-json \
		--set image.registry=$(KO_REGISTRY) \
		--set image.repository=$(PACKAGE)/$(CLI_DIR) \
		--set image.tag=$(GIT_SHA)

###########
# INSTALL #
###########

.PHONY: install-crds
install-crds: codegen-crds ## Install CRDs
	@echo Install CRDs... >&2
	@kubectl create -f ./$(CRDS_PATH)

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
