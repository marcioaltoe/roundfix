---
task: task_04
spec: 0079-one-door-for-fleet-knowledge
status: pending
type: backend
complexity: high
---

# Task 04: Refuse what the findings contract forbids

## Overview

The three repository-level consistency checks land in `internal/speccheck`
beside the loop-order rule that proved the shape: `SC-FINDING-LIFECYCLE` (an
active finding without its lifecycle frontmatter), `SC-ROLLUP-MEMBER` (a
rollup declaring a member that does not resolve), and `SC-ARCHIVE-LICENSE`
(an archived finding without a resolvable `absorbed_by`). They arrive on a
repository the sweep (task_03) already made conformant, so the gate goes
green the day the rules start demanding it.

## Requirements

1. MUST implement the three checks as repository-level rules with the byte-
   exact codes `SC-FINDING-LIFECYCLE`, `SC-ROLLUP-MEMBER`, and
   `SC-ARCHIVE-LICENSE`, reporting file, line, and fix line in the same
   shape as every existing rule.
2. MUST keep every check citation-only per ADR-0093: a member is a written
   path, a license is a written pointer, a status is a written value — no
   check infers relevance, theme, or intent.
3. MUST keep every check presence-aware per ADR-0094: a repository without a
   findings directory, without rollups, or without an archive skips the
   corresponding check silently instead of failing.
4. MUST validate lifecycle values against exactly `pending`, `partial`,
   `deferred`, and `done`; an unknown value is an `SC-FINDING-LIFECYCLE`
   failure with the offending value in the message.
5. MUST resolve rollup members against active and archived findings, and
   licenses against active findings (rollups), active Spec folders, and
   archived Spec folders — the same resolution the sweep's relations use.
6. MUST exercise each check through a dedicated fixture carrier in the
   loop-order idiom, never through live or archived Specs, with a
   red-then-green pair per check authored against the real report surface.
7. MUST include a presence-aware case per check proving the silent skip, and
   MUST keep the existing fixture corpus green with `TestCheckCorpusBudget`
   still passing — over-reach is the standing failure mode of new rules.

## Subtasks

- [ ] Implement the three rules and their report wiring.
- [ ] Build the findings fixture carrier; author red/green pairs.
- [ ] Add presence-aware skip cases; hold the corpus and budget green.

## Acceptance Criteria

- [ ] Each check fails its red fixture with its own code and passes its
      green fixture.
- [ ] Each check skips silently when its artifact class is absent.
- [ ] The pre-existing corpus reports no new diagnostics and the budget
      holds.
- [ ] `spec check` over this repository is green under the new rules.

## Context

- interface: internal/speccheck/citations.go
- interface: internal/speccheck/citations_test.go
- interface: internal/speccheck/report.go

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test -count=1 ./internal/speccheck -run 'FindingLifecycle|RollupMember|ArchiveLicense' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the named check tests exist, are selected, and pass —
  an empty selection cannot pass this gate.
- `output="$(go test -count=1 ./internal/speccheck -run '^TestCheckCorpusBudget$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the corpus budget guard still selects and passes.
- `grep -q 'SC-FINDING-LIFECYCLE' internal/speccheck/*.go && grep -q 'SC-ROLLUP-MEMBER' internal/speccheck/*.go && grep -q 'SC-ARCHIVE-LICENSE' internal/speccheck/*.go`
  — expected: exit 0; the three byte-exact codes exist on the check surface.
- `go run -buildvcs=false ./cmd/roundfix spec check > /dev/null`
  — expected: exit 0; this repository, swept by task_03, is conformant under
  its own new rules.

## References

- `_prd.md` → Core Feature 10; Goal 4.
- `_techspec.md` → Implementation Design (the three checks); API Contracts;
  Testing Approach; Build Order 3.
- ADR-0093, ADR-0094.
- `references/2026-08-06-findings-accumulate-faster-than-they-become-specs.md`
  → why advice was not enough.
