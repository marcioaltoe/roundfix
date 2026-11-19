---
task: task_05
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: completed
type: docs
complexity: high
---

# Task 05: This repository marks its decisions and publishes only its README, its user guide and its lessons

## Overview

This repository applies the rules to itself. The Agent marks every ADR in the tree with the Knowledge Mark by the test of the export rule, writes the dated notices that ADR-0438 and ADR-0440 give ADR-0248, ADR-0278 and ADR-0412, and runs the Managed Refresh so the update fills their Amendment Lists and names the ADRs marked true in the Standard Export List. The export file ends with this repository's own lines, which drop the glossary, the references, the design notes and every ADR, so that only the README, the user guide and the lessons are published; its contract test follows the new rule, and the coverage record is recorded again.

## Requirements

1. MUST give every file under `docs/adr/` in the tree when the Task runs, the ADRs of other Specs that merged before it included, the Knowledge Mark `knowledge: true` or `knowledge: false` on the line right after `status` in its header, by the test of `_techspec.md` → API Contract 24 and `clause.secondbrain.export-03-knowledge-mark`: for this repository the product is Roundfix, so a decision about what a command, the Daemon, the Delivery Queue, the Baseline or a skill does for the people and repositories that use Roundfix is `true`, and a decision about this repository's tests, build, CI, release procedure or internal structure is `false`. The mark is a judgment made one ADR at a time; the Task adds exactly one line to each file for it, never reflows or re-saves a file, never moves `updated_at`, and changes no other line (`_techspec.md` → Risks & Considerations). The ADR files it marks, other than the three of requirement 2, are recorded by the Daemon.
2. MUST write the three dated notices of `_techspec.md` → API Contract 23, word for word, each after the text it changes: ADR-0248's after its paragraph that ends "A documentation contract holds both directions.", ADR-0278's after its paragraph that begins "**Lessons leave before the text.**", and ADR-0412's at the end of its decision text. No other line of the three changes by hand, and the decision text of no accepted ADR changes.
3. MUST replace every line of `.secondbrain-export` outside the Standard Export List with the own lines of API Contract 22: after the end marker, one blank line and exactly these seven lines, in this order, so that `CHANGELOG.md` and the lines the conversion kept leave the file (`_prd.md` → Decisions, "What does this repository's own variation drop"):
   - `# This repository publishes its README, its user guide and its lessons (maintainer, 2026-11-18).`
   - `# The glossary, the references and the design notes describe how the tool is built.`
   - `!CONTEXT.md`
   - `!docs/references/`
   - `!docs/design/`
   - `# Its decision records are about the tool; the user guide carries what other projects use.`
   - `!docs/adr/*`
   The lines between the two markers are never edited by hand.
4. MUST then run `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` until a run reports nothing to do: the run fills the Amendment Lists of ADR-0248, ADR-0278 and ADR-0412 from the new notices and rewrites the Standard Export List so it names the ADRs marked true, and changes no line outside the markers. Where an ADR's header already holds `amended_by`, as ADR-0248's and ADR-0278's do since task_02's refresh, the run replaces that one line with the longer list that keeps every identifier it held (`_techspec.md` → ADR header and the lesson form, "writes the Amendment List by inserting or replacing one line"); that is the only line of an ADR the Task may remove.
5. MUST rewrite `internal/docscontract/secondbrain_export_test.go` as API Contract 25 says, keeping its build tag `docscontract` and its `//verify:always` line: it holds `TestThisRepositoryExportsNoProcessPath`, which fails when the file's own lines publish a process path or lack a reason, through `ExportVariationDepartures`, which task_03 adds, and `TestTheExportRuleRefusesAProcessPath`, which proves on a temporary directory, in the subtests `docs_whole` and `agents_md`, that the same rule refuses `docs/` whole and `AGENTS.md`; and it no longer holds `TestSecondbrainExportFollowsTheHistoryForm` or `TestSecondbrainExportRuleRefusesBothMismatches`. Neither new test asserts that every ADR carries the mark or that the Standard Export List is current (`_prd.md` → Decisions).
6. MUST record the coverage record again, never by hand, with `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`, because the two replaced tests are named in it.
7. MUST NOT change any Go file other than `internal/docscontract/secondbrain_export_test.go`, any skill, any guide under `docs/agents/`, any Baseline module, template, profile, golden or Source Baseline, any file under `docs/user-guide/`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, `AGENTS.md`, `CHANGELOG.md`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration, and creates no file under `docs/lessons/`. Of the Governed Paths it changes only `docs/references/coverage-record.json`, listed in `_authorization.md` `paths`.

## Subtasks

- [ ] Mark every ADR in the tree, one line each.
- [ ] Write the three amendment notices.
- [ ] Write this repository's own lines of the export file and run the Managed Refresh until it reports nothing to do.
- [ ] Rewrite the export contract test and record the coverage record again.

## Acceptance Criteria

- [ ] Every ADR in the tree carries the Knowledge Mark on the line after `status`, and no ADR lost a line other than an `amended_by` line the refresh replaced with a list that keeps every identifier it held.
- [ ] ADR-0248, ADR-0278 and ADR-0412 carry their dated notices and list the ADRs that amend them in their headers.
- [ ] This repository's update reports `current` with a current documentation part and no unmarked ADR, and its export file is the Standard Export List followed by exactly the own lines of API Contract 22.
- [ ] The export contract tests hold that this repository exports no process path and that the rule refuses one, and the coverage record is current.

## Context

- instruction: `docs/adr/0439-an-adrs-header-says-whether-the-secondbrain-wants-it-and-which-adrs-amend-it.md`
- instruction: `docs/agents/secondbrain.md`
- interface: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- interface: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- interface: `docs/adr/0412-planning-references-and-the-release-preparation-commit-set-no-release-impact.md`
- interface: `.secondbrain-export`
- interface: `internal/docscontract/secondbrain_export_test.go`
- interface: `docs/references/coverage-record.json`

## Consumers

- This repository's export file and its contract tests — consumers: `./internal/docscontract`
- The coverage record — consumers: `./internal/spec`

## Verification

- `for f in docs/adr/*.md; do awk 'NR == 1 && $0 == "---" { fm = 1; next } fm && $0 == "---" { exit } fm && previous ~ /^status:/ { found = 1; ok = ($0 ~ /^knowledge: (true|false)$/) } fm { previous = $0 } END { exit((found && ok) ? 0 : 1) }' "$f" || { echo "no Knowledge Mark on the line after status: $f" >&2; exit 1; }; done; changed="$(git diff --numstat HEAD -- docs/adr | awk '$2 != 0 { print $3 }')" || exit 1; for f in $changed; do old="$(git diff -U0 HEAD -- "$f" | awk '/^--- / { next } /^-/ { print }')" || exit 1; test "$(printf '%s\n' "$old" | grep -c '^-amended_by: \[.*\]$')" -eq 1 && test "$(printf '%s\n' "$old" | wc -l | tr -d ' ')" -eq 1 || { printf 'ADR lines removed: %s\n' "$f" >&2; exit 1; }; new="$(grep '^amended_by: \[' "$f")" || { printf 'Amendment List removed: %s\n' "$f" >&2; exit 1; }; for id in $(printf '%s\n' "$old" | grep -oE 'ADR-[0-9]+'); do printf '%s\n' "$new" | grep -qE -- "$id[],]" || { printf 'Amendment List lost %s: %s\n' "$id" "$f" >&2; exit 1; }; done; done` — expected: exit 0. Before this Task no ADR carries the mark; afterwards every ADR carries it on the line after `status`, and no ADR lost a line other than an `amended_by` line the refresh replaced with a list that keeps every identifier it held.
- `for pair in 'docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md|The export no longer mirrors' 'docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md|A lesson is no longer sent as a message:' 'docs/adr/0412-planning-references-and-the-release-preparation-commit-set-no-release-impact.md|Extended by [ADR-0438](0438-'; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; grep -q '^amended_by: \[.*ADR-0438' docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md && grep -q '^amended_by: \[.*ADR-0440' docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md && grep -q '^amended_by: \[.*ADR-0438.*ADR-0440' docs/adr/0412-planning-references-and-the-release-preparation-commit-set-no-release-impact.md || { echo 'an Amendment List lacks its new identifier' >&2; exit 1; }; actual="$(awk 'after { print } /^# roundfix:standard-export:end$/ { after = 1 }' .secondbrain-export)"; expected="$(printf '%s\n' '' '# This repository publishes its README, its user guide and its lessons (maintainer, 2026-11-18).' '# The glossary, the references and the design notes describe how the tool is built.' '!CONTEXT.md' '!docs/references/' '!docs/design/' '# Its decision records are about the tool; the user guide carries what other projects use.' '!docs/adr/*')"; test "$actual" = "$expected" || { printf 'unexpected own lines:\n%s\n' "$actual" >&2; exit 1; }; bin_dir="$(mktemp -d)" || exit 1; home_dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$bin_dir/roundfix" ./cmd/roundfix || { rm -rf "$bin_dir" "$home_dir"; exit 1; }; out="$(HOME="$home_dir" "$bin_dir/roundfix" baseline update --repo . --no-skills --no-history --format json </dev/null 2>&1)"; status=$?; rm -rf "$bin_dir" "$home_dir"; test "$status" -eq 0 || { printf '%s\n' "$out"; exit 1; }; for want in '"state":"current"' '"documentation":{"status":"current"' '"unmarked":0}'; do printf '%s\n' "$out" | grep -qF -- "$want" || { printf 'missing %s in:\n%s\n' "$want" "$out" >&2; exit 1; }; done` — expected: exit 0. Before this Task the three notices and the own lines are absent and every ADR is unmarked; afterwards the notices are written and listed in the headers, the export file ends with exactly the own lines of API Contract 22, and this repository's update under a fresh home reports `current` with a current documentation part and no unmarked ADR.
- `if grep -q 'func TestSecondbrainExportFollowsTheHistoryForm' internal/docscontract/secondbrain_export_test.go; then echo 'the history-form test is still present' >&2; exit 1; fi; out="$(go test -count=1 -timeout 20m -v -tags docscontract -run '^(TestThisRepositoryExportsNoProcessPath|TestTheExportRuleRefusesAProcessPath)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThisRepositoryExportsNoProcessPath TestTheExportRuleRefusesAProcessPath/docs_whole TestTheExportRuleRefusesAProcessPath/agents_md; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the file holds the history-form tests and not the new ones; afterwards the two new tests pass, the coverage record matches the tests, and the documentation contracts pass.

## References

- `_prd.md` → Goals; User Stories 1, 5, 6 and 7; Core Features 6, 7, 8 and 19; Success Metric 6; Paths that leave a mirror
- `_techspec.md` → This repository; API Contract 22; API Contract 23; API Contract 24; API Contract 25; Risks & Considerations; Build Order 5
- ADR-0438, ADR-0439, ADR-0440, ADR-0248, ADR-0278, ADR-0412

## Result

Skills loaded: implement-task, golang-testing, testing-boss, tech-writer, setup-context-driven, roundfix, evidence-gate.

- All 376 ADRs carry a Knowledge Mark immediately after status (323 true, 53 false); three already carried marks. A focused comparison against HEAD found no removed ADR lines except replaced Amendment Lists retaining their identifiers.
- Added the exact dated notices to ADR-0248, ADR-0278 and ADR-0412. Managed Refresh populated their Amendment Lists, preserving the existing identifiers.
- Replaced the repository's own export lines with the seven specified lines. The required Managed Refresh updated the Standard Export List and three headers; a subsequent run reported zero file changes. Its JSON documentation part reported current with 376 ADRs and zero unmarked.
- Replaced the export contract tests; focused runs of `TestThisRepositoryExportsNoProcessPath` and `TestTheExportRuleRefusesAProcessPath` passed. Ran the required coverage-record generator successfully after replacing the tests. Full Daemon Verification remains the Daemon's responsibility.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/adr/0001-daemon-owns-round-commits.md`
- `docs/adr/0002-yaml-user-and-project-config.md`
- `docs/adr/0004-central-run-database.md`
- `docs/adr/0006-fetch-creates-a-run.md`
- `docs/adr/0007-resolve-uses-downloaded-review-issues.md`
- `docs/adr/0008-run-event-payload-stores-raw-producer-json.md`
- `docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md`
- `docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md`
- `docs/adr/0011-acp-runtimes-run-batches-with-full-access.md`
- `docs/adr/0017-roundfix-drives-acp-runtimes-through-acpx.md`
- `docs/adr/0019-watch-clean-means-merge-ready.md`
- `docs/adr/0020-parsed-prompt-result-outranks-acpx-exit-code.md`
- `docs/adr/0022-stop-requests-travel-through-the-run-database.md`
- `docs/adr/0023-runs-execute-in-per-run-worktrees.md`
- `docs/adr/0024-run-integration-is-porcelain-only.md`
- `docs/adr/0025-spec-runs-schedule-ready-tasks-concurrently.md`
- `docs/adr/0026-task-integration-is-a-serialized-cherry-pick-queue.md`
- `docs/adr/0027-removed-config-keys-degrade-to-warnings.md`
- `docs/adr/0028-operational-runs-can-detach-from-their-caller.md`
- `docs/adr/0029-review-artifacts-live-with-the-spec.md`
- `docs/adr/0030-agent-run-logs-are-opt-in.md`
- `docs/adr/0031-roundfix-ships-through-npm-platform-packages.md`
- `docs/adr/0032-roundfix-detects-and-avoids-an-unsafe-codex.md`
- `docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md`
- `docs/adr/0034-run-worktrees-can-bootstrap-their-environment.md`
- `docs/adr/0035-spec-root-is-configurable-and-external-artifacts-stay-uncommitted.md`
- `docs/adr/0036-review-artifacts-are-committed-in-a-separate-docs-commit.md`
- `docs/adr/0037-roundfix-owns-agent-model-and-reasoning-selection.md`
- `docs/adr/0038-daemon-allows-one-verification-repair.md`
- `docs/adr/0044-orphaned-run-locks-are-reclaimed-on-proven-owner-death.md`
- `docs/adr/0049-agent-selection-profiles-are-atomic.md`
- `docs/adr/0050-configured-fallbacks-activate-after-notification.md`
- `docs/adr/0051-tasks-and-qa-own-agent-sessions.md`
- `docs/adr/0052-run-completion-is-compare-and-set.md`
- `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md`
- `docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md`
- `docs/adr/0057-daemon-exclusively-owns-implement-task-status.md`
- `docs/adr/0058-baseline-upgrades-fail-closed-on-unaccounted-rule-removal.md`
- `docs/adr/0059-generated-output-is-formatter-stable-in-the-target-repository.md`
- `docs/adr/0060-source-baselines-are-exhaustive-and-project-agnostic.md`
- `docs/adr/0061-standard-typescript-monorepo-is-opinionated.md`
- `docs/adr/0062-roundfix-owned-versions-restart-at-zero.md`
- `docs/adr/0063-repositories-own-their-http-contract.md`
- `docs/adr/0064-baseline-readoption-uses-byte-exhaustive-structural-inventory.md`
- `docs/adr/0066-context-driven-baseline-execution-belongs-to-the-cli.md`
- `docs/adr/0067-custom-baseline-profiles-are-repository-owned.md`
- `docs/adr/0068-baseline-command-uses-one-confirmation-gated-workflow.md`
- `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md`
- `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
- `docs/adr/0071-baseline-plans-are-portable-and-preimage-bound.md`
- `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- `docs/adr/0073-baseline-apply-uses-a-recoverable-multi-file-transaction.md`
- `docs/adr/0074-repository-rules-use-hybrid-semantic-ownership.md`
- `docs/adr/0075-profile-divergence-uses-confirmed-repository-owned-adaptation.md`
- `docs/adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md`
- `docs/adr/0077-new-specs-carry-a-readable-project-constraint-snapshot.md`
- `docs/adr/0078-confirmed-root-rules-move-to-semantic-owners.md`
- `docs/adr/0079-independent-reasoning-controls-make-model-identifiers-opaque.md`
- `docs/adr/0080-qa-verdicts-distinguish-environment-blocked-rows.md`
- `docs/adr/0081-sanctioned-digest-regeneration-is-fallout-of-the-authorized-edit.md`
- `docs/adr/0083-adopted-sources-move-to-their-owning-spec.md`
- `docs/adr/0085-a-regeneration-run-is-not-gated-on-the-pins-it-rewrites.md`
- `docs/adr/0086-category-removal-is-declared-by-flag-not-by-fragment-shape.md`
- `docs/adr/0087-executable-discovery-resolves-links-without-executing.md`
- `docs/adr/0088-the-qa-gate-is-authored-into-the-graph.md`
- `docs/adr/0089-code-under-test-takes-its-environment-explicitly.md`
- `docs/adr/0090-repository-facts-are-read-in-batches-never-cached-across-mutations.md`
- `docs/adr/0091-the-qa-gate-is-a-task-node-of-its-own-type.md`
- `docs/adr/0092-triage-intent-lives-in-a-typed-backlog.md`
- `docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md`
- `docs/adr/0094-the-consistency-check-is-artifact-presence-aware.md`
- `docs/adr/0096-the-qa-gate-proves-machine-facts-before-it-spends-an-agent-turn.md`
- `docs/adr/0097-a-qa-row-carries-forward-only-on-declared-unmoved-evidence.md`
- `docs/adr/0098-run-events-append-in-batches-with-durability-per-batch.md`
- `docs/adr/0099-retention-accounting-is-mechanical-classification-is-not.md`
- `docs/adr/0100-a-managed-refresh-proves-preservation-instead-of-backing-it-up.md`
- `docs/adr/0103-an-applied-refresh-republishes-the-setup-manifest.md`
- `docs/adr/0104-a-spec-accepts-on-evidence-it-did-not-author.md`
- `docs/adr/0105-a-capability-projection-bounds-what-it-keeps-not-what-a-runtime-advertises.md`
- `docs/adr/0107-profile-readiness-covers-every-configured-work-category.md`
- `docs/adr/0108-roundfix-warms-an-opencode-session-to-apply-a-reasoning-effort.md`
- `docs/adr/0110-a-negative-control-is-authored-and-recorded-never-synthesised.md`
- `docs/adr/0111-an-unobserved-verification-is-unknown-not-a-verdict.md`
- `docs/adr/0113-a-run-outcome-is-derived-from-unresolved-work-not-from-a-failed-batch.md`
- `docs/adr/0114-opening-an-agent-session-is-not-agent-work.md`
- `docs/adr/0115-a-superseded-run-branch-is-discarded-through-roundfix.md`
- `docs/adr/0116-a-citation-is-checked-against-what-it-cites.md`
- `docs/adr/0117-a-defect-is-checked-by-the-stage-that-can-produce-it.md`
- `docs/adr/0118-a-repository-constant-that-varies-becomes-a-decision.md`
- `docs/adr/0120-retired-documentation-lives-under-one-history-root.md`
- `docs/adr/0121-an-archive-move-is-a-ledger-of-identities-never-of-content.md`
- `docs/adr/0122-a-proposed-decision-is-pending-not-retired.md`
- `docs/adr/0125-a-spawned-fixture-is-a-compiled-binary-never-a-written-script.md`
- `docs/adr/0126-the-suite-proves-its-own-isolation-per-package.md`
- `docs/adr/0127-process-residue-is-a-readiness-fact-not-a-run-record.md`
- `docs/adr/0130-the-audit-judges-governed-paths-and-history-keeps-the-set-honest.md`
- `docs/adr/0131-the-tooling-row-states-applicability-not-mutation.md`
- `docs/adr/0132-a-refused-gate-records-the-refusal-as-its-terminal-row.md`
- `docs/adr/0133-a-diagnostic-names-the-literal-it-requires.md`
- `docs/adr/0134-the-gate-performs-the-repairs-its-task-names.md`
- `docs/adr/0135-an-absent-diagnostic-is-a-reported-state.md`
- `docs/adr/0136-a-repeated-failure-is-recognised-by-a-normalised-signature.md`
- `docs/adr/0137-the-run-budget-is-explained-where-it-is-configured.md`
- `docs/adr/0138-spec-run-git-policy-commit-per-task-push-only-at-clean.md`
- `docs/adr/0139-one-active-run-per-work-target-keyed-in-the-run-database.md`
- `docs/adr/0140-an-agent-selection-is-proven-as-the-exact-advertised-tuple.md`
- `docs/adr/0141-review-runs-execute-in-a-clean-tracked-user-checkout.md`
- `docs/adr/0142-head-bound-review-source-evidence-decides-the-watch-outcome.md`
- `docs/adr/0143-release-planning-is-read-only-including-its-reset-mode.md`
- `docs/adr/0144-setup-is-declarative-manifest-identity-and-declared-decision-effects.md`
- `docs/adr/0145-a-managed-region-is-trusted-by-its-marker-and-refreshed-when-unrecorded.md`
- `docs/adr/0146-release-publication-is-all-or-nothing-with-a-bounded-token-fallback.md`
- `docs/adr/0147-the-adapters-refusal-fires-first-and-the-catalogue-is-the-net.md`
- `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`
- `docs/adr/0149-one-regeneration-declaration-the-grant-names-the-command-the-tree-names-outputs.md`
- `docs/adr/0150-work-branches-express-purpose-and-run-branches-express-ownership.md`
- `docs/adr/0151-configured-review-profiles-select-the-independent-reviewer.md`
- `docs/adr/0153-pre-pr-review-is-an-explicit-provider-policy.md`
- `docs/adr/0154-a-qa-archive-override-records-user-authority-not-a-pass.md`
- `docs/adr/0155-the-qa-task-declares-the-gate-matrix.md`
- `docs/adr/0156-a-promise-a-spec-declares-names-a-consuming-task.md`
- `docs/adr/0157-tracked-source-is-work-output-whatever-its-mode.md`
- `docs/adr/0158-a-budget-expired-run-settles-budget-exceeded-and-stays-recoverable.md`
- `docs/adr/0159-verification-runs-past-a-failure-only-when-declared-independent.md`
- `docs/adr/0160-only-the-frozen-authorization-opens-a-red-repository-gate.md`
- `docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md`
- `docs/adr/0163-review-artifact-retirement-reads-recorded-evidence.md`
- `docs/adr/0164-an-implement-run-budget-renews-at-each-task-settlement.md`
- `docs/adr/0165-a-blocking-review-after-archive-parks-publication-for-a-corrective-spec.md`
- `docs/adr/0166-the-daemon-records-the-paths-a-task-changed-without-declaring-them.md`
- `docs/adr/0167-the-pre-pr-pull-request-row-never-decides-a-qualifying-partial.md`
- `docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md`
- `docs/adr/0169-the-pre-pr-review-diffs-the-candidate-from-its-merge-base.md`
- `docs/adr/0170-a-task-already-completed-on-its-target-is-nothing-to-carry.md`
- `docs/adr/0171-a-retention-prune-reports-only-what-it-reclaimed.md`
- `docs/adr/0172-doctor-reports-reclaimable-run-storage-from-cheap-reads.md`
- `docs/adr/0173-a-baseline-plan-reports-the-citations-its-history-relocations-break.md`
- `docs/adr/0174-a-pre-pr-review-reads-the-final-message-and-keeps-a-record-per-checkout.md`
- `docs/adr/0176-citation-checks-read-only-what-a-specs-authors-wrote.md`
- `docs/adr/0177-the-relocation-citation-scan-reads-through-a-repository-root.md`
- `docs/adr/0178-a-task-commit-is-authorized-by-the-grant-it-ran-under.md`
- `docs/adr/0179-an-explicitly-empty-paths-list-grants-operations-only.md`
- `docs/adr/0180-built-in-selections-derive-from-one-dated-recommended-profile.md`
- `docs/adr/0181-a-configuration-is-compared-with-the-recommendation-only-where-a-person-asked.md`
- `docs/adr/0182-a-task-settles-on-the-facts-its-gate-will-check.md`
- `docs/adr/0183-an-attributed-claim-carries-a-receipt-the-check-proves.md`
- `docs/adr/0184-a-techspec-states-a-command-surface-as-a-transcript.md`
- `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- `docs/adr/0187-the-roundfix-skill-and-the-command-reference-are-read-one-command-at-a-time.md`
- `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- `docs/adr/0190-a-stack-rule-names-the-language-or-workspace-it-governs.md`
- `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- `docs/adr/0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md`
- `docs/adr/0193-a-spec-names-its-prerequisites-and-the-owner-waits-for-them.md`
- `docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md`
- `docs/adr/0195-rows-that-read-the-gate-or-the-commits-are-observed-on-every-pass.md`
- `docs/adr/0196-a-pre-pr-review-finding-parks-only-after-validation.md`
- `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- `docs/adr/0198-a-run-counts-the-tokens-its-adapters-report-by-each-reports-own-scope.md`
- `docs/adr/0199-a-queue-token-ceiling-stops-new-starts-and-never-a-running-task.md`
- `docs/adr/0200-spec-authoring-gets-an-advisory-judge-that-never-gates.md`
- `docs/adr/0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md`
- `docs/adr/0202-a-stack-rule-carries-force-and-a-preference-names-what-licenses-the-exception.md`
- `docs/adr/0203-a-repository-rule-becomes-a-baseline-clause-only-when-it-recurs.md`
- `docs/adr/0204-a-composed-profile-takes-a-setup-composed-from-upstream-setups-by-name.md`
- `docs/adr/0205-a-repository-records-its-frontend-layout-and-an-unrecorded-layout-follows-the-suggestion.md`
- `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`
- `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`
- `docs/adr/0209-the-judge-suggests-which-sources-share-a-spec.md`
- `docs/adr/0210-an-evidence-snapshot-records-one-digest-per-declared-input.md`
- `docs/adr/0211-the-agent-environment-drops-a-node-preload-whose-file-is-missing.md`
- `docs/adr/0212-a-merged-spec-supersedes-its-runs-spec-directory-work-and-declared-leftovers.md`
- `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`
- `docs/adr/0214-why-a-verification-failed-is-classified-by-deterministic-signatures.md`
- `docs/adr/0215-a-measurement-lands-as-a-record-and-archived-evidence-leaves-only-with-approval.md`
- `docs/adr/0217-a-cursor-selection-names-the-advertised-model-value-and-the-login-stays-the-maintainers.md`
- `docs/adr/0219-a-profile-draft-binds-to-the-closest-built-in-profile-it-adapts.md`
- `docs/adr/0220-doctor-checks-what-delivery-needs-and-only-a-definite-answer-fails.md`
- `docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md`
- `docs/adr/0222-a-baseline-guide-says-only-what-holds-for-the-repository-that-reads-it.md`
- `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- `docs/adr/0224-this-repository-holds-its-upstream-skills-to-the-snapshot-and-the-lock-not-to-a-pinned-digest.md`
- `docs/adr/0225-a-queue-item-runs-the-roundfix-binary-its-branch-builds.md`
- `docs/adr/0226-a-verification-command-the-shell-cannot-parse-is-never-honest.md`
- `docs/adr/0227-a-review-through-claude-keeps-its-verdict-and-fits-its-window.md`
- `docs/adr/0228-an-unrecorded-branch-prefix-follows-the-commit-type-rule.md`
- `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`
- `docs/adr/0230-an-archived-spec-keeps-its-relative-links.md`
- `docs/adr/0231-the-jev-ceiling-is-a-user-config-value.md`
- `docs/adr/0232-merge-evidence-releases-the-runs-and-item-branches-of-a-merged-spec.md`
- `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`
- `docs/adr/0235-openai-and-anthropic-models-run-only-through-the-codex-and-claude-subscriptions.md`
- `docs/adr/0236-a-failed-check-is-judged-on-a-merge-with-the-current-default-branch.md`
- `docs/adr/0237-a-delivery-retry-records-a-merge-made-outside-the-queue.md`
- `docs/adr/0238-a-light-tier-runs-low-complexity-tasks-on-an-open-model.md`
- `docs/adr/0239-each-openrouter-stage-reads-its-own-roundfix-key-first.md`
- `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- `docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md`
- `docs/adr/0242-a-blocked-review-names-the-protocol-step-and-retries-once-before-the-prompt.md`
- `docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`
- `docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md`
- `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- `docs/adr/0249-a-qa-step-formats-its-qa-directory-and-reopen-sees-a-late-dependency.md`
- `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`
- `docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`
- `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`
- `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`
- `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- `docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`
- `docs/adr/0258-a-spec-is-authored-against-the-qa-rerun-classes-and-a-retirement-writes-reduced-history.md`
- `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- `docs/adr/0263-the-better-auth-base-path-is-a-repository-decision-a-baseline-update-records.md`
- `docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md`
- `docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md`
- `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`
- `docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`
- `docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md`
- `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- `docs/adr/0274-a-turn-that-did-not-finish-is-not-accepted.md`
- `docs/adr/0276-a-row-the-authorization-keeps-offline-qualifies-and-a-byte-invariant-compares-with-the-delivery-base.md`
- `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`
- `docs/adr/0280-a-guide-task-records-its-own-surface-fingerprint-a-prd-records-who-decided-and-the-glossary-is-held-to-its-own-rule.md`
- `docs/adr/0281-a-test-replaces-an-executable-it-rebuilds-instead-of-writing-into-it.md`
- `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- `docs/adr/0283-a-run-marker-keeps-the-judge-silent-inside-a-run.md`
- `docs/adr/0284-carry-forward-binds-a-task-to-the-planning-text-it-references.md`
- `docs/adr/0285-every-test-file-compiles-for-windows-and-a-repository-contract-proves-it.md`
- `docs/adr/0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md`
- `docs/adr/0287-a-plan-is-checked-for-undeclared-verification-paths-and-every-settlement-checks-the-setup-manifest.md`
- `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`
- `docs/adr/0289-a-delivered-specs-dirty-run-is-superseded-under-its-declared-directories-and-a-kept-one-names-why.md`
- `docs/adr/0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md`
- `docs/adr/0291-qa-evidence-holds-text-and-images-and-the-qa-report-commit-refuses-any-other-binary.md`
- `docs/adr/0292-a-task-commit-is-scanned-for-secrets-and-a-retry-accepts-a-secret-rewrite.md`
- `docs/adr/0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md`
- `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- `docs/adr/0295-reachability-reads-the-fetched-default-branch-and-origin-names-it-when-origin-head-is-unset.md`
- `docs/adr/0296-the-storage-report-reconciles-every-page-of-one-storage-snapshot.md`
- `docs/adr/0297-the-linux-lane-reads-the-repository-history-read-only.md`
- `docs/adr/0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md`
- `docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md`
- `docs/adr/0301-gc-sanitation-classifies-the-artifact-root-each-run-wrote-to.md`
- `docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md`
- `docs/adr/0303-a-linux-lane-skipped-for-its-environment-is-shown-where-the-delivery-is-read.md`
- `docs/adr/0304-a-queue-item-waits-for-a-run-slot-up-to-runs-slot-wait.md`
- `docs/adr/0306-archive-and-cancel-refuse-a-path-whose-resolved-target-leaves-the-repository.md`
- `docs/adr/0308-an-archive-rebases-the-unproven-links-it-records-and-the-record-link-repair-points-a-promoted-target-at-its-copy.md`
- `docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md`
- `docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md`
- `docs/adr/0312-a-new-delivery-queue-resumes-an-archived-item-branch-at-review.md`
- `docs/adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md`
- `docs/adr/0314-a-light-tier-that-cannot-serve-the-run-yields-to-the-category-profile.md`
- `docs/adr/0315-a-test-selector-another-argument-of-its-invocation-contains-is-refused-at-authoring.md`
- `docs/adr/0316-the-linux-lane-names-a-failed-go-list-apart-from-go-test.md`
- `docs/adr/0317-a-result-status-note-of-a-completed-task-is-dismissed-without-the-validator.md`
- `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- `docs/adr/0320-a-question-call-carries-up-to-three-independent-questions.md`
- `docs/adr/0322-a-python-profile-imports-the-published-upstream-python-preset.md`
- `docs/adr/0323-a-profile-starter-file-is-written-only-where-absent-and-only-on-a-confirmed-preview.md`
- `docs/adr/0324-a-test-never-gives-a-running-executable-a-second-name.md`
- `docs/adr/0326-a-code-test-never-reads-the-live-documentation-trees.md`
- `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
- `docs/adr/0334-every-implement-run-and-the-queue-gate-trim-a-declared-shared-build-cache.md`
- `docs/adr/0335-a-stop-that-is-not-the-work-s-fault-does-not-cost-the-run.md`
- `docs/adr/0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md`
- `docs/adr/0340-an-untracked-executable-rides-a-task-commit-only-when-the-spec-authorization-names-it.md`
- `docs/adr/0341-a-profile-change-carries-the-recorded-decisions-and-a-restore-keeps-the-claude-skill-links.md`
- `docs/adr/0342-every-task-settles-on-the-repository-verification-and-a-run-without-a-qa-gate-checks-its-integrated-tree.md`
- `docs/adr/0343-settle-follows-the-branch-plan-and-a-declared-row-matches-the-criterion-it-names.md`
- `docs/adr/0344-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails-and-the-next-proof-closes-it.md`
- `docs/adr/0345-a-derived-record-writer-names-its-command-to-the-suite-guard.md`
- `docs/adr/0346-a-run-agent-loads-the-skills-the-dispatch-selects-for-its-task.md`
- `docs/adr/0347-a-row-citing-an-empty-grant-declares-no-mutation-and-a-completed-task-never-collides.md`
- `docs/adr/0348-an-absent-spec-root-inside-the-repository-holds-no-spec.md`
- `docs/adr/0349-a-spec-release-ceiling-answers-the-release-plan-approval.md`
- `docs/adr/0350-an-approved-baseline-update-deletes-the-retired-skills-it-proves-unedited.md`
- `docs/adr/0351-a-qa-report-settles-only-when-its-local-links-resolve-in-the-committed-tree.md`
- `docs/adr/0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md`
- `docs/adr/0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md`
- `docs/adr/0354-the-spec-skills-announce-spec-events-only-through-the-users-own-instructions.md`
- `docs/adr/0355-a-precondition-refusal-is-never-imported-as-a-prior-qa-pass.md`
- `docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md`
- `docs/adr/0359-an-installed-owned-skill-keeps-an-unpinned-lock-entry-and-doctor-names-a-stale-one.md`
- `docs/adr/0360-the-vendored-catalog-skills-follow-the-catalog-through-one-pinned-commit.md`
- `docs/adr/0361-a-task-commit-secret-match-returns-as-verification-feedback-and-names-its-cause.md`
- `docs/adr/0362-only-a-git-entry-that-git-accepts-marks-a-git-root.md`
- `docs/adr/0363-runs-list-prints-a-table-at-a-terminal-and-a-short-run-id-names-a-run.md`
- `docs/adr/0364-a-required-external-skill-is-checked-only-for-presence.md`
- `docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md`
- `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`
- `docs/adr/0367-a-numbered-file-collision-with-an-earlier-queue-item-parks-before-the-run.md`
- `docs/adr/0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`
- `docs/adr/0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md`
- `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`
- `docs/adr/0371-a-session-starts-with-an-update-check-and-the-skill-applies-the-update.md`
- `docs/adr/0372-the-baseline-names-the-archive-destination-a-profiles-missing-decision-and-a-guide-it-does-not-own.md`
- `docs/adr/0373-a-proven-owner-identity-mismatch-means-the-owner-is-gone-for-every-liveness-reader.md`
- `docs/adr/0374-a-rerun-completes-a-release-set-from-the-same-content.md`
- `docs/adr/0375-a-grouped-select-option-projects-its-groups-values-and-capability-evidence-stays-fail-closed.md`
- `docs/adr/0376-a-qa-report-relies-only-on-evidence-its-own-run-produced.md`
- `docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md`
- `docs/adr/0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md`
- `docs/adr/0380-a-skill-the-catalog-replaces-is-followed-to-its-replacements-and-retired.md`
- `docs/adr/0381-a-results-row-outside-every-results-table-refuses-the-qa-report.md`
- `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- `docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md`
- `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md`
- `docs/adr/0386-a-user-level-copy-of-a-roundfix-owned-skill-older-than-the-binary-is-reported-never-changed.md`
- `docs/adr/0387-the-suite-time-budget-is-measured-on-ci-and-judged-after-the-contract-runs.md`
- `docs/adr/0389-an-adr-alternatives-line-marks-an-adr-as-not-applying-in-english-or-portuguese.md`
- `docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md`
- `docs/adr/0391-a-conflicting-merge-of-the-default-branch-is-a-review-delta-from-its-conflicted-carry.md`
- `docs/adr/0392-an-older-roundfix-refuses-a-repository-whose-baseline-a-newer-roundfix-wrote.md`
- `docs/adr/0393-roundfix-reads-skills-lock-versions-1-and-2-and-keeps-the-version-it-found.md`
- `docs/adr/0394-the-makefile-follows-the-go-build-cache-set-with-go-env-and-falls-back-small.md`
- `docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md`
- `docs/adr/0396-a-message-that-asks-to-be-told-when-its-work-is-done-receives-one-completion-notice.md`
- `docs/adr/0398-a-parked-item-always-has-a-way-back-through-a-delivery-retry.md`
- `docs/adr/0399-the-baseline-inventories-a-nested-carrier-only-when-git-lists-it.md`
- `docs/adr/0400-a-spec-is-sized-by-one-context-and-by-tasks-that-earn-their-verification-cycle.md`
- `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`
- `docs/adr/0402-the-research-subagent-reads-and-never-runs-a-command-or-writes-a-file.md`
- `docs/adr/0403-a-feature-is-judged-for-the-idea-stage-before-its-prd-and-the-judgment-never-refuses.md`
- `docs/adr/0404-the-idea-stage-judge-reads-one-item-that-is-not-yet-a-spec.md`
- `docs/adr/0405-a-user-level-skill-copy-is-reported-with-the-command-that-removes-it.md`
- `docs/adr/0406-the-authoring-verification-probe-prepares-its-checkout-as-the-daemon-prepares-a-task-worktree.md`
- `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`
- `docs/adr/0408-a-repository-that-publishes-an-authoring-skill-receives-no-header-for-it-and-a-fix-is-authored-in-the-session.md`
- `docs/adr/0409-a-queue-forecast-is-a-likely-time-and-a-limit-from-the-repository-own-deliveries-or-a-named-condition.md`
- `docs/adr/0410-the-status-view-lists-specs-being-authored-and-open-backlog-entries.md`
- `docs/adr/0411-a-repository-rule-outside-the-specific-repository-guide-is-named-never-moved.md`
- `docs/adr/0414-a-supervisor-session-ends-only-when-the-maintainer-asks.md`
- `docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md`
- `docs/adr/0416-the-update-check-names-each-configured-profile-that-differs-from-the-recommended-profile.md`
- `docs/adr/0417-a-correction-carries-a-declared-defect-class-and-a-repository-requires-it-only-by-its-own-setting.md`
- `docs/adr/0418-a-phase-entry-holds-the-model-and-reasoning-effort-that-ran-an-authoring-phase.md`
- `docs/adr/0419-the-corrections-report-prints-counts-and-sets-two-models-side-by-side-only-above-twenty-specs-each.md`
- `docs/adr/0420-the-corrections-record-is-one-append-only-file-per-spec-under-roundfix-home.md`
- `docs/adr/0421-a-phase-entry-and-the-corrections-record-name-their-repository-by-its-git-lineage.md`
- `docs/adr/0422-a-commit-without-a-mark-of-a-run-or-of-the-queue-on-the-item-branch-is-a-correction-sign.md`
- `docs/adr/0423-an-undeclared-correction-sign-parks-the-item-at-its-merge-as-a-finding.md`
- `docs/adr/0424-corrections-are-declared-and-reported-by-their-own-command-group-and-named-c-n.md`
- `docs/adr/0425-a-session-passes-its-model-to-the-phase-command-and-only-the-name-its-harness-states.md`
- `docs/adr/0426-the-corrections-report-counts-the-tasks-the-plan-held-on-the-default-branch-and-starts-no-fetch.md`
- `docs/adr/0427-an-authoring-phase-falls-back-to-its-declared-fallback-when-the-declared-model-s-usage-limit-runs-out.md`
- `docs/adr/0428-a-spec-whose-phase-fell-back-keeps-the-fallback-for-the-limited-model-until-the-spec-ends.md`
- `docs/adr/0429-roundfix-reads-the-limit-from-the-error-a-session-passes-and-prints-the-model-an-author-subagent-starts-on.md`
- `docs/adr/0430-the-recommended-authoring-declaration-carries-a-fallback-for-each-authoring-phase.md`
- `docs/adr/0431-the-baseline-tells-each-root-instruction-file-how-an-authoring-phase-keeps-its-model.md`
- `docs/adr/0432-claude-md-is-always-a-symbolic-link-to-agents-md.md`
- `docs/adr/0433-a-checkout-without-symbolic-links-keeps-the-link-git-records-for-claude-md.md`
- `docs/adr/0434-a-renamed-spec-keeps-its-spec-number-and-is-one-row-of-the-spec-queue.md`
- `docs/adr/0435-a-delivery-pull-request-behind-its-base-branch-is-brought-up-to-date-by-a-merge-of-the-default-branch.md`
- `docs/adr/0436-an-authoring-phase-questions-a-rule-that-blocks-the-goal-and-the-rule-binds-until-the-user-answers.md`
- `docs/adr/0437-a-rule-the-user-keeps-carries-a-kept-line-and-is-asked-about-again-only-with-a-new-fact.md`
- `docs/adr/0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md`
- `docs/adr/0439-an-adrs-header-says-whether-the-secondbrain-wants-it-and-which-adrs-amend-it.md`
- `docs/adr/0440-a-lesson-learned-is-written-only-when-a-test-holds-and-leaves-when-a-rule-states-it.md`
- `docs/adr/0441-the-recommended-authoring-declaration-gives-the-research-subagent-the-built-in-web-tools.md`
- `docs/adr/0442-an-open-authoring-question-waits-unless-the-project-declares-that-the-session-follows-the-author.md`
- `docs/adr/0447-in-an-adopter-a-change-to-the-files-roundfix-manages-or-to-a-context-md-needs-no-release.md`
