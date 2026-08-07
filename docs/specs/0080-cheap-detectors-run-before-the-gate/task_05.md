---
task: task_05
spec: 0080-cheap-detectors-run-before-the-gate
status: pending
type: backend
complexity: high
---

# Task 05: Carry a row forward only on unmoved evidence

## Overview

The one mechanism in this Spec that can make something green wrongly, so it is
designed to fail closed in every direction and tested that way: the refusal
cases outnumber the happy path six to one, because the risk is entirely on the
permissive side.

A row is carried because nothing it depends on moved — never because it passed
recently. Recency would inherit a verdict; this inherits an observation, and
says so in the report.

## Requirements

1. MUST implement `Carriable` with every condition of ADR-0097 holding at
   once: the earlier report established the row as `pass`; the row declared a
   non-empty typed input list and every entry is `repository_path`; the
   earlier report's head is an ancestor of the current head; no declared input
   appears in the changed-path set between those heads; and every cited
   evidence path still resolves with unchanged content.
2. MUST refuse to carry a row that failed, was blocked by any cause, was
   skipped, or declared no inputs.
3. MUST refuse a mixed input list: any `external_repository`, `live_service`,
   or `elapsed_time` entry makes the row re-observe, even alongside
   repository paths.
4. MUST record, on every carried row, the report and the head that established
   it, so a reader sees which evidence is fresh and which is inherited.
5. MUST NOT make any verdict more permissive than a fresh observation would;
   ADR-0080 keeps verdict semantics and a carried row cannot change a count's
   meaning.
6. MUST reuse the existing changed-path primitive rather than adding a second
   way to compute what moved.

## Subtasks

- [ ] Implement the resolver and its ancestry and changed-path checks.
- [ ] Author the refusal suite before the happy path.
- [ ] Emit the establishing citation into the carried row.

## Acceptance Criteria

- [ ] Each refusal condition has its own named test and each refuses.
- [ ] A row meeting every condition carries, citing its establishing report
      and head.
- [ ] A mixed input list never carries.
- [ ] No verdict computation changed.

## Rehearsal Cases

- Case: prior status is not `pass`; Observation: `Carriable` returns false and
  the row is re-observed.
- Case: row declares no inputs; Observation: refused, so declaration is
  opt-in.
- Case: a declared `repository_path` appears in the changed-path set;
  Observation: refused.
- Case: the establishing head is not an ancestor of the current head;
  Observation: refused.
- Case: a cited evidence file resolves but its content changed; Observation:
  refused.
- Case: inputs mix `repository_path` with `elapsed_time`; Observation:
  refused.
- Case: every condition holds; Observation: carried, with the establishing
  report path and head recorded on the row.

## Context

- interface: internal/speccheck/citations.go
- interface: internal/worktree/worktree.go
- interface: internal/daemon/task_engine.go

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test -count=1 ./internal/... -run 'Carriable|CarryForward' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the carry-forward tests are selected and pass.
- `output="$(go test -count=1 ./internal/... -run 'Carriable|CarryForward' -v 2>&1)"; printf '%s\n' "$output" | grep -c -- '--- PASS' | { read n; [ "$n" -ge 7 ]; }`
  — expected: exit 0; at least the seven rehearsal cases are present as
  passing subtests, so the refusal suite cannot be reduced to one happy path.
- `grep -rq 'PriorChangedFiles' internal/ && go test -count=1 ./internal/daemon/... ./internal/speccheck/...`
  — expected: exit 0; the existing changed-path primitive is the one in use and
  both packages stay green.

## References

- `_prd.md` → Core Feature 5; User Stories 2, 6; Success Metrics.
- `_techspec.md` → Implementation Design (Carry-forward); Testing Approach;
  Build Order 5.
- ADR-0097, ADR-0080.
