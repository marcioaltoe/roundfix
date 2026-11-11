---
task: task_02
spec: 0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered
status: pending
type: docs
complexity: medium
---

# Task 02: The Baseline's receiving clause answers every direct question and bounds the declared precedence

## Overview

The sender's and the receiving clauses of the core module's cross-project
continuity rule keep their ids and change their words: the receiving project
answers each direct question a message asks, blocking or not, and only the
oldest blocking message of each sender goes ahead of pending work, without
replacing the maintainer's judgment. Every adopter's managed refresh writes
the new words into its agent instructions guide, and this repository's own
guide and setup manifest are refreshed the same way.

## Requirements

1. MUST change, in `internal/baseline/assets/modules/core.json`, only the
   `guidance` of `clause.core.message-the-owning-project` and
   `clause.core.consume-messages-after-recording` to the texts of
   `_techspec.md` → API Contract 6, word for word, and raise
   `rule.core.cross-project-continuity`'s `version` from `2` to `3`. MUST keep
   each clause's id, `enforcement` and `replaces` list,
   `clause.core.messages-are-self-contained`, the file's two-space JSON
   formatting, and every other clause, rule and module unchanged; MUST NOT
   add or remove a clause.
2. MUST add to `internal/baseline/message_reply_clause_test.go`
   `TestTheReceivingClauseAnswersQuestionsAndBoundsTheDeclaredPrecedence`,
   calling `t.Parallel()` first, with the phrases and the two absent phrases
   `_techspec.md` → Testing Approach lists for `internal/baseline`. MUST NOT
   change `TestTheMessageClausesAskForAReplyOnlyForBlockedWork`.
3. MUST regenerate, in this order and never by hand: the module record with
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   (it writes the next version of `core`; the Spec names no number); the
   derived Baseline files with `make baseline-digests`; and, as the Task's
   last step, this repository's managed guide and setup manifest with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
   whose second run reports no change. Specs 0377, 0376, 0378 and 0379
   deliver first, so the derived files are measured on the tree the Task runs
   on: the files the scratch prototype rewrote are declared below, and any
   other file those commands rewrite is recorded by the Daemon.
4. MUST NOT change any other module, template or profile source by hand,
   `.roundfixrc.yml`, or a Source Baseline. Of the Governed Paths the Task
   may change only `internal/baseline/assets/modules/core.json`,
   `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
   `docs/agents/agent-instructions.md` and `docs/agents/setup-context.json`,
   all through the commands above.

## Subtasks

- [ ] Reword the two clauses and raise the rule version.
- [ ] Add the clause test.
- [ ] Record the module version and regenerate the derived files.
- [ ] Refresh this repository's guide and setup manifest, and confirm a second refresh changes nothing.

## Acceptance Criteria

- [ ] A resolved profile's agent instructions answer every direct question and give the precedence to one blocking message per sender, once each, and no longer carry the old triage and reply sentences.
- [ ] The three message clause ids keep their force and Source Baseline replacements.
- [ ] The derived Baseline files and this repository's guide and setup manifest are current, and a second managed refresh reports no change.

## Context

- instruction: `docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md`
- instruction: `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- instruction: `internal/baseline/cross_project_continuity_test.go`
- instruction: `internal/baseline/retirement_rule_clauses_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/message_reply_clause_test.go`
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

- The core module's two reworded clauses and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`
- This repository's managed guide and setup manifest — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheReceivingClauseAnswersQuestionsAndBoundsTheDeclaredPrecedence|TestTheMessageClausesAskForAReplyOnlyForBlockedWork|TestEveryInboxClauseAnAdopterHoldsIsReplaced|TestTheGuidesSendCrossProjectWorkThroughAiMemory|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheReceivingClauseAnswersQuestionsAndBoundsTheDeclaredPrecedence ' || { echo 'missing pass: TestTheReceivingClauseAnswersQuestionsAndBoundsTheDeclaredPrecedence' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 30m -run 'Baseline|Guidance|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the clause test does not exist; afterwards the clauses render, their replacements and the module version are recorded, and the Baseline suite and the CLI's Baseline tests pass.

- `tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qE -- 'Answer each direct question a message asks with .memory_message_send. as soon as the answer exists, whether or not the message declared blocked work' || { echo 'agent instructions lack the reply rule' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qE -- 'takes precedence for one message per sender at a time' || { echo 'agent instructions lack the bound' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qE -- 'never replaces the maintainer.s judgment' || { echo 'agent instructions lack the judgment' >&2; exit 1; }; f="$(mktemp)"; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json > "$f"; status=$?; grep -q '"state":"current"' "$f" || { cat "$f"; rm -f "$f"; exit 1; }; rm -f "$f"; test "$status" -eq 0 && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide has neither the reply rule nor the bound; afterwards this repository's guide carries both, a managed refresh reports `current`, and the documentation contracts, the setup manifest included, pass.

## References

- `_prd.md` → Goals 1, 2, 3 and 5; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Success Metric 1
- `_techspec.md` → Current behavior 5; API Contract 6; Invariant 5; Declared breaks; Testing Approach; Build Order 2
- ADR-0390
