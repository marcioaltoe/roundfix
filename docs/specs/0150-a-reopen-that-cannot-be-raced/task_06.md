---
task: task_06
spec: 0150-a-reopen-that-cannot-be-raced
status: pending
type: test
complexity: medium
---

# Task 06: Public CLI evidence that builds what it audits

## Overview

Pre-PR review found that this Spec's QA Report records `verdict: pass` while its
own frontmatter says the binary it exercised is stale:

```
auditing_binary: "0.14.1 (b1c8a33d, built 2026-09-19 17:30:06 -0300)"
auditor_staleness: "stale: commit ancestry: build commit predates audited tree"
```

That binary predates Task 05, so the rows meant to prove the public command's
behavior proved an older command's behavior. The focused unit tests did run
against the delivered tree, but the public-CLI evidence did not.

The durable repair is not to re-run an Agent with better luck. It is to make the
public-CLI evidence a test that builds the binary it audits, so the proof lives
in the suite and cannot drift from the tree again. The repository already does
this elsewhere: `internal/cli/baseline_plan_test.go` builds `./cmd/roundfix` and
runs it, and `buildRoundfixBinaryForMacro` exists for the same reason.

## Requirements

1. MUST build the command from the tree under test and exercise `reopen`
   through that binary, not through an in-process entry point.
2. MUST cover, through the built binary: a stale gate reopening, the refusal
   when the gate is not stale, and a write through a symlinked Task path landing
   on the target while the link survives.
3. MUST assert the prior QA Report is byte-identical after a successful reopen.
4. MUST NOT weaken or duplicate the unit coverage Tasks 01, 02 and 05 deliver;
   this is the public-surface layer above them.

## Subtasks

- [ ] Add the built-binary test following the existing build-and-run precedent.
- [ ] Cover the three journeys and the report byte-identity assertion.

## Acceptance Criteria

- [ ] The test builds `./cmd/roundfix` and invokes the resulting binary.
- [ ] A stale gate reopens through it, and the QA Report bytes are unchanged.
- [ ] A healthy gate is refused through it.
- [ ] A symlinked Task path has its target rewritten and remains a symlink.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/reopen.go`
- interface: `internal/cli/baseline_plan_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReopenThroughTheBuiltBinary" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReopenThroughTheBuiltBinary"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `grep -q "cmd/roundfix" internal/cli/reopen_test.go` — expected: exit 0; the reopen tests build the command under test. Before this Task they do not, so the command fails.

## References

- [_techspec.md](_techspec.md) — Testing Approach
