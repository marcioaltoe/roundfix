---
task: task_04
spec: 0420-roundfix-stops-verifying-on-linux
status: pending
type: docs
complexity: medium
---

# Task 04: ADR-0452 supersedes the lane's records, and the glossary and export list follow

## Overview

This Task delivers `_techspec.md` → Build Order 4: the rules change with the work, as the maintainer's decision of 2026-11-24, "Tirar de vez", retires them. It writes the accepted ADR-0452, deletes the five decision records of the lane it supersedes, adds the amendment notices to the six records that change in part, removes the four glossary entries of the lane and moves the citations of Base Test Tools and External Test Tool to ADR-0452, and updates the Standard Export List. Its documentation contract also sweeps the live documents for any trace of the lane that task_01 to task_03 left. It is verifiable on its own through two documentation contract tests, the glossary checks and the existing documentation contracts.

## Requirements

1. MUST create `docs/adr/0452-roundfix-verifies-on-macos-and-the-delivery-queue-has-no-linux-lane.md`, accepted, with a true Knowledge Mark, `created_at` and `updated_at` the date of the change, and the title `Roundfix verifies on macOS, and the Delivery Queue has no Linux Lane`, as `_techspec.md` → API Contract 11 states:
   - after the title it names what it supersedes and what it amends, each with the date of its `created_at`, beginning `Supersedes ADR-0288, ADR-0294, ADR-0297, ADR-0303 and ADR-0316 (<date>); all are retired, and Git keeps their text.`;
   - it states as rules API Contracts 1 to 8 and the first item of API Contract 9, and that the suite-time budget keeps its number until the macOS runner measures the suite;
   - it restates the rule of the tests' tools with its words: a test that needs an External Test Tool skips naming it, a tool outside the Base Test Tools is resolved through `testfixture.RequireTool`, and the rule in `internal/testfixture` refuses any other way;
   - it records the options given up: a light Linux job that only builds and vets, no test in continuous integration, the whole release on the macOS runner, three release jobs that validate the tag first, the moving label `macos-latest`, a larger runner, and a class of `unclassified` for the stored blocker;
   - it names no client and cites no Spec folder path.
2. MUST delete, in the same change, the five superseded records: `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`, `docs/adr/0294-the-linux-lane-runs-like-ci.md`, `docs/adr/0297-the-linux-lane-reads-the-repository-history-read-only.md`, `docs/adr/0303-a-linux-lane-skipped-for-its-environment-is-shown-where-the-delivery-is-read.md` and `docs/adr/0316-the-linux-lane-names-a-failed-go-list-apart-from-go-test.md` (`docs/agents/docs-layout.md`: "a superseded ADR once an accepted successor names it with its date").
3. MUST give each of ADR-0252, ADR-0285, ADR-0295, ADR-0384, ADR-0385 and ADR-0423 `amended_by` holding `ADR-0452` in its header (added, or appended to a list it already has), the date of the change in `updated_at`, and one paragraph that begins `**Amended by [ADR-0452](0452-roundfix-verifies-on-macos-and-the-delivery-queue-has-no-linux-lane.md) (<date>).**` after the text `_techspec.md` → ADR Alternatives quotes from it, saying in one or two sentences what changed: the gate and continuous integration run on macOS (ADR-0252), the contract compiles the tests for Linux as well (ADR-0285), the example's lane is retired (ADR-0295), the gate's unit covers the repository Verification alone (ADR-0384), the lane's container no longer receives the budget (ADR-0385), and the lane's key is retired (ADR-0423). Every other byte of those six files stays, so every sentence this Spec quotes from them, the plain-text `(ADR-0288)` of ADR-0423 included, stays where it is (invariant 12). Each stays `accepted`.
4. MUST revise the glossary through `domain-modeling`: remove from `CONTEXT.md` the entries **Linux Lane**, **Lane Environment Skip**, **Lane Container Failure** and **Lane Change Failure** whole, and change **Base Test Tools** and **External Test Tool** only by replacing the link to ADR-0294 with a link to `docs/adr/0452-roundfix-verifies-on-macos-and-the-delivery-queue-has-no-linux-lane.md` as ADR-0452, keeping their definitions and their `_Avoid_` lines (`_prd.md` → Glossary). No other entry of `CONTEXT.md` changes.
5. MUST remove from `.secondbrain-export` the lines of ADR-0288, ADR-0303 and ADR-0316, and add by hand the path of ADR-0452 in its place by number, after the last line of a decision record (`docs/agents/secondbrain.md`: "the only line added by hand is the path of an ADR, in order"). No other line changes.
6. MUST add `internal/docscontract/linux_verification_decision_record_test.go`, with the `docscontract` build tag and the `verify:always` marker the package's other files carry, holding:
   - `TestTheLinuxVerificationDecisionSupersedesTheLaneRecords`, which finds the record by its title, never by its number, takes its `ADR-NNNN` from its file name, and asserts in the subtests `decision_record` (status `accepted`, Knowledge Mark true, the supersession sentence naming the five records with dates), `superseded_records_are_deleted` (no file under `docs/adr/` begins with `0288-`, `0294-`, `0297-`, `0303-` or `0316-`), `amended_records` (each of the six is `accepted`, its `amended_by` holds the record's id and its body holds one `**Amended by [` paragraph linking the record's file), `glossary` (none of the four removed terms is defined and the two changed entries link the record's file) and `export_list` (the record's path is listed once, in ascending order, and none of the five deleted paths is);
   - `TestNoLiveDocumentDescribesTheLinuxLane`, which reads, line by line, every Markdown file under `docs/user-guide/`, `.agents/skills/roundfix/`, `skills/roundfix/` and `docs/agents/`, and `AGENTS.md`, `CLAUDE.md`, `README.md`, the `Makefile`, every file under `.github/workflows/` and `.roundfixrc.yml`, and fails naming the file and line for any line that holds `Linux Lane`, `deliver lane`, `Linux gate`, `linux-lane` or `DR-LANE-`, except a line that holds `is deprecated and ignored; the Linux Lane was retired, so remove it` (`_prd.md` → Success Metric 8).
7. MUST raise by one, from the value it finds, the number of `docscontract` files that `TestRepositoryContractTestsDeclareTheirRelevance` pins in `internal/verifyselect/contracts_repository_test.go`, for the file of requirement 6; change nothing else in that file.
8. MUST NOT change any other decision record, any guide under `docs/agents/`, any skill, any user-guide page, any Go source other than the two test files named above, or any file `_authorization.md` → Limits names. When the sweep of requirement 6 finds a line task_01 to task_03 left, stop and report it rather than editing a file this Task does not declare.

## Subtasks

- [ ] Write ADR-0452 and delete the five records it supersedes.
- [ ] Add the six amendment notices and headers.
- [ ] Remove the four glossary entries and move the two citations.
- [ ] Update the Standard Export List.
- [ ] Add the decision record and sweep tests, and raise the pinned count.

## Acceptance Criteria

- [ ] ADR-0452 is accepted, carries a true Knowledge Mark and supersedes the five records, which are absent from `docs/adr/`.
- [ ] The six amended records name ADR-0452 in `amended_by` with a dated notice, and keep every other byte of their bodies.
- [ ] The glossary defines none of the four lane terms; Base Test Tools and External Test Tool keep their definitions and cite ADR-0452.
- [ ] The Standard Export List lists ADR-0452 in order and none of the deleted records.
- [ ] No live document describes the lane outside the two warning lines; the documentation contracts pass.

## Context

- instruction: `internal/docscontract/delivery_recovery_decision_record_test.go`
- instruction: `internal/docscontract/secondbrain_export_test.go`
- instruction: `internal/testfixture/external_tool_test.go`
- interface: `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`
- interface: `docs/adr/0285-every-test-file-compiles-for-windows-and-a-repository-contract-proves-it.md`
- interface: `docs/adr/0295-reachability-reads-the-fetched-default-branch-and-origin-names-it-when-origin-head-is-unset.md`
- interface: `docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md`
- interface: `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md`
- interface: `docs/adr/0423-an-undeclared-correction-sign-parks-the-item-at-its-merge-as-a-finding.md`
- interface: `CONTEXT.md`
- interface: `.secondbrain-export`
- interface: `internal/verifyselect/contracts_repository_test.go`
- creates: `docs/adr/0452-roundfix-verifies-on-macos-and-the-delivery-queue-has-no-linux-lane.md`
- creates: `internal/docscontract/linux_verification_decision_record_test.go`
- deletes: `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`
- deletes: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- deletes: `docs/adr/0297-the-linux-lane-reads-the-repository-history-read-only.md`
- deletes: `docs/adr/0303-a-linux-lane-skipped-for-its-environment-is-shown-where-the-delivery-is-read.md`
- deletes: `docs/adr/0316-the-linux-lane-names-a-failed-go-list-apart-from-go-test.md`

## Consumers

- this repository's decision records, the glossary and the Standard Export List — consumers: `./internal/docscontract`
- the pinned number of `docscontract` files — consumers: `./internal/verifyselect`

## Verification

- `out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestTheLinuxVerificationDecisionSupersedesTheLaneRecords|TestNoLiveDocumentDescribesTheLinuxLane)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheLinuxVerificationDecisionSupersedesTheLaneRecords/decision_record TestTheLinuxVerificationDecisionSupersedesTheLaneRecords/superseded_records_are_deleted TestTheLinuxVerificationDecisionSupersedesTheLaneRecords/amended_records TestTheLinuxVerificationDecisionSupersedesTheLaneRecords/glossary TestTheLinuxVerificationDecisionSupersedesTheLaneRecords/export_list TestTheLinuxVerificationDecisionSupersedesTheLaneRecords TestNoLiveDocumentDescribesTheLinuxLane; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, ADR-0452 supersedes and amends as stated and no live document describes the lane.
- `grep -F -- '**Base Test Tools**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Base Test Tools\*\*(:| —( |$))' && grep -F -- '**External Test Tool**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*External Test Tool\*\*(:| —( |$))' && for term in 'Base Test Tools' 'External Test Tool'; do awk -v header="**$term**:" '$0 == header { found=1; next } found && $0 == "" { exit } found { print }' CONTEXT.md | grep -qF -- 'roundfix-verifies-on-macos-and-the-delivery-queue-has-no-linux-lane' || { echo "entry does not cite the new record: $term" >&2; exit 1; }; done && ! grep -qE '^\*\*(Linux Lane|Lane Environment Skip|Lane Container Failure|Lane Change Failure)\*\*' CONTEXT.md` — expected: exit 0, the two changed terms keep their definitions and cite the new record, and the four lane terms are gone.
- `f="$(grep -lxF '# Roundfix verifies on macOS, and the Delivery Queue has no Linux Lane' docs/adr/*.md)" || exit 1; test -n "$f" && grep -qx 'status: accepted' "$f" && grep -qx 'knowledge: true' "$f" && grep -qxF -- "$f" .secondbrain-export && for n in 0288 0294 0297 0303 0316; do if ls docs/adr | grep -q "^$n-"; then echo "superseded record still present: ADR-$n" >&2; exit 1; fi; if grep -q "^docs/adr/$n-" .secondbrain-export; then echo "export list still names $n" >&2; exit 1; fi; done && for n in 0252 0285 0295 0384 0385 0423; do grep -q '^amended_by:.*ADR-0452' docs/adr/$n-*.md || { echo "ADR-$n does not name ADR-0452" >&2; exit 1; }; done && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 20m -run '^TestRepositoryContractTestsDeclareTheirRelevance$' ./internal/verifyselect` — expected: exit 0, the record is accepted, marked and exported, the five superseded records and their export lines are gone, the six amended records name it, and the documentation contracts and the pinned count pass.

## References

- `_prd.md` → Goal 5; User Story 6; Core Feature 11; Success Metric 8; Glossary; Rule challenges; Decisions
- `_techspec.md` → The decision record, the glossary, the guide and the skill; API Contract 11; API Contract 13; Invariants; ADR Alternatives; Testing Approach; Build Order 4; Risks & Considerations; Coverage Map (Deviations from the PRD)
- [`references/2026-11-24-techspec-measurements.md`](references/2026-11-24-techspec-measurements.md), sections 2 and 5
- ADR-0288, ADR-0294, ADR-0297, ADR-0303, ADR-0316, ADR-0252, ADR-0285, ADR-0295, ADR-0384, ADR-0385, ADR-0423
