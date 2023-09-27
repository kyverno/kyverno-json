.DEFAULT_GOAL: build-all

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.20.0
CONTROLLER_GEN                     := $(TOOLS_DIR)/controller-gen
CONTROLLER_GEN_VERSION             := v0.12.0
CLIENT_GEN                         := $(TOOLS_DIR)/client-gen
LISTER_GEN                         := $(TOOLS_DIR)/lister-gen
INFORMER_GEN                       := $(TOOLS_DIR)/informer-gen
OPENAPI_GEN                        := $(TOOLS_DIR)/openapi-gen
REGISTER_GEN                       := $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       := $(TOOLS_DIR)/deepcopy-gen
DEFAULTER_GEN                      := $(TOOLS_DIR)/defaulter-gen
APPLYCONFIGURATION_GEN             := $(TOOLS_DIR)/applyconfiguration-gen
CODE_GEN_VERSION                   := v0.28.0
GEN_CRD_API_REFERENCE_DOCS         := $(TOOLS_DIR)/gen-crd-api-reference-docs
GEN_CRD_API_REFERENCE_DOCS_VERSION := latest
GO_ACC                             := $(TOOLS_DIR)/go-acc
GO_ACC_VERSION                     := latest
GOIMPORTS                          := $(TOOLS_DIR)/goimports
GOIMPORTS_VERSION                  := latest
HELM                               := $(TOOLS_DIR)/helm
HELM_VERSION                       := v3.12.3
HELM_DOCS                          := $(TOOLS_DIR)/helm-docs
HELM_DOCS_VERSION                  := v1.11.0
KO                                 := $(TOOLS_DIR)/ko
KO_VERSION                         := v0.14.1
KUTTL                              := $(TOOLS_DIR)/kubectl-kuttl
KUTTL_VERSION                      := v0.0.0-20230914072640-e3af68e47317
TOOLS                              := $(KIND) $(CONTROLLER_GEN) $(CLIENT_GEN) $(LISTER_GEN) $(INFORMER_GEN) $(OPENAPI_GEN) $(REGISTER_GEN) $(DEEPCOPY_GEN) $(DEFAULTER_GEN) $(APPLYCONFIGURATION_GEN) $(GEN_CRD_API_REFERENCE_DOCS) $(GO_ACC) $(GOIMPORTS) $(HELM) $(HELM_DOCS) $(KO) $(KUTTL)
ifeq ($(GOOS), darwin)
SED                                := gsed
else
SED                                := sed
endif
COMMA                              := ,

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)

$(CONTROLLER_GEN):
	@echo Install controller-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

$(CLIENT_GEN):
	@echo Install client-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/client-gen@$(CODE_GEN_VERSION)

$(LISTER_GEN):
	@echo Install lister-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/lister-gen@$(CODE_GEN_VERSION)

$(INFORMER_GEN):
	@echo Install informer-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/informer-gen@$(CODE_GEN_VERSION)

$(OPENAPI_GEN):
	@echo Install openapi-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/openapi-gen@$(CODE_GEN_VERSION)

$(REGISTER_GEN):
	@echo Install register-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(CODE_GEN_VERSION)

$(DEEPCOPY_GEN):
	@echo Install deepcopy-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GEN_VERSION)

$(DEFAULTER_GEN):
	@echo Install defaulter-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/defaulter-gen@$(CODE_GEN_VERSION)

$(APPLYCONFIGURATION_GEN):
	@echo Install applyconfiguration-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/applyconfiguration-gen@$(CODE_GEN_VERSION)

$(GEN_CRD_API_REFERENCE_DOCS):
	@echo Install gen-crd-api-reference-docs... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/ahmetb/gen-crd-api-reference-docs@$(GEN_CRD_API_REFERENCE_DOCS_VERSION)

$(GO_ACC):
	@echo Install go-acc... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/ory/go-acc@$(GO_ACC_VERSION)

$(GOIMPORTS):
	@echo Install goimports... >&2
	@GOBIN=$(TOOLS_DIR) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

$(HELM):
	@echo Install helm... >&2
	@GOBIN=$(TOOLS_DIR) go install helm.sh/helm/v3/cmd/helm@$(HELM_VERSION)

$(HELM_DOCS):
	@echo Install helm-docs... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/norwoodj/helm-docs/cmd/helm-docs@$(HELM_DOCS_VERSION)

$(KO):
	@echo Install ko... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/google/ko@$(KO_VERSION)

$(KUTTL):
	@echo Install kuttl... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/kyverno/kuttl/cmd/kubectl-kuttl@$(KUTTL_VERSION)

.PHONY: install-tools
install-tools: $(TOOLS) ## Install tools

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

###########
# CODEGEN #
###########

PACKAGE                     ?= github.com/eddycharly/tf-kyverno
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
	@mkdir -p $(GOPATH_SHIM)/src/github.com/eddycharly && ln -s -f ${PWD} $(PACKAGE_SHIM)

# TODO: doesn't work because structs lack DeepCopy
# .PHONY: codegen-register
# codegen-register: $(PACKAGE_SHIM) $(REGISTER_GEN) ## Generate types registrations
# 	@echo Generate registration... >&2
# 	@GOPATH=$(GOPATH_SHIM) $(REGISTER_GEN) \
# 		--go-header-file=./scripts/boilerplate.go.txt \
# 		--input-dirs=$(INPUT_DIRS)

# TODO: doesn't work with interface{}
# .PHONY: codegen-deepcopy
# codegen-deepcopy: $(PACKAGE_SHIM) $(DEEPCOPY_GEN) ## Generate deep copy functions
# 	@echo Generate deep copy functions... >&2
# 	@GOPATH=$(GOPATH_SHIM) $(DEEPCOPY_GEN) \
# 		--go-header-file=./scripts/boilerplate.go.txt \
# 		--input-dirs=$(INPUT_DIRS) \
# 		--output-file-base=zz_generated.deepcopy

.PHONY: codegen-crds
codegen-crds: $(CONTROLLER_GEN) ## Generate CRDs
	@echo Generate crds... >&2
	@$(CONTROLLER_GEN) crd paths=./pkg/apis/... crd:crdVersions=v1 output:dir=$(CRDS_PATH)

.PHONY: codegen-all
codegen-all: codegen-crds ## Rebuild all generated code

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
