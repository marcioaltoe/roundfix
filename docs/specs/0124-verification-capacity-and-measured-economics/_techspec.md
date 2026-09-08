---
spec: 0124-verification-capacity-and-measured-economics
prd: _prd.md
created: 2026-09-08
---

# Measured verification capacity and economics — Technical candidate

## Executive Summary

Measure the repository gate in the contexts that actually execute it before changing concurrency, cache policy or deadlines. Reuse existing test-budget and compiled fixture infrastructure. The design accepts investigation outcomes with explicit unknowns rather than manufacturing a performance improvement or assigning a flaky failure to load without its evidence.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The cross-Spec order and remaining
decisions are in [the portfolio plan](../../workflow/2026-09-08-pending-work-plan.md).

## Project Constraints

- Identifier strategy: not applicable — retain existing Run, Task, and diagnostic identities; measurements do not introduce a new persisted identity scheme. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — verification scheduling and local/CI test execution introduce no authentication or HTTP policy. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0056 separates per-Run Task and Verification Capacity and expressly does not coordinate the entire machine. ADR-0117 requires each defect to be checked by the stage that can produce it. Preserve these boundaries unless a measured, explicitly approved design revises them. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0119 is an archived historical decision cited by an adopted measurement. Preserve that historical evidence; it is not a new operative authorization or a reason to weaken current refusal behavior.
  ADR-0147 supersedes the historical refusal decision: preserve an adapter-origin refusal and use the catalog as the net, not as authority to claim a better message or successful selection.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `Makefile`, `.github/workflows/ci-verify.yml`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Measurement record | `docs/references and existing test-budget/fixture evidence seams` | Describe toolchain, machine/package/test parallelism, cache mode, process ownership and command outcomes. |
| Failure capture | `existing Force Stop/worktree fixtures and event diagnostics` | Retain full failure output and ownership evidence before retrying. |
| Capacity policy | `Makefile Verification composition and existing task/verification capacity boundaries` | Apply only the measured, approved resource policy without changing Task meanings. |
| Cache and analyzer gate | `Makefile and CI Verification workflow` | Preserve complete/incremental semantics and add vet after0123 fixes the actual copied-lock defect. |
| Historical proposal disposition | `adopted source records and public measurement report` | Make evidence-backed adopt/defer decisions for optional detectors and test-mass proposals. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Comparable experiments

Author an experiment matrix that distinguishes warm incremental verification,
forced fresh test execution, genuinely cold compilation/cache and loaded
Daemon/Agent execution. Record command, toolchain/formatter, operating system,
package parallelism, test-binary parallelism, configured capacities, elapsed
time and exit. `-parallel` is not a machine-wide concurrency limit. Use isolated
resources and a bounded measurement window; no infinite rerun-until-green.
Preserve the unfiltered exit status and full failing output before summarizing.

The Go1.27.1/Go1.26.7 formatter difference observed during PR176 is an example
of execution-context evidence: this project declares Go1.26, and fresh gates
must use its toolchain rather than silently changing CI or format policy.
Do not claim the global developer environment has changed merely because one
command selected the correct toolchain.

### Failures, controls and decisions

The Force Stop and worktree-load observations retain unknown causes until an
actual failing run names them. Keep independent negative controls and record
when the failure does not recur; non-observation is not a fix. Preserve the
already implemented contextual diagnostics, serialized store writes, compiled
fixtures and six-Run tests. Do not revive the historical HOME-isolation idea.

Choose bounded capacity changes from measured bottlenecks. Distinguish Task
Capacity, Verification Capacity, package parallelism and whole-Spec delivery;
one cannot silently stand in for another. Any changed Makefile/CI/deadline
policy uses the exact grant and before/after evidence. Keep the existing CI
complete gate and its budget rather than raising it to accommodate a regression.

A cold-cache convenience command must name the cache the gate actually uses and
must not become the daily default. Add `go vet` only after0123 removes the copied
state ownership defect; an injected negative control proves the analyzer is in
the gate. For fixture deadline guards, external-record probes, fragile literals
and test-mass infrastructure, publish explicit adopt/defer dispositions with
reasons. A historical suggestion is not automatic approval to add a detector.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Store matched measurement context, samples, phase/resource observations and original failures. Retain the distinction between controlled measurements and inconclusive load-dependent observations.

### API Contracts

Keep current Verification entry points and output meaning until measurements justify an approved change. Reports identify their environment, original diagnostic and what remains unproved; no new success state conceals omitted work.

## Coverage Map

- PRD Goal 1 → Measurement record, Capacity policy.
- PRD Goal 2 → Failure capture.
- PRD Goal 3 → Cache and analyzer gate.
- PRD Goal 4 → Cache and analyzer gate.
- Core Feature 1 → Measurement record.
- Core Feature 2 → Failure capture.
- Core Feature 3 → Capacity policy.
- Core Feature 4 → Cache and analyzer gate.
- Core Feature 5 → Cache and analyzer gate.
- Core Feature 6 → Historical proposal disposition.

The Testing Approach below describes the observations that must settle these
contracts. Task IDs and actual evidence are deliberately not invented during
proposal authoring; approved Task decomposition must assign every contract and
success metric before execution.

## Integration Points

Local repository evidence and the adopted sources define the concrete seams.
The [owned source index](references/_index.md) records each primary source.
The [portfolio plan](../../workflow/2026-09-08-pending-work-plan.md) records
secondary consumers and prerequisite Specs. External research was read through
Exa and compared with local Secondbrain history; the PRD and portfolio plan
retain links and describe its effect. Published interfaces support feasibility,
not a claim that the proposed runtime or behavior already exists.

## Testing Approach

Use focused tests at the named package seams for deterministic rules, real
Git/store/process boundaries for integration behavior, and the authored public
QA Task for user-visible acceptance. Do not infer a terminal pass from source
inspection or a focused fixture. Required observations:

1. Identical sample definitions across incremental/fresh/cold contexts, including actual cache path and toolchain evidence.
2. Injected command failure remains nonzero with complete output and no summary pipeline hiding it.
3. Failure capture precedes any bounded repeat; missing reproduction retains unknown cause and original artifacts.
4. An intentional copied-lock negative control fails the approved analyzer gate after the real runtime repair passes.
5. Before/after measurements preserve existing functional and negative tests and report unchanged or worse outcomes honestly.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Bounded experiment and evidence contract (depends on: none).
2. Failure capture and matched baseline measurements (depends on: 1).
3. Measured capacity/cache decision and authorized implementation (depends on: 1, 2).
4. Analyzer gate after0123 and negative control (depends on: 1, 2).
5. Proposal dispositions, public report and terminal QA (depends on: 2, 3, 4).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Measurements are sensitive to contention and cache state. The experiment budget and window remain pending; record observations instead of inventing thresholds. The existing240-second CI gate is not changed by this plan.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge is confirmed only with independent review and required checks approved for the current candidate; releases, tags and paid consumption are not implied.
- Reviewer selection follows the [confirmed portfolio policy](../../workflow/2026-09-08-pending-work-plan.md); this Spec does not introduce a separate override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
