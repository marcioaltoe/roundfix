---
task: task_16
spec: 0119-spec-contained-authorization
status: pending
type: backend
complexity: medium
---

# Task 16: Stop the suite guard from reparsing the archive per package

## Overview

Widening grant discovery to archived Specs made the suite guard walk
`docs/history/specs` on every guarded test package: 105 archived Specs and
1763 markdown files, read and YAML-parsed once per test binary. Under the full
concurrent suite that cost pushes the Daemon's Task-cycle deadlines past their
budgets. The discovery is correct and stays; paying for it repeatedly is the
defect.

## Requirements

1. MUST stop reading and parsing markdown files that cannot carry an
   authorization record, so the walk's cost scales with records rather than
   with every file under the Spec roots.
2. MUST resolve the sanctioned regeneration set at most once per process and
   reuse it, because an archived record is immutable within a run and reparsing
   it yields the same answer at full cost.
3. MUST keep every discovery outcome the current implementation reaches: active
   Spec grants, archived Spec grants, preserved legacy records without
   frontmatter, multi-Spec consuming lists, and the refusal of proposed,
   null-dated, malformed or unrelated records.
4. MUST NOT narrow the discovered roots. Dropping the archive would trade this
   cost for the disagreement between readers that Task 09 removed.
5. MUST bring the Daemon package's wall clock back to the unchanged target's
   range under the full concurrent suite. Measured on 2026-09-10: the package
   runs 4.1s at the delivery target and 6.8s on the candidate in isolation, and
   the candidate misses its deadlines under `-parallel 16`.
6. MUST NOT raise a time budget, lengthen a deadline, reduce parallelism, or
   skip a case. The cost goes, not the signal.

## Subtasks

- [ ] Filter the walk to files that can carry a record before reading them.
- [ ] Resolve the regeneration set once per process and reuse it.
- [ ] Prove every current discovery outcome still holds.
- [ ] Measure the Daemon package against the unchanged target.

## Acceptance Criteria

- [ ] The guard reads a number of files proportional to candidate records, not
      to every markdown file under the Spec roots; a test asserts the read count
      does not grow with an added non-record file.
- [ ] Resolving twice in one process performs the underlying walk once.
- [ ] Active, archived, legacy frontmatter-free and multi-Spec records all still
      resolve, and proposed, null-dated, malformed and unrelated records still
      contribute nothing.
- [ ] The Daemon package passes as one concurrent run, including the seven
      Task-cycle cases that miss their deadlines today.
- [ ] No time budget, deadline, parallelism setting or skip was changed, proven
      by the absence of such edits in the diff.

## Context

- interface: `internal/suiteguardcontract/regeneration.go`
- interface: `internal/suiteguard/suiteguard.go`

## Verification

- `grep -q 'func TestSanctionedRegenerationReadsOnlyCandidateRecords' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestSanctionedRegenerationReadsOnlyCandidateRecords$'` — the walk's read count tracks records rather than every markdown file.
- `grep -q 'func TestSanctionedRegenerationResolvesOncePerProcess' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestSanctionedRegenerationResolvesOncePerProcess$'` — a second resolution reuses the first walk.
- `grep -q 'func TestSanctionedRegenerationReadsOnlyCandidateRecords' internal/suiteguardcontract/regeneration_test.go || exit 1; go test -count=1 ./internal/suiteguardcontract` — every preserved discovery outcome still holds once the cheap walk is in place. The guard reads the new case, so this preservation check cannot pass before the work.
- `grep -q 'func TestSanctionedRegenerationResolvesOncePerProcess' internal/suiteguardcontract/regeneration_test.go || exit 1; go test -count=1 ./internal/daemon` — the Daemon package passes as one concurrent run once the cost is gone.
- `grep -q 'func TestSanctionedRegenerationReadsOnlyCandidateRecords' internal/suiteguardcontract/regeneration_test.go || exit 1; widened="$(grep -n 'time.After(' internal/daemon/task_engine_test.go | grep -v '2 \* time.Second' || true)"; test -z "$widened" || { printf '%s\n' "$widened"; exit 1; }; skips="$(git diff --unified=0 main -- internal/daemon internal/suiteguard | grep -E '^\+' | grep -F 't.Skip(' || true)"; test -z "$skips" || { printf '%s\n' "$skips"; exit 1; }` — every Task-cycle deadline still reads two seconds and no case was skipped, so the cost went rather than the signal. The check reads the deadline value and a skip call, not words: an earlier version matched `parallel` and flagged `t.Parallel()`, which enables parallelism, and matched `t.Skip` inside `result.Skipped`.

## References

- `_prd.md` → Decisions: Regression locks; Promoted work and known limitations.
- `_techspec.md` → Implementation Design: Audit and compatibility.
- `qa/qa-report-2026-09-10-02.md` → F-005.
