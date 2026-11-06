---
task: task_04
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
status: pending
type: docs
complexity: medium
---

# Task 04: The Baseline tells each session to run the Update Check before its first task

## Overview

The Baseline's core module gains one mandatory clause in its cross-project
continuity rule: at the start of each session, before its first task, run
`roundfix update check` once and relay it; in `auto`, follow its `Next:` line
through the Roundfix skill's `--update` mode; in a Daemon-assigned turn, skip
it. Every adopter's managed refresh writes it into its agent instructions
guide, and this repository's own guide and setup manifest are refreshed the
same way.

## Requirements

1. MUST add to `internal/baseline/assets/modules/core.json`, in
   `rule.core.cross-project-continuity`, after
   `clause.core.consume-messages-after-recording`, the clause
   `clause.core.check-for-updates-at-session-start` with enforcement
   `mandatory`, no `replaces`, and the guidance of `_techspec.md` → API
   Contract 9 word for word, and raise that rule's `version` from `1` to `2`.
   MUST keep the file's two-space JSON formatting and change no other clause,
   rule or module.
2. MUST add `"clause.core.check-for-updates-at-session-start": "mandatory"` to
   `characterizedBaselineForce` in
   `internal/baseline/clause_characterization_test.go`, gofmt-aligned, and add
   `TestTheSessionStartClauseRunsTheUpdateCheck` in
   `internal/baseline/update_check_clause_test.go`, built like
   `TestTheGuidesSendCrossProjectWorkThroughAiMemory`
   (`buildProjectDecisionPlan`, `planPostimage`), asserting that the planned
   `docs/agents/agent-instructions.md`, whitespace-normalized, contains
   "run `roundfix update check` once and relay its report",
   "With `update.mode` set to `auto`, run that mode before the first task" and
   "In a Daemon-assigned turn, skip the check.", each exactly once.
3. MUST regenerate, in this order and never by hand: the module record with
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   (it writes the next version of `core`; the Spec names no number); the
   derived Baseline files with `make baseline-digests`; and, as the Task's
   last step, this repository's managed guide and setup manifest with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
   whose second run reports no change. Spec 0362 delivers first and edits
   another module, so the derived files are measured on the tree the Task
   runs on: the files the scratch prototype rewrote are declared below, and
   any other file those commands rewrite is recorded by the Daemon.
4. MUST NOT write a harness hook, change any other module, template or
   profile source by hand, or change `.roundfixrc.yml`. Of the Governed Paths
   the Task may change only `internal/baseline/assets/modules/core.json`,
   `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
   `docs/agents/agent-instructions.md` and `docs/agents/setup-context.json`,
   all through the commands above.

## Subtasks

- [ ] Add the clause and raise the rule version.
- [ ] Characterize its force and add the clause test.
- [ ] Record the module version and regenerate the derived files.
- [ ] Refresh this repository's guide and setup manifest, and confirm a second refresh changes nothing.

## Acceptance Criteria

- [ ] A resolved profile's agent instructions carry the clause once, with its `notify`, `auto` and Daemon-assigned rules.
- [ ] The derived Baseline files and this repository's guide and setup manifest are current, and a second managed refresh reports no change.
- [ ] The force characterization names the clause as mandatory.

## Context

- instruction: `docs/adr/0371-a-session-starts-with-an-update-check-and-the-skill-applies-the-update.md`
- instruction: `internal/baseline/cross_project_continuity_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/clause_characterization_test.go`
- creates: `internal/baseline/update_check_clause_test.go`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/setup-context.json`

## Consumers

- The core module's clause and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`
- This repository's managed guide and setup manifest — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheSessionStartClauseRunsTheUpdateCheck|TestBaselineClauseForceIsCharacterized|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheSessionStartClauseRunsTheUpdateCheck ' || { echo 'missing pass: TestTheSessionStartClauseRunsTheUpdateCheck' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 30m -run 'Baseline|Guidance|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the clause test does not exist; afterwards the clause renders, its force and the module version are recorded, and the Baseline suite and the CLI's Baseline tests pass.

- `tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -q -- 'run .roundfix update check. once and relay its report in a few lines' || { echo 'agent instructions lack the clause' >&2; exit 1; }; grep -q 'clause.core.check-for-updates-at-session-start' internal/baseline/assets/modules/core.json || { echo 'core module lacks the clause' >&2; exit 1; }; f="$(mktemp)"; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json > "$f"; status=$?; grep -q '"state":"current"' "$f" || { cat "$f"; rm -f "$f"; exit 1; }; rm -f "$f"; test "$status" -eq 0 && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide has no such clause; afterwards this repository's guide carries it, a managed refresh reports `current`, and the documentation contracts, the setup manifest included, pass.

## References

- `_prd.md` → Goal 5; Core Feature 4; Success Metric 6; Acceptance evidence
- `_techspec.md` → Current behavior 5; API Contract 9; Declared breaks; Testing Approach; Build Order 4
- ADR-0371
