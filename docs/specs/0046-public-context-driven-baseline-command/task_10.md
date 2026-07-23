---
task: task_10
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: high
---

# Task 10: Classify root instructions through sealed ACP proposals

## Overview

Add optional semantic classification without granting the Agent access to the
checkout or mutation authority. Preferred and fallback selections analyze the
same sealed snapshot, while invalid or unavailable analysis returns the
deterministic workflow to manual classification.

## Requirements

1. MUST eagerly prove Codex `gpt-5.6-sol`/`xhigh` and Codex
   `gpt-5.5`/`xhigh` through Exact Agent Selection Proof without reading
   configurable Agent Work Categories.
2. MUST run each attempt in a fresh session and private empty directory with
   terminal and tools denied, one turn, bounded input, output, and deadline,
   and no checkout or Agent-log path.
3. MUST send byte-identical canonical snapshot input to preferred and fallback
   attempts and reject tools, extra prose, unknown or missing IDs, duplicate
   dispositions, digest mismatch, unsupported destinations, oversized output,
   timeout, or invalid JSON.
4. MUST discard raw and invalid ACP output, persist no Run or Run Event, and
   admit only a complete normalized proposal to consolidated human review.
5. MUST use manual classification when both selections fail and must not start
   fallback after parent cancellation or unproven session cleanup.

## Subtasks

- [x] Define the strict classification snapshot and proposal schemas.
- [x] Build the non-Run sealed ACP execution adapter.
- [x] Implement exact preferred/fallback supervision and strict cleanup.
- [x] Integrate validated proposals with root-instruction planning.
- [x] Add adapter, fallback, cancellation, privacy, and no-mutation tests.

## Acceptance Criteria

- [x] A valid preferred proposal prevents a fallback attempt.
- [x] Invalid preferred output is discarded and fallback receives identical snapshot bytes.
- [x] Tool use or checkout access is impossible under the emitted ACPX arguments.
- [x] Both unavailable or invalid selections return complete manual-classification destinations.
- [x] Parent cancellation closes the active session and never activates fallback.
- [x] No outcome creates repository changes, Run rows, Run Events, or Agent logs.
- [x] Equivalent accepted classifications produce the same Plan Digest regardless of proposal source.

## Context

- instruction: `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md`
- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/acp_stream.go`

## Verification

- `rtk go test -count=1 ./internal/baselineacp ./internal/baseline -run 'TestSealedClassification|TestPreferredFallback|TestProposalValidation|TestManualClassificationFallback|TestAnalyzerCancellation|TestAnalyzerNoPersistence'` — expected: proof, sandbox, validation, fallback, cleanup, privacy, and deterministic-plan cases pass.
- `rtk go test -count=1 ./internal/baselineacp -run TestACPXReadOnlyArguments` — expected: arguments deny tools and terminal, bound the attempt, omit full access, and expose no checkout path.
- `rtk make verify` — expected: the full repository gate passes without requiring a live model.

## References

- `_prd.md` → User Story 5; Core Features 8–9, 17, 19; Non-Goals / Out of Scope.
- `_techspec.md` → Interfaces: ProposalAnalyzer; API Contracts: ACP classification; Integration Points: ACPX/Codex ACP; Build Order 6.
- ADR-0039 → disposable exact-selection proof.
- ADR-0069 → supervised read-only semantic analysis.

## Result

Implemented strict
`roundfix/baseline-classification-snapshot/v1` and
`roundfix/baseline-classification-proposal/v1` contracts. The canonical
snapshot binds the complete Source Baseline evidence, allowed classifications,
two supported root dispositions, response fields, exact-source-byte rules,
and a domain-separated Snapshot Digest. Proposal parsing rejects duplicate or
unknown JSON fields, extra prose, stale digests, missing, unknown, or duplicate
entry IDs, unsupported destinations, changed proposed bytes, incomplete
reasons, invalid JSON, and output above 512 KiB. Snapshots reject more than
256 entries or 2 MiB of canonical input.

Added a non-Run ACPX path that accepts no Run, Run Event sink, checkout, or
Agent-log path. Each attempt uses a fresh random Agent Session and a fresh
empty mode-0700 directory. Session creation and prompting emit `--deny-all`,
`--non-interactive-permissions fail`, an empty `--allowed-tools`,
`--no-terminal`, `--max-turns 1`, zero prompt retries, and two-minute timeout
and TTL bounds. ACP tool events and unsupported updates fail closed. Raw ACP
JSONL, thought, and invalid model output remain ephemeral; only the bounded
Agent message reaches the strict proposal parser.

The supervisor directly and eagerly proves Codex `gpt-5.6-sol`/`xhigh` and
Codex `gpt-5.5`/`xhigh` without consulting Agent Work Categories. It sends the
same canonical snapshot bytes to each eligible attempt, accepts the preferred
proposal without activating fallback, retries invalid, timed-out, oversized,
or tool-using preferred output with the proven fallback, and returns a
complete deterministic manual proposal when analysis is unavailable or
invalid. Parent cancellation and any unproven Session or private-directory
cleanup stop fallback. Normalized proposals convert into the existing strict
Decision Document used by root-instruction planning; proposal origin and raw
ACP data never enter the Plan Digest.

Acceptance evidence:

- `TestPreferredFallbackValidPreferredPreventsFallback` proved both fixed
  selections before one preferred attempt and observed no fallback attempt.
- `TestPreferredFallbackInvalidPreferredUsesIdenticalSnapshotBytes` discarded
  invalid preferred output, activated fallback, and compared both complete
  prompt byte slices with each other and with `AnalysisSnapshot.CanonicalBytes`.
- `TestACPXReadOnlyArguments`,
  `TestSealedACPXPromptRejectsToolUseAndClosesSession`, and
  `TestAnalyzerNoPersistenceOrRepositoryMutation` proved the denied-tools,
  no-terminal ACPX arguments, tool-event rejection, private empty directories,
  absent checkout path, unchanged checkout sentinel, and no Run/Event/log
  persistence surface.
- `TestManualClassificationFallbackWhenSelectionsUnavailableOrInvalid` and
  `TestManualClassificationFallbackReturnsCompleteDestinations` proved every
  Source Baseline Entry receives either a complete Repository-Specific
  Normative Rules destination or a reasoned rejection after both semantic
  paths fail.
- `TestAnalyzerCancellationClosesActiveSessionWithoutFallback` and
  `TestSealedACPXPromptCancellationCancelsAndClosesSession` proved parent
  cancellation cancels and closes the active session before returning and
  never activates fallback.
- `TestPreferredFallbackDoesNotStartAfterProofCleanupFailure` and
  `TestPreferredFallbackRejectsToolUseAndStopsAfterCleanupFailure` proved no
  later proof or attempt starts after cleanup becomes unproven.
- `TestSealedClassificationEquivalentProposalsProduceSamePlanDigest` converted
  equivalent manual and parsed semantic proposals through the strict Decision
  Document boundary and produced the same portable Plan Digest.

Verification:

- `rtk go test -count=1 ./internal/baselineacp ./internal/baseline -run
  'TestSealedClassification|TestPreferredFallback|TestProposalValidation|TestManualClassificationFallback|TestAnalyzerCancellation|TestAnalyzerNoPersistence'`
  passed 24 tests.
- `rtk go test -count=1 ./internal/baselineacp -run
  TestACPXReadOnlyArguments` passed.
- `rtk go test -count=1 ./internal/agent ./internal/baselineacp
  ./internal/baseline` passed 408 tests.
- `rtk go test -race -count=1 ./internal/agent ./internal/baselineacp
  ./internal/baseline` passed.
- `rtk make verify` passed: 1,939 Go tests across 22 packages, both 256-test
  setup-context-driven suites, asset loading, Roundfix skill validation, and
  the final binary build.
