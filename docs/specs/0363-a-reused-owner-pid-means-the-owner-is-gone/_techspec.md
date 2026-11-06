---
spec: 0363-a-reused-owner-pid-means-the-owner-is-gone
prd: _prd.md
created: 2026-11-06
---

# A reused owner PID means the owner is gone — Technical Spec

## Executive Summary

The store gains one question, `OwnerGone`, answered by the existing
ownership proof: the owner is gone when its PID is absent or when the live
start identity at that PID was read, is comparable and differs from the
recorded one. Doctor's `Residue` skips a Run whose lineage read fails with
the existing `ErrOwnerProcessIdentityMismatch`, and the three readers that
asked only `ProcessAlive` ask `OwnerGone` instead. The trade-off this design
accepts is one more start-identity read for a live PID with a recorded
identity, in exchange for one rule that every liveness reader shares; the
signal paths keep their own stricter proof and are not touched.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0373 (new) records the rule.
  ADR-0313: "A proven owner identity mismatch releases a terminal Run's owner
  record without a signal", which `OwnerGone` reuses. ADR-0044: "Orphaned Run
  locks are reclaimed only on proven owner death", kept by counting a
  mismatch as that proof. ADR-0127: "Process residue is a readiness fact, not
  a Run record", so Doctor writes nothing. ADR-0279: "Status reports a Queue
  Owner proven gone", whose readers stay unchanged. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 4 rules answer of 2026-11-06 quoted in the authorization
  record; no Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0363-a-reused-owner-pid-means-the-owner-is-gone/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## System Architecture

- `internal/store/process.go`: `OwnerProcessControl.proveOwner` already
  returns `(absent, nil)` for an absent owner and an error wrapping
  `ErrOwnerProcessIdentityMismatch` for a proven mismatch; the new method and
  package function `OwnerGone` read it and never signal.
- `internal/cli/doctor.go`, `Residue`: a `ReadLineage` error that
  `errors.Is` `store.ErrOwnerProcessIdentityMismatch` makes the Run
  contribute nothing; every other error keeps today's text.
- `internal/cli/cli.go`, `reclaimOrphanedActiveRun` and
  `orphanedActiveRunReason`; `internal/store/store.go`,
  `Store.ReclaimOrphanedRun` (its own guard); `internal/cli/deliver_workflow.go`,
  `RunSlots`; `internal/cli/deliver.go`, `runDeliverStart` (the existing
  release before `CreateDeliveryQueueWithLimits`): each replaces
  `store.ProcessAlive(pid)` with `OwnerGone(ctx, pid, recordedIdentity)`.
- Unchanged by design (ADR-0373): `forceStopRun`, `runDeliverStop`,
  `TerminateAndWait`, `TerminateTreeAndWait`, `SignalProvenProcess`,
  `inspectReconcileProcesses`, `deliveryOwnerLine`, `runDeliverRetry`,
  `runDeliverResume`, `runs list`, `gc`, worktree disposal and
  `inspectStagingRegistration`.
- Docs: `docs/user-guide/commands/doctor.md`,
  `docs/user-guide/commands/deliver.md`, and the Roundfix skill's
  `references/setup.md` and `references/deliver.md`.

## Implementation Design

### Interfaces

```go
// internal/store/process.go
// OwnerGone reports whether the recorded owner no longer runs: its PID is
// absent, or alive with a proven different start identity (reused). It sends
// no signal; any answer short of that proof reports the owner present.
func (controller *OwnerProcessControl) OwnerGone(ctx context.Context, pid int, recordedIdentity string) (gone, reused bool)
func OwnerGone(ctx context.Context, pid int, recordedIdentity string) (gone, reused bool)

// internal/cli/cli.go
func orphanedActiveRunReason(pid int, reused bool) string

// internal/cli/deliver.go
func releaseGoneDeliveryQueueOwner(ctx context.Context, runStore *store.Store, gitRoot string) error
```

### Data Models

No changes. Runs and Delivery Queues already record `owner_pid` and the
owner's start identity.

### API Contracts

1. API Contract: `OwnerGone`. `pid <= 0` and the caller's own PID →
   `(false, false)`. Proof `(absent, nil)` → `(true, false)`. Proof error
   that `errors.Is` `ErrOwnerProcessIdentityMismatch` → `(true, true)`. Proof
   `(false, nil)`, including a live PID with an empty recorded identity (the
   legacy PID-only proof), and every other error (unreadable identity,
   foreign-platform token, process-table failure) → `(false, false)`.
2. API Contract: Doctor residue. A terminal Run whose `ReadLineage` returns
   an error matching `ErrOwnerProcessIdentityMismatch` adds no residue and no
   read error. With no other Run reporting, the check prints
   `residue: ok (no process residue found)`. Every other line keeps its
   bytes, including `could not read process table for Run <id>: recorded
   owner identity is unproven` and the wrapped read errors.
3. API Contract: the Active Run lock. `reclaimOrphanedActiveRun` reclaims
   when `OwnerGone` reports gone, with reason
   `owner process <pid> gone (pid reused); lock reclaimed` when reused and
   `owner process <pid> not running; lock reclaimed` otherwise;
   `Store.ReclaimOrphanedRun` returns without a write unless `OwnerGone(ctx,
   pid, run.OwnerIdentity)` reports gone, and its default reason follows the
   same two forms.
4. API Contract: Run slots and queue start. `RunSlots` excludes a holder
   whose recorded owner `OwnerGone` reports gone. The release that
   `runDeliverStart` performs before creating a queue moves into
   `releaseGoneDeliveryQueueOwner(ctx, runStore, gitRoot) error`, which
   releases an existing queue's owner when `OwnerGone(ctx, existing.OwnerPID,
   existing.OwnerIdentity)` reports gone and keeps today's error text when
   the release fails. No output text changes.

### Surface Transcripts

None. The only changed lines are an existing Doctor line that drops a
`partial` result and the reclaim reason inside the existing warning, which
`internal/cli` tests assert through fake process controllers and the test's
own parent process; no command transcript can reproduce a reused PID on
demand.

## Coverage Map

- Goal 1, Story 1, Core Feature 2, Success Metric 1 → API Contract 2, Build
  Order 2.
- Goal 2, Success Metric 2 → API Contract 2 (unchanged lines), Build Order 2.
- Core Feature 1, Success Metric 3 → API Contract 1, Build Order 2.
- Goal 3, Stories 2-3, Core Feature 3, Success Metric 4 → API Contracts 3-4,
  Build Order 3.
- Goal 4, Success Metric 5 → System Architecture (unchanged by design),
  Testing Approach, Build Order 3.
- Core Feature 4 → Build Order 1.

Deviations from the PRD: None.

## Absence Evidence

1. `OwnerGone` (the store question) — command: `grep -rn 'func OwnerGone\|ProcessAlive(' internal --include=*.go --exclude=*_test.go`; result: no `OwnerGone`; `ProcessAlive` is called on the PID alone by `Store.ReclaimOrphanedRun`, `reclaimOrphanedActiveRun`, `RunSlots`, `runDeliverStart` and the `deliver` status, retry and resume paths, the last three after an identity proof.
2. The Doctor mismatch branch — command: `grep -n 'ErrOwnerProcessIdentityMismatch' internal/cli/doctor.go`; result: no match; `Residue` wraps every `ReadLineage` error as `could not read process table`.

## ADR Alternatives

- ADR-0373: the design takes none of its rejected options: no residue line
  for a reused owner, no tree walk after a mismatch, no unreadable identity
  counted as gone, and no release from Doctor.
- ADR-0313: rejected releasing from Doctor and other read commands; Doctor
  still releases nothing, and `reconcile --apply` keeps the release.
- ADR-0044: rejected reclaiming a lock on anything short of proof of death;
  a proven mismatch is that proof, and an unreadable identity still blocks.
- ADR-0127: rejected inventing Run rows for residue; Doctor still writes no
  record.
- ADR-0279: rejected changing a queue from `status`; status stays read-only
  and its `Owner:` lines keep their bytes.

## Integration Points

The host process table only, through the existing platform readers
(`kern.proc` on macOS, `/proc` on Linux); no signal is sent.

## Testing Approach

- `internal/store/process_owner_gone_test.go` (new): `OwnerGone` over a
  controller whose `processAbsent` and `processStartIdentity` are fakes, as
  `process_identity_mismatch_test.go` builds them, for every case of API
  Contract 1.
- `internal/store/reclaim_reused_owner_unix_test.go` (new, `//go:build
  unix`): `ReclaimOrphanedRun` with the test's live parent PID and the
  identity `runtime.GOOS + ":0.0"` reclaims with the reuse reason; with the
  parent's real identity (`OwnerProcessIdentity`) it writes nothing. Skips
  naming the reason when the platform has no start identity.
- `internal/cli/doctor_test.go`: two new cases through
  `testDoctorReportsProcessResidueCase`: the mismatch error gives the `ok`
  line; a mismatch Run beside a Run with an unreadable identity gives only
  the unreadable `partial` line.
- `internal/cli/owner_pid_reused_unix_test.go` (new, `//go:build unix`): the
  parent-PID technique through `reclaimOrphanedActiveRun`, `RunSlots` and
  `releaseGoneDeliveryQueueOwner`, each with the mismatching and the
  matching identity.
- No test signals a real process; the parent PID is only read.
- Existing seams that must keep passing unchanged: the doctor residue tests,
  `TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity`,
  the `orphan_unix_test.go` and `reclaim_unix_test.go` reclaim tests,
  `TestRunSlotsExcludeAHolderWhoseOwnerIsGone`, the reconcile owner-record
  tests and the Force Stop tests.

## Build Order

1. Documentation: the doctor and deliver guides, the Roundfix setup and
   deliver references, the skill version and the Behavior Surface Record.
   It comes first so each code Task names its guide through a dependency.
2. `OwnerGone` and its tests, and Doctor's residue branch and its tests
   (depends on: 1).
3. The Active Run lock reclaim, `RunSlots` and `deliver start` on
   `OwnerGone`, with their tests (depends on: 2).
4. Final QA gate (depends on: 3).

## Risks & Considerations

- A process group member that left the owner's group but stayed in its
  session could outlive the owner while the PID is reused; Doctor already
  cannot report it today (the proof fails first), and walking the new
  holder's session would report strangers, so the design accepts that gap.
- Rebase order: Spec 0360 changes `deliver_workflow.go`, the deliver guide
  and reference and the doctor guide before this Spec; this Spec edits only
  `RunSlots` and the paragraphs it names, and regenerates every record by
  its command.
- The unix-only tests carry `//go:build unix`; the Windows compile contract
  (`TestEveryTestCompilesForWindows`) holds because no other file depends on
  them.

## Glossary

None.

## Decisions

- One store function instead of four inline checks. See ADR-0373.
- Doctor adds no line for a reused owner.
- The signal paths and the readers that already read a mismatch as gone stay
  as they are.
- Three implementation Tasks and the QA gate, in series, documentation
  first.
