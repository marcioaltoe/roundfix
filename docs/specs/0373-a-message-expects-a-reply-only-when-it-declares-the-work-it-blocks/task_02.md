---
task: task_02
spec: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
status: pending
type: docs
complexity: medium
---

# Task 02: The Baseline's message clauses ask for a reply only when the sender's work is blocked

## Overview

The three message clauses of the core module's cross-project continuity rule
keep their ids and change their words: a message is one-way unless the
sender's own work depends on the answer, a sender that depends on it says so
in the subject and first lines and names the blocked work, and the receiver
handles such a message next without interrupting running work, answers only
it, and sends no notice to the rest. Every adopter's managed refresh writes
the new words into its agent instructions guide, and this repository's own
guide and setup manifest are refreshed the same way.

## Requirements

1. MUST change, in `internal/baseline/assets/modules/core.json`, only the
   `guidance` of `clause.core.message-the-owning-project`,
   `clause.core.messages-are-self-contained` and
   `clause.core.consume-messages-after-recording` to the texts of
   `_techspec.md` → API Contract 8, word for word, and raise
   `rule.core.cross-project-continuity`'s `version` from `2` to `3`. MUST keep
   each clause's id, `enforcement` and `replaces` list, the file's two-space
   JSON formatting, and every other clause, rule and module unchanged; MUST
   NOT add or remove a clause.
2. MUST create `internal/baseline/message_reply_clause_test.go` with
   `TestTheMessageClausesAskForAReplyOnlyForBlockedWork`, calling
   `t.Parallel()` first, built like
   `TestTheGuidesSendCrossProjectWorkThroughAiMemory`
   (`buildProjectDecisionPlan`, `newProjectDecisionPlanRepository`,
   `standardTypeScriptDecisions("make verify")`, `planPostimage`), asserting
   that the planned `docs/agents/agent-instructions.md`, whitespace-normalized,
   contains exactly once each of "A message expects no reply unless the
   sending project's own work depends on the answer.", "say so in the
   message's subject and first lines, in the message's own words and
   language", "after the work already running, which is never interrupted for
   it" and "send no acknowledgment, adoption notice or completion notice to a
   message that declared none", and contains neither "triage them oldest
   first" nor "when the message asks for a reply". No existing test changes.
3. MUST regenerate, in this order and never by hand: the module record with
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   (it writes the next version of `core`; the Spec names no number); the
   derived Baseline files with `make baseline-digests`; and, as the Task's
   last step, this repository's managed guide and setup manifest with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
   whose second run reports no change. Specs 0372 and 0371 deliver first, so
   the derived files are measured on the tree the Task runs on: the files the
   scratch prototype rewrote are declared below, and any other file those
   commands rewrite is recorded by the Daemon.
4. MUST NOT change any other module, template or profile source by hand,
   `.roundfixrc.yml`, or a Source Baseline. Of the Governed Paths the Task
   may change only `internal/baseline/assets/modules/core.json`,
   `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
   `docs/agents/agent-instructions.md` and `docs/agents/setup-context.json`,
   all through the commands above.

## Subtasks

- [ ] Reword the three clauses and raise the rule version.
- [ ] Add the clause test.
- [ ] Record the module version and regenerate the derived files.
- [ ] Refresh this repository's guide and setup manifest, and confirm a second refresh changes nothing.

## Acceptance Criteria

- [ ] A resolved profile's agent instructions carry the sender, content and receiver rules once each, and no longer say "triage them oldest first" or "when the message asks for a reply".
- [ ] The three clause ids keep their force and Source Baseline replacements.
- [ ] The derived Baseline files and this repository's guide and setup manifest are current, and a second managed refresh reports no change.

## Context

- instruction: `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- instruction: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
- instruction: `internal/baseline/cross_project_continuity_test.go`
- instruction: `internal/baseline/retirement_rule_clauses_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/module-versions.json`
- creates: `internal/baseline/message_reply_clause_test.go`
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

- The core module's three clauses and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`
- This repository's managed guide and setup manifest — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheMessageClausesAskForAReplyOnlyForBlockedWork|TestEveryInboxClauseAnAdopterHoldsIsReplaced|TestTheGuidesSendCrossProjectWorkThroughAiMemory|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheMessageClausesAskForAReplyOnlyForBlockedWork ' || { echo 'missing pass: TestTheMessageClausesAskForAReplyOnlyForBlockedWork' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 30m -run 'Baseline|Guidance|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the clause test does not exist; afterwards the clauses render, their replacements and the module version are recorded, and the Baseline suite and the CLI's Baseline tests pass.

- `tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qE -- 'A message expects no reply unless the sending project.s own work depends on the answer' || { echo 'agent instructions lack the sender rule' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qE -- 'after the work already running, which is never interrupted for it' || { echo 'agent instructions lack the receiver rule' >&2; exit 1; }; f="$(mktemp)"; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json > "$f"; status=$?; grep -q '"state":"current"' "$f" || { cat "$f"; rm -f "$f"; exit 1; }; rm -f "$f"; test "$status" -eq 0 && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide has neither rule; afterwards this repository's guide carries both, a managed refresh reports `current`, and the documentation contracts, the setup manifest included, pass.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; User Story 3; Core Feature 1; Core Feature 2; Core Feature 3; Core Feature 4; Success Metric 1
- `_techspec.md` → Current behavior 5; API Contract 8; Invariants 4 and 5; Declared breaks; Testing Approach; Build Order 2
- ADR-0383
