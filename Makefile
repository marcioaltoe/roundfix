ifneq (,$(wildcard .env))
include .env
export
endif

RTK := $(shell command -v rtk >/dev/null 2>&1 && echo rtk)
GO := $(RTK) go
GOFMT := $(RTK) gofmt

GOCACHE ?= $(CURDIR)/.gocache
export GOCACHE

APP := roundfix
CMD := ./cmd/roundfix
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP)
PKGS := ./...
# Local build identity for `roundfix --version`: short commit (plus -dirty
# when the tree has changes) and local build time. The release workflow
# stamps only app.Version from the tag and leaves these empty.
BUILD_COMMIT := $(shell commit=$$(git rev-parse --short HEAD 2>/dev/null) || exit 0; dirty=$$(git status --porcelain --untracked-files=all 2>/dev/null); if test -n "$$dirty"; then dirty=-dirty; else dirty=; fi; printf '%s%s' "$$commit" "$$dirty")
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S %z')
STAMP_LDFLAGS := -X 'roundfix/internal/app.BuildCommit=$(BUILD_COMMIT)' -X 'roundfix/internal/app.BuildTime=$(BUILD_TIME)'
BUILD_FLAGS ?= -buildvcs=false -ldflags "$(STAMP_LDFLAGS)"
RUN_FLAGS ?= -buildvcs=false
TARGET ?= project
GO_FILES := $(shell find . -name '*.go' -not -path './.git/*')

.DEFAULT_GOAL := help

.PHONY: help bootstrap verify fmt fmt-check test test-race baseline-digests build install run version clean deps skills-check skills-install skills-link skills-sync skills-sync-check

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

verify: fmt-check test skills-sync-check skills-check build ## Run the required local verification gate

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

DERIVED_DIGEST_PATHS := internal/baseline/assets/setups internal/baseline/testdata internal/baseline/assets/source-baselines internal/baseline/assets/formatter-fixtures internal/baseline/assets/profiles

baseline-digests: ## Regenerate derived Baseline digest artifacts
	@snapshot=$$(mktemp -t rf-digests) && \
	find $(DERIVED_DIGEST_PATHS) -type f -exec shasum {} + | sort > "$$snapshot" && \
	$(GO) test ./skills -run TestAuthorialSkillSync -update -count=1 && \
	$(GO) test ./internal/baseline -run TestCatalogCompatibility -update -count=1 && \
	$(GO) test ./internal/baseline -run TestBaselineCompatibilityCorpus -update -count=1 && \
	$(GO) test ./internal/baseline -run TestReadoptionCompatibilityMaintainedFixture -update -count=1 && \
	$(GO) test ./internal/baseline -run TestFormatterComposition -update -count=1 && \
	changed=$$(find $(DERIVED_DIGEST_PATHS) -type f -exec shasum {} + | sort | comm -13 "$$snapshot" - | awk '{print $$2}') && \
	rm -f "$$snapshot" && \
	if [ -z "$$changed" ]; then \
		echo "baseline-digests: no changes; derived artifacts already match their canonical sources"; \
	else \
		echo "baseline-digests: regenerated"; \
		echo "$$changed" | sed 's/^/  /'; \
	fi

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

skills-update: ## Install missing skills and update existing ones to latest (reads skills-lock.json)
	bunx skills experimental_install
	bunx skills update -p -y
	bun run fmt
	
skills-check: ## Validate shipped Roundfix skill artifacts
	$(GO) run $(RUN_FLAGS) $(CMD) skills check

# The Roundfix-owned skill bundle shipped in the binary: the operational
# roundfix skill plus the 13 authorial workflow skills. Everything else is
# managed by the external skills-lock.json origin and is never embedded.
OWNED_SKILLS := roundfix write-idea write-prd write-techspec write-tasks setup-context-driven implement-task implement-spec brainstorming council business-analyst archive-spec qa-gate evidence-gate

skills-sync: ## Regenerate the embedded skills/ bundle from canonical .agents/skills/
	@for s in $(OWNED_SKILLS); do rm -rf "skills/$$s"; cp -R ".agents/skills/$$s" "skills/$$s"; done

skills-sync-check: ## Fail when the embedded bundle drifts from canonical .agents/skills/
	@for s in $(OWNED_SKILLS); do \
		diff -r ".agents/skills/$$s" "skills/$$s" >/dev/null || { \
			echo "skills/$$s drifts from .agents/skills/$$s; run 'make skills-sync'"; exit 1; }; \
	done
	$(GO) test -count=1 ./skills -run 'TestNoPythonBaselineRuntime|TestThinSetupSkill|TestCheckRejectsExecutableSetupEngineArtifacts|TestRecommendedSkillsMatchLock'

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
