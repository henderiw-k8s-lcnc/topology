
# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= latest
REGISTRY ?= yndd
PROJECT ?= topology

KPT_BLUEPRINT_CFG_DIR ?= blueprint/fn-config
KPT_BLUEPRINT_PKG_DIR ?= blueprint/${PROJECT}
# IMAGE_TAG_BASE defines the docker.io namespace and part of the image name for remote images.
# This variable is used to construct full image tags for ndd packages.
IMAGE_TAG_BASE ?= $(REGISTRY)/${PROJECT}

# Image URL to use all building/pushing image targets
IMG ?= $(IMAGE_TAG_BASE)-controller:$(VERSION)

# Package
PKG ?= $(IMAGE_TAG_BASE)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Private Github REPOs
GITHUB_TOKEN ?= $(shell cat ~/.config/github.token)

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
##@ Development


.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen kpt kptgen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	mkdir -p ${KPT_BLUEPRINT_CFG_DIR}
	mkdir -p ${KPT_BLUEPRINT_PKG_DIR}/crd/bases
	mkdir -p ${KPT_BLUEPRINT_PKG_DIR}/app
	$(CONTROLLER_GEN) crd paths="./..." output:crd:artifacts:config=${KPT_BLUEPRINT_PKG_DIR}/crd/bases
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	kpt pkg init ${KPT_BLUEPRINT_PKG_DIR} --description "${PROJECT} controller"
	kpt pkg init ${KPT_BLUEPRINT_PKG_DIR}/crd --description "${PROJECT} crd"
	kpt pkg init ${KPT_BLUEPRINT_PKG_DIR}/app --description "${PROJECT} app"
	kptgen apply config ${KPT_BLUEPRINT_PKG_DIR} --fn-config-dir ${KPT_BLUEPRINT_CFG_DIR}
	rm ${KPT_BLUEPRINT_PKG_DIR}/package-context.yaml
	rm ${KPT_BLUEPRINT_PKG_DIR}/crd/package-context.yaml
	rm ${KPT_BLUEPRINT_PKG_DIR}/app/package-context.yaml

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: manifests generate fmt vet ## Run tests.
	##KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./... -coverprofile cover.out

##@ Build

##@ Deployment

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
KPT ?= $(LOCALBIN)/kpt
KPTGEN ?= $(LOCALBIN)/kptgen

## Tool Versions
KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_TOOLS_VERSION ?= v0.9.0
KPT_VERSION ?= main
KPTGEN_VERSION ?= v0.0.10
K8S_VERSION ?= v0.25.3


KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: kpt
kpt: $(KPT) ## Download kpt locally if necessary.
$(KPT): $(LOCALBIN)
	test -s $(LOCALBIN)/kpt || GOBIN=$(LOCALBIN) go install -v github.com/GoogleContainerTools/kpt@$(KPT_VERSION)

.PHONY: kptgen
kptgen: $(KPTGEN) ## Download kptgen locally if necessary.
$(KPTGEN): $(LOCALBIN)
	test -s $(LOCALBIN)/kptgen || GOBIN=$(LOCALBIN) go install -v github.com/henderiw-kpt/kptgen@$(KPTGEN_VERSION)

.PHONY: codegen
codegen: 
	go get -d $(CODE_GENERATOR)@$(K8S_VERSION)
	mkdir -p ${SCRIPT_DIR}
	cp $(GOPATH)/pkg/mod/$(CODE_GENERATOR)@$(K8S_VERSION)/generate-groups.sh $(SCRIPT_DIR)/
	chmod 755 ${SCRIPT_DIR}/generate-groups.sh