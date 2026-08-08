---
task: task_02
spec: 0088-a-third-runtime-that-can-run
status: pending
type: backend
complexity: high
---

# Task 02: Retain capability values by relevance instead of by position

## Overview

Make the capability projection accept a large advertised option and bound what it
keeps. Today a `select` option above 64 values is discarded whole, which leaves
no model state and no Exact Agent Selection Proof on any runtime that advertises
a catalog. After this Task, an OpenCode-sized payload projects into capabilities
whose model list contains the requested model, and a model that is genuinely
unadvertised still fails closed.

## Requirements

1. MUST add a retention input naming the Agent Selection being proven — its model
   and its reasoning effort — and MUST thread it from both capability read paths
   so the projection never infers it.
2. MUST keep every advertised value unchanged when an option advertises at or
   below the retained bound, so Codex and Claude projections are byte-identical
   to today.
3. MUST, above the retained bound, keep the option's current value, keep every
   value that binds to the requested model by exact advertised value or by
   canonical prefix before a trailing bracketed effort, and fill the remainder in
   advertised order up to the bound.
4. MUST record, per option, how many values the adapter advertised alongside how
   many were retained, so a diagnostic can state both.
5. MUST keep a hard ceiling above which an advertised option is still refused as
   `too_many_option_values`, and MUST set that ceiling far enough above the
   measured catalog that a legitimate catalog does not reach it.
6. MUST raise the capability response byte ceiling to accommodate the measured
   payload with headroom, and MUST state the measured size in a comment beside
   the constant.
7. MUST leave the fail-closed contract intact: a requested model absent from the
   advertised values produces `SelectionUnsupportedError`, never
   `capability_evidence_invalid`.
8. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [ ] Introduce the retention input and thread it through both read paths.
- [ ] Replace the size-as-malformation rule with the retention rule.
- [ ] Record advertised-versus-retained counts on the projected option.
- [ ] Move the three bounds and comment the measured payload size.
- [ ] Edit the break-half characterization test that pinned the old refusal, and
      declare the break in this Task's Result.
- [ ] Re-record the coverage record if a test name moved.

## Acceptance Criteria

- [ ] A `model` option advertising more values than the bound, whose values
      include `opencode-go/kimi-k3`, projects into capabilities whose model list
      contains `opencode-go/kimi-k3`.
- [ ] The same payload with the requested model removed produces
      `SelectionUnsupportedError` and no `capability_evidence_invalid`.
- [ ] A five-value Claude-shaped option, including a bracketed variant, projects
      exactly as it does today, value for value and in the same order.
- [ ] A requested canonical model whose advertised form carries a bracketed
      effort suffix is retained when the option is above the bound.
- [ ] The projected option reports the advertised count and the retained count,
      and they differ for an oversized option.
- [ ] An option above the hard ceiling still fails with
      `too_many_option_values`.
- [ ] The invariant half of the characterization corpus passes unmodified.

## Context

- interface: `internal/agent/selection_capabilities.go`
- interface: `internal/agent/selection_assignment.go`

## Bounded scope

This Task may create or modify only:

- `internal/agent/selection_capabilities.go`
- `internal/agent/selection_capabilities_test.go`
- `internal/agent/selection_capabilities_characterization_test.go`
- `internal/agent/acpx_runner_test.go`
- `docs/references/coverage-record.json`
- `docs/specs/0088-a-third-runtime-that-can-run/task_02.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/agent -count=1` — expected: exits 0.
- `go test ./internal/agent -run 'Retention' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `go test ./internal/agent -run 'CharacterizationInvariant' -count=1 -v` — expected: exits 0, proving the invariant half still holds.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0, proving no test was retired without re-recording.
- `grep -q 'maxAdvertisedCapabilityValues' internal/agent/selection_capabilities.go` — expected: exits 0, proving the hard ceiling exists.
- `grep -q 'AdvertisedCount' internal/agent/selection_capabilities.go` — expected: exits 0, proving the counts are recorded.

## References

- `_prd.md` → Goal 1; Core Features: relevance-bounded capability retention.
- `_techspec.md` → Implementation Design: Interfaces, Data Models; Build Order 2.
- `references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md`
  → the 417-value catalog and the 50,590-byte payload.
- ADR-0105.
