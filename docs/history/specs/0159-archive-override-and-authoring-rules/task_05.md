---
task: task_05
spec: 0159-archive-override-and-authoring-rules
status: completed
type: backend
complexity: medium
---

# Task 05: The override waives what normal archive would refuse

## Overview

Corrective Task from the pre-PR review of 2026-09-24. The override refuses whenever the newest QA Report qualifies, without looking at the QA Task's status, so a Spec whose QA Task is failed or pending while a stale report says `pass` can be archived neither normally nor with the override, contradicting ADR-0154. And an unreadable report stores an absolute local path in `qa_override_qa_outcome`.

## Requirements

1. MUST refuse the override only when a normal archive of the same Spec would succeed: every Task completed and the newest report eligible.
2. MUST accept the override for a failed or pending QA Task whatever the newest report says, still requiring every non-QA Task completed.
3. MUST record an unreadable report's outcome with a path relative to the Spec folder, never an absolute path.
4. MUST keep the archive-spec skill, its mirror and `docs/user-guide/commands.md` true to these rules, regenerating the mirror with `make skills-sync`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A failed QA Task with a `pass` report archives under the override.
- [ ] The override is refused when a normal archive would succeed.
- [ ] The recorded outcome of an unreadable report holds no absolute path.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/archive.go`
- interface: `.agents/skills/archive-spec/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestArchiveQAOverrideAcceptsAFailedQATaskWithAPassReport|TestArchiveQAOverrideRefusedOnlyWhenNormalArchiveSucceeds|TestArchiveQAOverrideRecordsARelativeOutcome)$" ./internal/spec 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestArchiveQAOverrideAcceptsAFailedQATaskWithAPassReport TestArchiveQAOverrideRefusedOnlyWhenNormalArchiveSucceeds TestArchiveQAOverrideRecordsARelativeOutcome; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The override

## Result

The override now tests the complete normal-archive predicate before refusing:
every Task must be `completed` and the newest QA Report must be eligible. A
failed or pending QA Task therefore remains override-eligible even when the
newest report says `pass`, while every non-QA Task must still be completed.
Unreadable-report provenance re-renders the typed `QAReportError` with a path
relative to the Spec folder before storing `qa_override_qa_outcome`.

The archive skill now distinguishes normal Task completion from override Task
completion, states the same refusal boundary as the implementation, and
documents Spec-relative unreadable-report outcomes. `make skills-sync`
regenerated `skills/archive-spec/SKILL.md`, and the command guide carries the
same rules.

Acceptance evidence:

- Failed QA Task with a `pass` report: before the production edit,
  `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run '^TestArchiveQAOverrideAcceptsAFailedQATaskWithAPassReport$' ./internal/spec`
  failed because the eligible report alone refused the override. After the
  edit, the same focused check passed and confirmed that the QA Task and report
  moved byte-identically.
- Refusal only when normal archive succeeds:
  `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run '^TestArchiveQAOverrideRefusedOnlyWhenNormalArchiveSucceeds$' ./internal/spec`
  passed for both a `pass` report and a qualifying declared `partial`, with all
  Tasks completed and refused evidence left unchanged.
- Relative unreadable-report outcome: before the production edit,
  `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run '^TestArchiveQAOverrideRecordsARelativeOutcome$' ./internal/spec`
  failed because `qa_override_qa_outcome` contained the temporary Spec's
  absolute path. A strengthened broken-symlink case then exposed the same path
  inside the nested `os.PathError`; after both typed path layers were
  relativized, the focused check passed and found only
  `qa/qa-report-2026-09-24.md` in the stored diagnostic.

Additional focused-check evidence:

- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 ./internal/spec`
  passed, including the prior Task 01 override cases and the new regressions.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 ./skills`
  passed, including the identical canonical QA settlement contract.
- `rtk cmp -s .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md`
  passed after `rtk make skills-sync`.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache make verify-incremental`
  first reached the full suite but the sandbox denied process-table access to
  two unrelated force-stop integration tests. Re-running the same gate with
  process-table permission passed vet, all tests, skill synchronization and
  validation, and the build.
- `rtk git diff --check` passed.

The Daemon-owned `## Verification` command was not run in this Agent turn.
