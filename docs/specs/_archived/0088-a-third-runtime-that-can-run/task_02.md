---
task: task_02
spec: 0088-a-third-runtime-that-can-run
status: completed
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

- [x] Introduce the retention input and thread it through both read paths.
- [x] Replace the size-as-malformation rule with the retention rule.
- [x] Record advertised-versus-retained counts on the projected option.
- [x] Move the three bounds and comment the measured payload size.
- [x] Edit the break-half characterization test that pinned the old refusal, and
      declare the break in this Task's Result.
- [x] Re-record the coverage record if a test name moved.

## Acceptance Criteria

- [x] A `model` option advertising more values than the bound, whose values
      include `opencode-go/kimi-k3`, projects into capabilities whose model list
      contains `opencode-go/kimi-k3`.
- [x] The same payload with the requested model removed produces
      `SelectionUnsupportedError` and no `capability_evidence_invalid`.
- [x] A five-value Claude-shaped option, including a bracketed variant, projects
      exactly as it does today, value for value and in the same order.
- [x] A requested canonical model whose advertised form carries a bracketed
      effort suffix is retained when the option is above the bound.
- [x] The projected option reports the advertised count and the retained count,
      and they differ for an oversized option.
- [x] An option above the hard ceiling still fails with
      `too_many_option_values`.
- [x] The invariant half of the characterization corpus passes unmodified.

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

## Result

The capability projection now bounds what it retains. A 417-value catalog
projects; a model nobody advertises still fails closed.

**What changed.** `ParseSessionConfigOptions` and
`ParseSessionCapabilitySnapshot` take a `SelectionRetention` naming the
requested Agent Selection, and both read paths supply it from the RuntimeSpec
through `RetentionFor`, so the projection never infers it. At or below
`maxRetainedCapabilityValues` every advertised value survives in advertised
order, which is why the Codex and Claude fixtures project unchanged. Above it,
`retainAdvertisedValues` keeps the current value, every value binding to the
requested model — by exact advertised value or by canonical prefix before a
trailing bracketed effort, the same rule `modelsForCanonical` applies — and then
fills to the bound in advertised order. `SelectCapability.AdvertisedCount`
records what the adapter offered. `maxCapabilityValues` became
`maxRetainedCapabilityValues`; `maxAdvertisedCapabilityValues` (4096) is the new
absolute refusal ceiling, so `too_many_option_values` still fires on an
implausible payload rather than becoming dead vocabulary.
`maxCapabilityResponseBytes` moved from 64 KiB to 1 MiB, with the measured
50,590-byte payload named in the comment beside it.

**Declared breaks.**

1. `TestCharacterizationTodayRefusesAnOversizedAdvertisedOption` became
   `TestCharacterizationDeclaredBreakOversizedOptionRetainsInsteadOfRefusing`.
   The old assertion — the `too_many_option_values`, `missing_model_state`,
   `contradictory_response` triple — is preserved in that test's comment as the
   behavior this Task removed, and the measured provenance stays with it.
2. `TestSelectionCapabilitiesOfficialAndLegacyFixtures` gained
   `AdvertisedCount: 6` in its expected reasoning option, because
   `SelectCapability` gained a field and the test compares the whole struct.

**Commands and outcomes.**

- `go build -buildvcs=false ./...` — exit 0.
- `go test ./internal/agent -count=1` — exit 0.
- `go test ./internal/agent -run 'Retention' -count=1 -v` — exit 0; four tests.
- `go test ./internal/agent -run 'CharacterizationInvariant' -count=1 -v` — exit 0; four tests, the invariant half unchanged.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — exit 0.
- `grep -q 'maxAdvertisedCapabilityValues' internal/agent/selection_capabilities.go` — exit 0.
- `grep -q 'AdvertisedCount' internal/agent/selection_capabilities.go` — exit 0.
- `make verify` — exit 0 after `go clean -testcache`.

**Evidence per acceptance criterion.**

- Requested model retained past the bound:
  `TestRetentionKeepsRequestedModelPastTheBound` advertises the model **last**
  among 192 values, so a projection keeping an advertised prefix would drop it;
  the test then plans the assignment and asserts the adapter model.
- Unadvertised model still unsupported:
  `TestRetentionLeavesUnadvertisedModelUnsupported` asserts
  `SelectionUnsupportedError` and asserts the absence of
  `CapabilityEvidenceError`.
- Small advertised sets unchanged:
  `TestCharacterizationInvariantRetainsEveryValueAtOrBelowTheBound` compares
  four model values and six effort values in advertised order, and the whole
  pre-existing `internal/agent` suite passes.
- Bracketed variant of a canonical request retained:
  `TestRetentionKeepsBracketedVariantOfRequestedCanonicalModel` requests `opus`
  and asserts `opus[1m]` survives.
- Advertised and retained counts differ and are both reported:
  `TestRetentionRecordsAdvertisedAndRetainedCounts` asserts 192 advertised and
  64 retained.
- Absolute ceiling still fails closed:
  `TestCharacterizationInvariantOversizedOptionStillFailsClosedAboveTheCeiling`
  advertises 4097 values and asserts `too_many_option_values`.

**Follow-ups.** The coverage record was re-recorded in this commit, adding the
thirteen tests Tasks 01 and 02 introduced, so a future removal of any of them is
caught rather than logged.
