---
task: task_04
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
status: pending
type: docs
complexity: medium
---

# Task 04: The settlement skills, write-tasks and the user guide state both acceptance rules

## Overview

The QA gate writes the Network-Prohibited Row only if its skill names it, and
an author avoids a Pinned Byte Invariant only if `write-tasks` says how. This
Task writes the `_techspec.md` → Exact texts into `qa-gate`, `archive-spec`,
`roundfix` and its archive, settle and spec references, and `write-tasks`,
raises their versions, syncs the mirrors, and updates the four user guide
pages, the adoption guide and the Behavior Surface Record.

## Requirements

1. MUST answer both adopted Backlog Entries of 2026-10-13,
   [outside evidence blocked by the Spec's own authorization does not qualify](references/2026-10-13-outside-evidence-blocked-by-the-spec-authorization-does-not-qualify.md)
   and
   [a byte invariant pinned to a past commit fails once other Specs merge first](references/2026-10-13-a-byte-invariant-pinned-to-a-past-commit-fails-after-other-merges.md),
   by writing each text of `_techspec.md` → Exact texts into the canonical
   skill files under `.agents/skills/` and the user guide pages it names.
2. MUST add the `### Rows the Spec's authorization keeps offline` section to
   `qa-gate`, `archive-spec` and `roundfix` immediately after the QA
   settlement table, before the heading that follows it today, so the
   `### QA settlement` section stays byte-identical and
   `TestSettlementGuidanceIsOneTable` passes unchanged. MUST NOT add or change
   any text inside that section.
3. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   so each of the four changed owned skills raises both version fields and the
   record lists it, and every mirror under `skills/` matches its canonical
   copy.
4. MUST add `skills/acceptance_rules_test.go` with
   `TestTheSettlementSkillsNameTheNetworkProhibitedRow`, which reads the new
   section in the three canonical settlement skills and the new sentence in
   the roundfix archive and settle references, and
   `TestWriteTasksComparesByteInvariantsWithTheDeliveryBase`, which reads the
   write-tasks section and the `SC-QA-BYTES-PINNED` bullet of the spec
   reference. Both compare whitespace-normalized text and call `t.Parallel()`
   first.
5. MUST re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
   after the user guide pages change, and never hand-edit it.
6. MUST repeat the wording sweep: search `docs/`, `skills/`, `.agents/skills/`,
   `CONTEXT.md` and `AGENTS.md`, outside `docs/history/` and this Spec's
   folder, for `network denied: <host>` and for the exempt-row lists that name
   only the Pull Request and network-denied rows, and update every live copy
   this Task owns. The two Baseline guides are task_03's.

## Subtasks

- [ ] Write the qa-gate, archive-spec, roundfix and write-tasks texts.
- [ ] Write the user guide texts.
- [ ] Sync the mirrors, record the versions and re-record the Behavior Surface Record.
- [ ] Add the skill test.

## Acceptance Criteria

- [ ] The three settlement skills carry the new section after an unchanged
      QA settlement table.
- [ ] The qa-gate skill tells the gate when to write the network-prohibited
      status and to check offline sources instead of blocking them.
- [ ] `write-tasks`, the spec reference and the spec command page name
      `SC-QA-BYTES-PINNED` and the Delivery Base.
- [ ] Every mirror matches, every version is recorded, and the Behavior
      Surface Record is current.

## Consumers

- The owned skill texts and their version record — consumers: `./skills`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0276-a-row-the-authorization-keeps-offline-qualifies-and-a-byte-invariant-compares-with-the-delivery-base.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/qa-report.md`
- interface: `docs/user-guide/commands/settle.md`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`
- creates: `skills/acceptance_rules_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheSettlementSkillsNameTheNetworkProhibitedRow|TestWriteTasksComparesByteInvariantsWithTheDeliveryBase|TestSettlementGuidanceIsOneTable|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSettlementSkillsNameTheNetworkProhibitedRow TestWriteTasksComparesByteInvariantsWithTheDeliveryBase TestSettlementGuidanceIsOneTable TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for skill in qa-gate archive-spec roundfix write-tasks; do diff -r ".agents/skills/$skill" "skills/$skill" || { printf 'mirror differs: %s\n' "$skill" >&2; exit 1; }; done; for file in docs/user-guide/commands/archive.md docs/user-guide/commands/qa-report.md docs/user-guide/commands/settle.md docs/user-guide/context-driven-development.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'blocked (environment: network prohibited by Spec authorization)' || { printf 'missing row in %s\n' "$file" >&2; exit 1; }; done; grep -qF -- 'SC-QA-BYTES-PINNED' docs/user-guide/commands/spec.md || { printf 'the spec command page lacks SC-QA-BYTES-PINNED\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract` — expected: exit 0; before this Task neither skill test exists and no user guide page names the new row, so the command fails at its first pass check.

## References

- [_prd.md](_prd.md) — Goals 2-3; User Stories 2-4; Core Features 2-4; Success Metric 5
- [_techspec.md](_techspec.md) — Exact texts; API Contract 4; Invariants 6; Testing Approach; Build Order 4
- ADR-0276; ADR-0256; ADR-0259; ADR-0271
