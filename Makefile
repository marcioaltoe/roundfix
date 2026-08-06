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

.PHONY: help bootstrap verify spec-check spec-budget fmt fmt-check test test-race baseline-digests build install run version clean deps skills-check skills-install skills-link skills-sync skills-version-check skills-sync-check

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

verify: fmt-check test spec-budget skills-sync-check skills-check build spec-check ## Run the required local verification gate

spec-check: ## Check Spec artifact consistency
	$(BIN) spec check

spec-budget: ## Prove the Spec corpus sweep stays within its budget
	$(GO) test -count=1 -parallel=1 ./internal/speccheck -run '^TestCheckCorpusBudget$$'

fmt: ## Format Go files
	$(GOFMT) -w $(GO_FILES)

fmt-check: ## Check Go formatting without changing files
	@test -z "$$($(GOFMT) -l $(GO_FILES))" || { \
		echo "Go files need formatting:"; \
		$(GOFMT) -l $(GO_FILES); \
		exit 1; \
	}

# -parallel decides how many tests may overlap inside one package. Go defaults
# it to GOMAXPROCS, which is right for CPU-bound tests and wrong for these:
# they spend their time in temp directories, git subprocesses, and child
# processes, so a core count leaves them queueing on an idle machine. On a
# two-core runner internal/cli measured 129.3s at the default and 46.4s at 16;
# the curve is flat past that, and on a twelve-core machine it changes nothing.
GO_TEST_PARALLEL ?= 16

test: ## Run Go tests
	$(GO) test -parallel $(GO_TEST_PARALLEL) $(PKGS)

# The full suite's wall clock on a fresh run, in seconds. Derived from what the
# parallelised suite actually achieves on a GitHub runner, not from the old
# baseline: on 2026-08-03 the CI suite measured ~195s, dominated by
# internal/cli at 192.5s. The headroom absorbs runner variance while still
# catching a real regression — re-serialising internal/cli would put the suite
# near 470s on the same hardware. Change this deliberately, and say why.
# See docs/specs/_archived/0071-verification-cost/baseline/2026-08-03-after.md.
SUITE_BUDGET_SECONDS ?= 360

test-budget: ## Run Go tests once and fail if the suite exceeds SUITE_BUDGET_SECONDS
	@start=$$(date +%s); \
	$(GO) test -count=1 -parallel $(GO_TEST_PARALLEL) $(PKGS) || exit $$?; \
	elapsed=$$(( $$(date +%s) - start )); \
	if [ "$$elapsed" -gt "$(SUITE_BUDGET_SECONDS)" ]; then \
		echo "suite-time budget exceeded: the suite took $${elapsed}s against a budget of $(SUITE_BUDGET_SECONDS)s"; \
		echo "the budget is SUITE_BUDGET_SECONDS in the Makefile; raise it deliberately or find the regression"; \
		exit 1; \
	fi; \
	echo "suite-time budget: $${elapsed}s of $(SUITE_BUDGET_SECONDS)s"

test-race: ## Run Go tests with the race detector
	$(GO) test -race -parallel $(GO_TEST_PARALLEL) $(PKGS)

DERIVED_DIGEST_PATHS := internal/baseline/assets/setups internal/baseline/testdata internal/baseline/assets/source-baselines internal/baseline/assets/formatter-fixtures internal/baseline/assets/profiles
BASELINE_DIGEST_STEPS := \
	./internal/baseline:TestReadoptionCompatibilityMaintainedFixture \
	./skills:TestAuthorialSkillSync \
	./internal/baseline:TestFormatterComposition \
	./internal/baseline:TestBaselineCompatibilityCorpus \
	./internal/baseline:TestCatalogCompatibility \
	./internal/baseline:^TestCatalogDiagnosticCharacterization$$:-update-catalog-diagnostics \
	./internal/baseline:TestBaselinePlanCharacterization:-update-baseline-plan-characterization

baseline-digests: ## Regenerate derived Baseline digest artifacts
	@raw=""; snapshot=""; \
	err_code="temp_file_unavailable"; err_stage="snapshot-allocation"; err_retryable="true"; \
	err_next="Free space in the temporary directory or point TMPDIR at a writable one, then rerun make baseline-digests."; \
	finish() { status=$$?; rm -f "$$snapshot" "$$raw"; if [ "$$status" -ne 0 ]; then printf '{"schemaVersion":1,"type":"baseline-digests","ok":false,"changed":false,"errorCode":"%s","stage":"%s","retryable":%s,"nextSteps":"%s"}\n' "$$err_code" "$$err_stage" "$$err_retryable" "$$err_next"; fi; exit "$$status"; }; \
	trap finish EXIT; \
	snapshot=$$(mktemp "$${TMPDIR:-/tmp}/rf-digests.XXXXXX") || exit $$?; \
	raw=$$(mktemp "$${TMPDIR:-/tmp}/rf-digests-raw.XXXXXX") || exit $$?; \
	err_code="artifact_scan_failed"; err_stage="pre-scan"; err_retryable="false"; \
	err_next="Verify every path in DERIVED_DIGEST_PATHS exists and is readable, then rerun make baseline-digests."; \
	find $(DERIVED_DIGEST_PATHS) -type f -exec shasum {} + > "$$raw" || exit $$?; \
	err_code="snapshot_sort_failed"; err_stage="pre-scan-sort"; \
	err_next="Inspect the temporary directory for space or permission problems, then rerun make baseline-digests."; \
	sort "$$raw" > "$$snapshot" || exit $$?; \
	err_code="regeneration_failed"; err_retryable="false"; \
	err_next="Read the failing test output above, fix the canonical source it validates, then rerun make baseline-digests."; \
	for step in $(BASELINE_DIGEST_STEPS); do package=$${step%%:*}; test_spec=$${step#*:}; test_name=$${test_spec%%:*}; update_flag=$${test_spec#*:}; if [ "$$update_flag" = "$$test_spec" ]; then update_flag=-update; fi; err_stage="$$package:$$test_name"; $(GO) test "$$package" -run "$$test_name" "$$update_flag" -count=1 >&2 || { status=$$?; printf 'baseline-digests: regeneration failed at %s:%s\n' "$$package" "$$test_name" >&2; exit "$$status"; }; done; \
	err_code="artifact_scan_failed"; err_stage="post-scan"; \
	err_next="Verify every path in DERIVED_DIGEST_PATHS exists and is readable, then rerun make baseline-digests."; \
	find $(DERIVED_DIGEST_PATHS) -type f -exec shasum {} + > "$$raw" || exit $$?; \
	err_code="comparison_failed"; err_stage="comparison"; \
	err_next="Rerun make baseline-digests; if the comparison keeps failing, inspect the derived artifact paths for unreadable files."; \
	changed=$$(sort "$$raw" | comm -3 "$$snapshot" - | awk '{print $$2}' | sort -u) || exit $$?; \
	err_code="strict_validation_failed"; err_stage="strict-validation"; err_retryable="false"; \
	err_next="Read the strict catalog validation output above, repair any remaining inconsistency, then rerun make baseline-digests."; \
	$(GO) test ./internal/baseline -run TestCatalogCompatibility -count=1 >&2 || { status=$$?; printf 'baseline-digests: strict validation failed\n' >&2; exit "$$status"; }; \
	if [ -z "$$changed" ]; then result_changed=false; printf '%s\n' "baseline-digests: no changes; derived artifacts already match their canonical sources" >&2; else result_changed=true; printf '%s\n' "baseline-digests: regenerated" >&2; printf '%s\n' "$$changed" | sed 's/^/  /' >&2; fi; printf '{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":%s}\n' "$$result_changed"

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

skills-version-check: ## Fail when an owned skill omits its declared version
	@for s in $(OWNED_SKILLS); do \
		for root in .agents/skills skills; do \
			file="$$root/$$s/SKILL.md"; \
			awk 'NR == 1 && $$0 == "---" { frontmatter = 1; next } frontmatter && $$0 == "---" { exit } frontmatter && $$0 ~ /^version:[[:space:]]+[^[:space:]#]+([[:space:]]+#.*)?[[:space:]]*$$/ { found = 1 } END { exit !found }' "$$file" || { \
				echo "$$file does not declare a top-level version"; exit 1; }; \
		done; \
	done

skills-sync-check: skills-version-check ## Fail when the embedded bundle drifts from canonical .agents/skills/
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
