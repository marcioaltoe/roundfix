ifneq (,$(wildcard .env))
include .env
export
endif

RTK := $(shell command -v rtk >/dev/null 2>&1 && echo rtk)
GO := $(RTK) go
GOFMT := $(RTK) gofmt

APP := roundfix
CMD := ./cmd/roundfix
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP)
PKGS := ./...
BUILD_FLAGS ?= -buildvcs=false
RUN_FLAGS ?= -buildvcs=false
TARGET ?= project
GO_FILES := $(shell find . -name '*.go' -not -path './.git/*')

.DEFAULT_GOAL := help

.PHONY: help bootstrap verify fmt fmt-check test test-race build install run version clean deps skills-check skills-install skills-link

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n"} \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)


##@ Bootstrap

bootstrap: deps ## Download and verify Go modules

deps: ## Download, tidy, and verify Go modules
	$(GO) mod download
	$(GO) mod tidy
	$(GO) mod verify


##@ Quality & Testing

verify: fmt-check test skills-check build ## Run the required local verification gate

fmt: ## Format Go files
	$(GOFMT) -w $(GO_FILES)

fmt-check: ## Check Go formatting without changing files
	@test -z "$$($(GOFMT) -l $(GO_FILES))" || { \
		echo "Go files need formatting:"; \
		$(GOFMT) -l $(GO_FILES); \
		exit 1; \
	}

test: ## Run Go tests
	$(GO) test $(PKGS)

test-race: ## Run Go tests with the race detector
	$(GO) test -race $(PKGS)


##@ Build & Run

build: ## Build the CLI binary into bin/roundfix
	@mkdir -p $(BIN_DIR)
	$(GO) build $(BUILD_FLAGS) -o $(BIN) $(CMD)

install: build ## Build and install roundfix into Go bin for local testing
	$(GO) install $(BUILD_FLAGS) $(CMD)

run: ## Run the CLI; pass ARGS="--help" or another command
	$(GO) run $(RUN_FLAGS) $(CMD) $(ARGS)

version: ## Print the CLI version
	$(GO) run $(RUN_FLAGS) $(CMD) --version


##@ Cleanup

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)


##@ Agent Skills

skills-check: ## Validate shipped Roundfix skill artifacts
	$(GO) run $(RUN_FLAGS) $(CMD) skills check

skills-install: ## Install shipped Roundfix skills; pass TARGET=project|codex|claude|opencode|all
	$(GO) run $(RUN_FLAGS) $(CMD) skills install --target $(TARGET)

skills-link: ## Recreate .claude/skills symlinks from .agents/skills
	@mkdir -p .claude/skills
	@rm -f .claude/skills/*
	@for skill in .agents/skills/*/; do \
		name=$$(basename "$$skill"); \
		ln -s "../../.agents/skills/$$name" ".claude/skills/$$name"; \
	done
	@echo "Linked $$(ls .claude/skills | wc -l | tr -d ' ') skills from .agents/skills to .claude/skills"
