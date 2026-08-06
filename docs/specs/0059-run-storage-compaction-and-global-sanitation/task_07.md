---
task: task_07
spec: 0059-run-storage-compaction-and-global-sanitation
status: pending
type: backend
complexity: medium
---

# Task 07: Ship the compaction command and document all three surfaces

## Overview

Corrective Task from the QA gate's F-001 (`Blocks-Completion`) and F-002.

`internal/store` carries `PreviewCompaction` and `Compact` with every guard the
Spec asked for — Active Run, other writer, insufficient temporary capacity —
and none of it is reachable. `roundfix gc compact` exits `2` with
`unexpected argument "compact"`. The TechSpec declared the surface in its API
Contracts table and no Task requirement ever asked for it, so the feature
shipped as an unreachable library.

task_05 then refused to document it, correctly: writing the Skill contract for
a command that does not exist would have been a false contract. Its Result
names this exact follow-up. But task_05 still settled `completed` having
changed nothing, because all four of its Verification commands pass most easily
when no work happened — the shape `SC-VERIFY-WORK-INDEPENDENT` now refuses at
authoring time.

This Task ships the route and then the documentation, in that order, and its
Verification is written so that doing nothing fails it.

## Requirements

1. MUST add `roundfix gc compact` with a preview-first contract: the bare form
   previews bytes before, reclaimable, and projected after; `--apply` performs
   the guarded compaction and reports the same three numbers as measured.
2. MUST route every guard `internal/store` already implements: an Active Run,
   another writer, and insufficient temporary capacity each refuse by name,
   before any mutation.
3. MUST list the subcommand in `roundfix gc --help`.
4. MUST leave per-repository GC and `gc sanitize` behaving exactly as they do
   today.
5. MUST then document all three operator surfaces in the canonical Roundfix
   Skill and its mirror: `gc compact`, `gc sanitize`, and `storage report`,
   including compaction's three refusals and that it is explicit rather than an
   automatic side effect of a retention sweep.
6. MUST regenerate the mirror with `make skills-sync` and run the ADR-0081
   chain, including the two characterization corpora that
   `make baseline-digests` does not reach.
7. MUST change only `internal/cli/**`, `.agents/skills/roundfix/**`,
   `skills/roundfix/**`, this Task file, and the ADR-0081 digest fallout.

## Subtasks

- [ ] Add the `gc compact` route over the existing store API.
- [ ] Document all three surfaces in the Skill and regenerate the mirror.

## Acceptance Criteria

- [ ] `roundfix gc compact` previews three numbers and exits 0.
- [ ] `roundfix gc compact --apply` compacts and reports measured numbers.
- [ ] Each of the three refusals is reachable from the command and names its
      cause.
- [ ] `roundfix gc --help` lists the subcommand.
- [ ] Per-repository GC and `gc sanitize` are unchanged.
- [ ] The Skill documents `gc compact`, `gc sanitize`, and `storage report`.
- [ ] The mirror is byte-identical to the canonical Skill.

## Context

- interface: `internal/cli/gc.go`
- interface: `internal/store/journal.go`
- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

Every command below fails if this Task does nothing. That is deliberate:
task_05 settled `completed` having changed nothing because all four of its
checks passed most easily on an untouched repository.

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix gc --help | grep -q "gc compact"`
  — expected: exit 0; the shipped help lists the subcommand.
- `go run -buildvcs=false ./cmd/roundfix gc compact > /dev/null` — expected:
  exit 0; the preview runs rather than rejecting the argument.
- `grep -q "gc compact" .agents/skills/roundfix/SKILL.md && grep -q "gc sanitize" .agents/skills/roundfix/SKILL.md && grep -q "storage report" .agents/skills/roundfix/SKILL.md`
  — expected: exit 0; the Skill names all three surfaces.
- `output="$(go test ./internal/cli -count=1 -run 'GCCompact|Compact' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the command tests ran and passed.
- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `make verify` — expected: exit 0.

## References

- `qa/qa-report-2026-08-06.md` → F-001, F-002.
- `_prd.md` → Core Feature 1; User Stories 1 and 2.
- `_techspec.md` → API Contracts.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.
