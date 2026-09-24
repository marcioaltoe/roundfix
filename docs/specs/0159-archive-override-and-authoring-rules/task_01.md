---
task: task_01
spec: 0159-archive-override-and-authoring-rules
status: completed
type: backend
complexity: medium
---

# Task 01: The QA Archive Override

## Overview

The Baseline guidance and the archive skill describe archiving a Spec despite failed or missing QA on maintainer authority, stamped `qa_override: true`, but `roundfix archive` accepts no override input.

## Requirements

1. MUST accept `roundfix archive <slug> --qa-override --approval <source> --reason <text>`, and refuse `--approval` or `--reason` without `--qa-override` and `--qa-override` without both.
2. MUST still require every non-QA Task completed under an override.
3. MUST refuse an override when the newest QA Report already qualifies under the declared-acceptance policy.
4. MUST stamp `qa_override: true`, `qa_override_approval`, `qa_override_reason`, `qa_override_qa_outcome` (the observed verdict, `missing`, or the read error) and `qa_override_revision` (the archived `HEAD`), and report the archive as an override.
5. MUST move the QA Task file and every QA Report byte-identically, never changing a status or a verdict.
6. MUST leave archive without the flag exactly as today.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] An authorized override archives a Spec with failed QA and stamps every field, with the QA Task and report unchanged.
- [ ] An override is refused with a non-QA Task pending, with QA already qualifying, and without approval or reason.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestArchiveQAOverrideStampsProvenance|TestArchiveQAOverrideStillRequiresNonQATasks|TestArchiveQAOverrideRefusedWhenQAQualifies|TestArchiveCommandQAOverrideRequiresApprovalAndReason)$" ./internal/spec ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestArchiveQAOverrideStampsProvenance TestArchiveQAOverrideStillRequiresNonQATasks TestArchiveQAOverrideRefusedWhenQAQualifies TestArchiveCommandQAOverrideRequiresApprovalAndReason; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the four cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The override

## Result

Implemented the authorized QA Archive Override without changing the normal
archive path. `ArchiveRequest` now accepts approval, reason and revision
provenance; the override still requires every non-QA Task completed, refuses
an already-qualifying newest QA Report, records failed, missing or unreadable
QA as observed, and stamps all five override fields. The archive result and
CLI output identify the override disposition. QA Task files and every QA
Report move unchanged.

The archive command accepts the documented positional form with
`--qa-override`, `--approval` and `--reason`. It rejects either provenance flag
without the override and rejects the override unless both non-empty values are
present. The CLI resolves the archived repository `HEAD` and records its full
revision. Invocation without the override retains the existing eligibility,
metadata and output behavior.

Focused-check evidence:

- Pre-change: `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test -count=1 -run '^TestArchiveQAOverrideStampsProvenance$' ./internal/spec` failed to compile because `ArchiveRequest.QAOverride`, `QAArchiveOverride` and `ArchiveResult.QAOverride` did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test -count=1 -run '^TestArchiveQAOverrideStampsProvenance$' ./internal/spec` passed. Its failed, missing and unreadable QA cases assert every provenance field and compare the QA Task plus all QA Reports byte-for-byte before and after the move.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test -count=1 -run '^TestArchiveQAOverride(StillRequiresNonQATasks|RefusedWhenQAQualifies)$' ./internal/spec` passed. It covers a pending non-QA Task plus pass and qualifying declared-partial reports, with refused archives leaving QA evidence unchanged.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test -count=1 -run '^Test(ArchiveCommandQAOverrideRequiresApprovalAndReason|RunArchiveQAOverrideReportsOverride)$' ./internal/cli` passed. It covers every missing/stray flag combination and the exact accepted CLI form, full `HEAD` stamp, override output, and byte-identical QA Task and Report.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test -count=1 -run '^(TestArchive|TestRunArchive|TestSpec0058|TestArchived)' ./internal/spec ./internal/cli` passed, including the existing normal archive, supersession, declared partial, external-root, help and historical-corpus cases.
- `rtk git diff --check` passed.

A broader `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test
-count=1 ./internal/spec ./internal/cli` attempt was blocked when an unrelated
CLI test path tried to reach `api.github.com`; permission to expose that
unspecified request was denied, so the attempt produced no package verdict and
was not counted as evidence.

The Daemon-owned `## Verification` command was not run in this Agent turn.
