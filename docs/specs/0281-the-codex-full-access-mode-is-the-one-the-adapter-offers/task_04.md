---
task: task_04
spec: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
status: pending
type: test
complexity: low
---

# Task 04: The proof test tells the advertised mode from the first candidate

## Overview

Corrective Task 1 of 2. It answers finding F1 of the QA Report of
2026-10-19 (Run `run_20261019T205953Z_8c49586070d0475f`, row Q3,
task_03 Requirement 2). The QA gate made `chooseFullAccessMode` return the
first candidate whatever the session advertises, and
`TestProofSelectsTheAdvertisedCodexFullAccessMode` still passed: its only
session advertises `agent-full-access`, which is also the first Codex
candidate, so a selection that ignores the advertised modes produces the
same `set-mode`. That is the defect class this Spec fixes: a fixed mode name
sent without asking the session. The production code follows the TechSpec;
only the proof test is too weak to prove it.

## Requirements

1. MUST turn `TestProofSelectsTheAdvertisedCodexFullAccessMode` in
   `internal/agent/full_access_mode_test.go` into a table of two subtests
   that each run `ProveExactSelection` through the fake acpx harness, as the
   test does today:
   - `current_adapter`: the session advertises `read-only`,
     `workspace-write`, `agent` and `agent-full-access`; the proof sends one
     `set-mode agent-full-access` and reports the effective policy
     `full-access`. These are today's assertions, kept.
   - `legacy_adapter`: the session advertises `read-only` and `full-access`
     and not `agent-full-access`; the proof sends one `set-mode full-access`
     and reports the effective policy `full-access`.
2. MUST keep the test's name, so the QA gate's named negative control
   (task_03 Requirement 2) still names it, and keep
   `assertLastInvocationClosesDisposable` in both subtests.
3. MUST prove the negative control before handing back: in a temporary
   `git archive` copy of the tree outside the checkout, make
   `chooseFullAccessMode` return `candidates[0]` whenever candidates exist,
   run `go test -count=1 -run '^TestProofSelectsTheAdvertisedCodexFullAccessMode$' ./internal/agent`
   there, and observe it fail on `legacy_adapter`. Write nothing under the
   checkout for this proof, and leave no copy behind.
4. MUST NOT change production code, any other test, any skill, guide,
   record or planning file. The finding is a test sensitivity gap, not a
   production defect.

## Subtasks

- [ ] Split the proof test into the `current_adapter` and `legacy_adapter` subtests.
- [ ] Prove in a temporary copy that the first-candidate mutation fails `legacy_adapter`.

## Acceptance Criteria

- [ ] Both subtests pass on the intended code.
- [ ] A selection that returns the first candidate whatever the session advertises fails `legacy_adapter`.

## Consumers

None. The Task changes one test and no contract another package reads.

## Context

- instruction: `docs/adr/0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md`
- interface: `internal/agent/full_access_mode_test.go`
- instruction: `internal/agent/agent.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestProofSelectsTheAdvertisedCodexFullAccessMode$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestProofSelectsTheAdvertisedCodexFullAccessMode/current_adapter TestProofSelectsTheAdvertisedCodexFullAccessMode/legacy_adapter; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/agent || exit 1; GOOS=windows go vet ./internal/agent` — expected: exit 0. Before this Task the proof test has no subtests, so the command fails at the first missing pass; afterwards both subtests and the package pass.

## References

- `_techspec.md` → Invariant 1
- `task_03.md` → Requirement 2
