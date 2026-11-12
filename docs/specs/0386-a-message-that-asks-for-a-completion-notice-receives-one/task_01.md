---
task: task_01
spec: 0386-a-message-that-asks-for-a-completion-notice-receives-one
status: pending
type: docs
complexity: medium
---

# Task 01: The Baseline sends one completion notice to a message that asks for it

## Overview

The sender's and the receiving clauses of the core module's cross-project
continuity rule keep their ids and change their words: a message that asks
to be told when its work is done receives one completion notice once the
release that carries that work is published, or once it merges when it
releases nothing; the request is kept in the record of the message; a
Daemon-assigned turn never sends the notice; and a sender asks for it only
when it will act on it, without waiting or watching. Every adopter's managed
refresh writes the new words into its agent instructions guide, and this
repository's own guide and setup manifest are refreshed the same way.

## Requirements

1. MUST change, in `internal/baseline/assets/modules/core.json`, only the
   `guidance` of `clause.core.message-the-owning-project` and
   `clause.core.consume-messages-after-recording` to the texts of
   `_techspec.md` → API Contract 1, word for word, and raise
   `rule.core.cross-project-continuity`'s `version` from `3` to `4`. MUST keep
   each clause's id, `enforcement` and `replaces` list,
   `clause.core.messages-are-self-contained`, the file's two-space JSON
   formatting, and every other clause, rule and module unchanged; MUST NOT add
   or remove a clause.
2. MUST add to `internal/baseline/message_reply_clause_test.go`
   `TestTheReceivingClauseSendsOneCompletionNoticeOnlyWhenAMessageAsksForIt`,
   calling `t.Parallel()` first, with the six phrases and the absent phrase
   `_techspec.md` → Testing Approach lists for `internal/baseline`. MUST NOT
   change `TestTheMessageClausesAskForAReplyOnlyForBlockedWork` or
   `TestTheReceivingClauseAnswersQuestionsAndBoundsTheDeclaredPrecedence`;
   both must keep passing, because the new wording keeps every phrase they
   pin once.
3. MUST regenerate, in this order and never by hand: the module record with
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   (it writes the next version of `core`; the Spec names no number); the
   derived Baseline files with `make baseline-digests`; and, as the Task's
   last step, this repository's managed guide and setup manifest with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
   whose second run reports no change. Specs 0383, 0384 and 0385 deliver
   first, so the derived files are measured on the tree the Task runs on: the
   files the prototype rewrote are declared below, and any other file those
   commands rewrite is recorded by the Daemon.
4. MUST NOT change any other module, template or profile source by hand,
   `.roundfixrc.yml`, or a Source Baseline. Of the Governed Paths the Task may
   change only `internal/baseline/assets/modules/core.json`,
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

- [ ] A resolved profile's agent instructions give one completion notice to a message that asks for it, after the release or the merge, never from a Daemon-assigned turn, and keep the old prohibition for every other message.
- [ ] The three message clause ids keep their force and Source Baseline replacements, and the reply rule and bounded precedence tests pass unchanged.
- [ ] The derived Baseline files and this repository's guide and setup manifest are current, and a second managed refresh reports no change.

## Context

- instruction: `docs/adr/0396-a-message-that-asks-to-be-told-when-its-work-is-done-receives-one-completion-notice.md`
- instruction: `docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md`
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

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheReceivingClauseSendsOneCompletionNoticeOnlyWhenAMessageAsksForIt|TestTheMessageClausesAskForAReplyOnlyForBlockedWork|TestTheReceivingClauseAnswersQuestionsAndBoundsTheDeclaredPrecedence|TestEveryInboxClauseAnAdopterHoldsIsReplaced|TestTheGuidesSendCrossProjectWorkThroughAiMemory|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheReceivingClauseSendsOneCompletionNoticeOnlyWhenAMessageAsksForIt ' || { echo 'missing pass: TestTheReceivingClauseSendsOneCompletionNoticeOnlyWhenAMessageAsksForIt' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 30m -run 'Baseline|Guidance|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the clause test does not exist, so its pass line is missing; afterwards the clauses render, the earlier reply and precedence tests still pass, the replacements and the module version are recorded, and the Baseline suite and the CLI's Baseline tests pass.

- `tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- 'unless it asks to be told when its work is done' || { echo 'agent instructions lack the exception' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- 'a Daemon-assigned turn never sends one' || { echo 'agent instructions lack the sender of the notice' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- 'A message may also ask to be told when the work it asks for is done' || { echo 'agent instructions lack the sender clause' >&2; exit 1; }; f="$(mktemp)"; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json > "$f"; status=$?; grep -q '"state":"current"' "$f" || { cat "$f"; rm -f "$f"; exit 1; }; rm -f "$f"; test "$status" -eq 0 && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide has neither the exception nor the sender's request; afterwards this repository's guide carries both, a managed refresh reports `current`, and the documentation contracts, the setup manifest included, pass.

## References

- `_prd.md` → Goals 1, 2, 4 and 5; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior 1 and 2; API Contract 1; Invariants 2 to 4; Declared breaks; Testing Approach; Build Order 1
- ADR-0396
