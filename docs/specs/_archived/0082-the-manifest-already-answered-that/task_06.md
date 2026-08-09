---
task: task_06
spec: 0082-the-manifest-already-answered-that
status: completed
type: backend
complexity: medium
---

# Task 06: Ask only what the manifest does not already answer

## Overview

The interactive command announces `update` mode and then asks all twelve
decisions anyway, each offering the stored value as a default, before blocking
on a spawned ACP runtime that re-segments the root instruction corpus. This task
short-circuits that path: when the manifest resolves, the interactive workflow
skips the settled prompts and the classification step and goes straight to the
plan, prompting only for decisions the manifest does not carry. First adoption
is untouched.

## Requirements

1. MUST skip the preservation-mode prompt, the profile prompt, and every decision
   prompt whose answer the resolved manifest carries.
2. MUST prompt for exactly those decisions the current catalog requires and the
   manifest does not carry, and for nothing else.
3. MUST NOT invoke the semantic analyzer, and therefore MUST NOT spawn an ACP
   runtime, when the manifest resolves.
4. MUST keep the plan confirmation prompt and its revision flow, so mutation
   still happens only against a reviewed Plan Digest.
5. MUST leave first adoption — a repository with no manifest — behaviorally
   identical, including its preservation prompt and supervised classification.
6. MUST fall back to the current full interactive path when the manifest is
   present but unreadable or incompatible, rather than refusing.
7. MUST keep the profile-change route reachable, so a maintainer can still move
   a repository to a different Baseline Profile interactively.

## Subtasks

- [x] Short-circuit the settled prompts when the manifest resolves.
- [x] Prompt only for newly required decisions.
- [x] Skip classification on the resolved path.
- [x] Keep the profile-change route reachable from the short-circuited path.
- [x] Prove first adoption and the incompatible-manifest fallback are unchanged.

## Acceptance Criteria

- [x] On a repository with a complete current manifest, the interactive workflow
      reaches the plan confirmation having issued zero decision prompts.
- [x] On a repository whose manifest lacks a catalog-required decision, exactly
      that decision is prompted.
- [x] No semantic analyzer call occurs on the resolved path, proven by a test
      whose injected analyzer fails the test if it is called.
- [x] On a repository with no manifest, the prompt sequence matches the task_01
      characterization corpus exactly.
- [x] On a repository with an unreadable manifest, the full interactive path runs
      rather than the command refusing.
- [x] On a repository whose stored profile digest no longer matches the catalog
      but whose profile resolves and whose decisions validate, the workflow
      announces update rather than adoption and issues zero decision prompts.
- [x] A maintainer can still choose to change the Baseline Profile.

## Context

- interface: `internal/cli/baseline_human.go`
- interface: `internal/baselineacp/analyzer.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/cli/ -run 'BaselineHuman' -v > /tmp/task_06-1.log 2>&1 && grep -q '^--- PASS: .*BaselineHuman' /tmp/task_06-1.log` — expected: exits 0.
- `go test ./internal/cli/ -run 'BaselineHuman' -v 2>&1 | grep -q -i 'analyzer'` — expected: exits 0, proving the never-called-analyzer case ran.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0, with the task_01 corpus proving first adoption is unchanged.

## References

- `_techspec.md` → Build Order 8; System Architecture.
- `_prd.md` → Core Features 7 and 8; User Story 2; Goal 1; Non-Goals: first adoption, profile changes.
- ADR-0068, ADR-0069, ADR-0099.

## Result

### Implementation

- The human workflow now resolves the Setup Manifest through
  `ResolveManifestInput`. Resolved and incomplete manifests announce `update`;
  unreadable or incompatible manifests retain the existing full-interview
  fallback instead of returning a manifest-read error.
- A manifest-backed Plan starts from the stored profile and decisions, prompts
  only for decision IDs still required by `ResolveDecisionInput`, uses
  `managed-refresh`, and never enters preservation selection, profile selection,
  segmentation, or classification.
- The final Plan Digest confirmation and reject-and-revise loop remain shared by
  adoption and update. Choosing Baseline Profile from that loop re-enters the
  full preservation/profile interview before recalculation, preserving the
  interactive profile-change route.

### Focused checks

- The pre-change focused run of the six Task 06 cases failed on the intended
  signals: complete manifests still requested preservation/profile input,
  profile-digest drift and a missing decision fell back to adoption, and a
  non-regular manifest returned a read error before the interview.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-0082-task-06-split-gocache go test ./internal/cli -run '^(TestBaselineHumanResolvedManifestSkipsPromptsAndAnalyzer|TestBaselineHumanProfileChangeRemainsReachable|TestHumanBaselineProfileDigestDriftRemainsUpdate|TestHumanBaselinePromptsOnlyForManifestMissingDecision|TestHumanBaselineFirstAdoptionPromptSequenceCharacterization|TestHumanBaselineUnreadableManifestFallsBackToFullInterview)$' -count=1`
  exited 0: `ok roundfix/internal/cli`.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-0082-task-06-regression-gocache go test ./internal/cli -run '^(TestHumanBaselineAdoption|TestConsolidatedReview|TestHumanBaselineInvokesSemanticSegmentationAndClassification|TestRejectedPlanRevision|TestRepeatedPlanRevisionDeterminism|TestBaselineNoTTY)$' -count=1`
  exited 0: `ok roundfix/internal/cli`.
- The commands under `## Verification` were not run; the Daemon owns them.

### Acceptance evidence

1. `TestBaselineHumanResolvedManifestSkipsPromptsAndAnalyzer` observes only the
   final Plan Digest confirmation on a complete current manifest and proves a
   declined Plan writes no repository bytes.
2. `TestHumanBaselinePromptsOnlyForManifestMissingDecision` removes only
   `verification.gate` and observes exactly its prompt followed by final Plan
   confirmation.
3. The injected `forbiddenBaselineSemanticAnalyzer` fails immediately from
   either `Segment` or `Classify`; the resolved and incomplete update cases pass
   with that analyzer installed.
4. `TestHumanBaselineFirstAdoptionPromptSequenceCharacterization` asserts the
   twelve pre-Plan prompt labels for a repository with no manifest, while the
   existing adoption, consolidated-review, and semantic-classification cases
   remain green.
5. `TestHumanBaselineUnreadableManifestFallsBackToFullInterview` presents a
   non-regular manifest path and reaches the preservation prompt with an
   incompatible-adoption announcement instead of returning the read error.
6. `TestHumanBaselineProfileDigestDriftRemainsUpdate` changes only the stored
   profile digest, observes `update`, builds the managed-refresh Plan with zero
   prompts, and keeps the analyzer forbidden.
7. `TestBaselineHumanProfileChangeRemainsReachable` rejects the initial update
   Plan, selects the Baseline Profile revision area, changes to `rust-cli`,
   receives a newly computed complete Plan, and declines it without mutation.
