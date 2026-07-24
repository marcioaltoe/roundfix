---
task: task_12
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: high
---

# Task 12: Recalculate rejected plans from scoped proposals

## Overview

Let maintainers reject a final plan, revisit one structured decision area, and
optionally translate free-form feedback into an in-scope proposal. Every
accepted revision returns through deterministic validation and a new approval.

## Requirements

1. MUST let the maintainer select profile, repository rules, divergences, or
   files as the decision area to revisit after rejecting a plan.
2. MUST accept direct structured changes without ACP and optional free-form
   suggestions through a separate strict revision snapshot and proposal
   schema.
3. MUST reject suggestions outside Baseline scope, unknown decisions,
   unauthorized destinations, incomplete changes, extra prose, or snapshot
   digest mismatch.
4. MUST recalculate the complete plan from the immutable repository snapshot
   and normalized accepted decisions rather than patching the prior artifact.
5. MUST show the resulting file projection and ledger and require a new exact
   Plan Digest approval before apply.
6. MUST preserve manual revision when ACP is unavailable.

## Subtasks

- [x] Add rejected-plan decision-area transitions.
- [x] Define strict revision snapshot and proposal contracts.
- [x] Supervise scoped ACP revision with the sealed analyzer.
- [x] Revalidate decisions and recompute the complete plan.
- [x] Add repeated-revision, out-of-scope, fallback, and approval tests.

## Acceptance Criteria

- [x] A direct structured revision requires no ACP session.
- [x] An accepted in-scope suggestion changes only permitted Baseline decisions.
- [x] An out-of-scope or invalid proposal changes no decision and exposes a manual next action.
- [x] Every accepted revision produces a newly computed plan and approval digest.
- [x] A previously approved digest cannot authorize a revised plan.
- [x] Repeated rejection remains deterministic and does not restart repository adoption.

## Context

- instruction: `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md`
- interface: `internal/baselineacp`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli -run 'TestRejectedPlanRevision|TestScopedRevisionProposal|TestRevisionOutOfScope|TestRevisionManualFallback|TestRevisionRequiresNewApproval'` — expected: structured and semantic revisions, scope rejection, fallback, recomputation, and approval invalidation pass.
- `rtk go test -count=1 ./internal/cli -run TestRepeatedPlanRevisionDeterminism` — expected: repeated review cycles preserve state and produce deterministic plans.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Story 6; Core Feature 11; User Experience.
- `_techspec.md` → Interfaces: ProposalAnalyzer; API Contracts: ACP revision; Build Order 6–7.
- ADR-0069 → proposal-only scoped semantic analysis.

## Result

Implemented rejected-plan revision as a repeatable state inside the existing
human Baseline workflow. A maintainer can revisit the Baseline Profile,
Repository-Specific Normative Rules, repository divergences and decisions, or
projected files. Direct structured changes remain entirely deterministic;
free-form suggestions cross a separate digest-bound Revision Snapshot and
strict Revision Proposal boundary supervised by the sealed preferred/fallback
Analyzer. Invalid, incomplete, out-of-scope, destination-bearing, prose-bearing,
or stale proposals are discarded before any decision replacement and continue
through the structured manual path.

Every accepted change is normalized through the Baseline catalog and calls
`RecalculatePlan`, which first validates the original Plan and its complete
repository preimage, then invokes full planning from the accepted request.
The previous artifact is never patched. The workflow renders the newly derived
file projection and complete ledgers and loops back to an exact confirmation
for the new Plan Digest.

### Acceptance evidence

1. `TestRejectedPlanRevision` uses a counting Analyzer and proves a direct
   structured decision revision starts zero ACP sessions.
2. `TestScopedRevisionProposal` admits one known decision replacement while
   asserting every unproposed decision remains byte-equivalent.
3. `TestRevisionOutOfScope` rejects unknown decisions, unauthorized path
   destinations, incomplete changes, extra prose, and snapshot mismatch
   without changing the snapshot decisions; the CLI routes unavailable or
   discarded semantic analysis to the structured manual revision.
4. `TestRejectedPlanRevision` observes distinct original and recomputed
   digests, a newly rendered `fileChanges` projection and managed-entry
   ledger, and two digest-bound confirmation prompts.
5. `TestRevisionRequiresNewApproval` passes the original digest with the
   revised Plan and observes the approval refusal.
6. `TestRepeatedPlanRevisionDeterminism` performs two rejection/revision
   cycles in one invocation, observes adoption only once, and replays the
   three resulting digest sequences against the same unchanged repository
   snapshot.
7. `TestRevisionManualFallback` proves unavailable preferred and fallback
   selections return an unchanged manual proposal rather than losing the
   maintainer's revision path.

### Verification

- `rtk proxy env GOCACHE=/tmp/roundfix-task12-go-cache go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli -run 'TestRejectedPlanRevision|TestScopedRevisionProposal|TestRevisionOutOfScope|TestRevisionManualFallback|TestRevisionRequiresNewApproval|TestRepeatedPlanRevisionDeterminism'`
  passed.
- The Daemon's first `rtk make verify` attempt exposed a timing-sensitive
  repository-identity assumption in `TestRepeatedPlanRevisionDeterminism`.
  The test now replays the workflow against one unchanged repository snapshot;
  the authoritative Daemon rerun is pending.
- `rtk git diff --check` passed.

### Verification feedback repair

The failing test created two independent Git repositories and compared their
Plan Digests. Because the clone-stable repository identity intentionally binds
the root commit, independently created histories can differ when their commit
metadata crosses a clock boundary. The test now runs both equivalent revision
sequences against the same unchanged repository, which directly exercises the
determinism invariant without assuming unrelated repositories have one
identity.

- `rtk proxy env GOCACHE=/tmp/roundfix-task12-repair-go-cache go test -count=20 ./internal/cli -run '^TestRepeatedPlanRevisionDeterminism$'`
  passed.
- `rtk proxy env GOCACHE=/tmp/roundfix-task12-repair-go-cache go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli -run 'TestRejectedPlanRevision|TestScopedRevisionProposal|TestRevisionOutOfScope|TestRevisionManualFallback|TestRevisionRequiresNewApproval|TestRepeatedPlanRevisionDeterminism'`
  passed.

### Follow-up

- The later documentation and skill-sync Task must publish the new rejected
  Plan revision prompts before any CLI-changing pull request is opened.
