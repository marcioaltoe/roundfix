---
spec: 0300-reconcile-releases-a-stale-run-owner-record
prd: _prd.md
created: 2026-10-22
---

# Reconcile releases a stale Run owner record — Technical Spec

## Executive Summary

`OwnerProcessControl.proveOwner` in `internal/store/process.go` returns the
same `ErrOwnerProcessIdentityUnproven` for a comparable start identity that
differs from the recorded one, for PID 0 and for the inspecting process's own
PID. `inspectReconcileProcesses` in `internal/cli/reconcile.go` turns every
`InspectTree` error into a preserved candidate, so a terminal Run whose owner
PID the system gave to another process is preserved forever. A new sentinel,
`ErrOwnerProcessIdentityMismatch`, wraps the unproven error for the comparable
mismatch only, with the same message bytes. Reconcile reports that case as a
`runOwnerRecord` process candidate, and `--apply` releases it through a new
store method, `ReleaseStaleRunOwnerRecord`, which clears the Run's owner PID
and identity by compare-and-set and journals what it cleared. The release
never reaches the process controller, so it cannot send a signal.

The trade-off accepted: a released Run no longer names the process that owned
it in its row; the Run Event keeps the PID and identity, and nothing that reads
a terminal Run's owner needs them once the owner is proven gone.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the release clears two existing columns and journals an existing
  event kind. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary homes,
  repositories and fake process tables. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 1 and 2 follow
  ADR-0313, named with its title in `_prd.md` → Project
  Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AO approval of 2026-10-22 ("the same rules as cycle AN"), carrying the cycle
  AJ grant "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0300-reconcile-releases-a-stale-run-owner-record/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at `34b4e246`. Line numbers move when
other Specs merge; the function names do not.

1. `proveOwner` (`internal/store/process.go:112`) returns
   `ErrOwnerProcessIdentityUnproven` for a PID that is not positive or equals
   `os.Getpid()` (`:113-115`), and again when the live start identity was read,
   carries the `runtime.GOOS` prefix and differs from the recorded one
   (`:138-141`, message `live process start identity does not match the
   recorded owner identity`). A failed identity read, unless the process is
   then absent, and a recorded token without the platform's prefix return
   `ErrOwnerIdentityUnreadable` (`:131-137`). `InspectTree` (`:208`) returns
   the proof error before it enumerates anything; `TerminateAndWait` (`:150`)
   returns it before any signal.
2. Characterization with a fake process table, a `go test -overlay` probe in
   `internal/store` that wrote nothing to the repository: for a comparable
   different identity, the own PID and PID 0, `InspectTree` returned an error
   matching `ErrOwnerProcessIdentityUnproven` and not
   `ErrOwnerIdentityUnreadable`; for a read failure and a foreign token, the
   reverse. No case enumerated or signalled. The comparable mismatch is not
   distinguishable from the other two by any sentinel.
3. `inspectReconcileProcesses` (`internal/cli/reconcile.go:758`) skips Runs
   without an owner PID, preserves an Active Run's owner and a Run recorded
   with `OwnerIdentityUnproven`, then calls `InspectTree` (`:799`) and turns
   any error into a `preservedCandidates` entry of kind `process` with the
   refusal reason `owned process tree cannot be proven: <error>` (`:800-809`).
   An absent owner with no live tree adds nothing.
4. Characterization through the command, a second overlay probe in
   `internal/cli`: a live child process started by the test held the recorded
   PID, the Run recorded `darwin:1.1` as identity, and the real process table
   was used. `roundfix reconcile`, `--format=json` and `--apply --format=json`
   each exited 0 with `preserved=1`, `process-candidates=0`, and the refusal
   reason of the adopted Backlog Entry; the Run Database bytes were unchanged,
   the owner PID was still recorded, and the child was alive.
5. `applyReconcileProcesses` (`:1130`) calls `TerminateTreeAndWait` for every
   `processCandidates` entry. `printReconcileText` prints `Process candidate:
   Run <id> owner PID <pid>` per entry and one `Debris summary:` line (`:1344`)
   whose fields `TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec`
   (`internal/cli/reconcile_item_branch_test.go:65`) and
   `internal/cli/reconcile_staging_test.go:41` pin by substring.
   `reconcileDebrisSummary` has no field pinned by name in JSON.
6. The store clears an owner only for the Delivery Queue
   (`ReleaseDeliveryQueueOwner`, `internal/store/delivery.go:343`, a
   compare-and-set on PID and identity). Terminal reconciliation journals a
   `daemon.outcome` event with a payload `event` field
   (`ReconcileIntegration`, `internal/store/store.go:755`).

Sweep of every reader of the owner proof errors and of a terminal Run's owner
(`rg -n 'ErrOwnerProcessIdentityUnproven|ErrOwnerIdentityUnreadable'
internal -g '!*_test.go'` and `rg -n '\.OwnerPID' internal -g '!*_test.go'`):

- Force Stop (`internal/cli/cli.go`, `forceStop…` around `:1029-1047`)
  refuses any unproven owner and refuses `--owner-identity-unreadable` over
  `ErrOwnerProcessIdentityUnproven`; the mismatch wraps that error, so both
  stay as they are.
- `TerminateAndWait`'s post-grace re-proof (`process.go:179-183`) already
  treats an unproven owner as "a reused PID is a stranger" and returns; it
  keeps its behavior.
- The Delivery Queue owner checks in `internal/cli/deliver.go` (`:184`,
  `:466`, `:533`) already treat an unproven queue owner as proven gone
  (the `deliver status` owner rule); they read the queue's owner, not a Run's, and are left
  unchanged.
- `SignalProvenProcess` and `SignalProvenProcessGroup`
  (`internal/store/process_signal_unix.go`) return the unproven sentinel only
  for invalid arguments and send nothing on a mismatch; unchanged.
- `doctor`'s `Residue` (`internal/cli/doctor.go:275`) reports a lineage read
  error as `partial`; it skips a Run without an owner PID, so a released
  record leaves it. Unchanged.
- The carry-forward staging sweep already treats an `owner.json` identity
  that differs from the live one as stale (`internal/worktree/staging.go`);
  unchanged.
- `RunSlots` (`deliver_workflow.go`), `activeOwnerPID` and implement's owner
  recording read Active Runs only; a release touches terminal Runs only.
- Text: `docs/user-guide/commands/reconcile.md` and
  `.agents/skills/roundfix/references/reconcile.md` describe `process`,
  `runBranch` and `staging` candidates and never an owner record; the guide's
  staging paragraph uses "owner record" for the staging `owner.json`, which
  is why the new term says "Run owner record".

## System Architecture

- `internal/store/process.go`: `ErrOwnerProcessIdentityMismatch`; `proveOwner`
  returns it for the comparable mismatch.
- `internal/store/run_owner_record.go` (new): `ReleaseStaleRunOwnerRecord`.
- `internal/cli/reconcile.go`: the `runOwnerRecord` candidate in
  `inspectReconcileProcesses`, its release in `applyReconcileProcesses`, the
  text heading and the `runOwnerRecordsReleased` count.
- `docs/user-guide/commands/reconcile.md`, the Roundfix `reconcile` reference
  and its mirror, the Roundfix skill's versions, `CONTEXT.md`.

No package, import edge or schema changes.

## Implementation Design

### Interfaces

```go
// ErrOwnerProcessIdentityMismatch wraps ErrOwnerProcessIdentityUnproven for a
// live process whose comparable start identity differs from the recorded one.
var ErrOwnerProcessIdentityMismatch error

// ReleaseStaleRunOwnerRecord clears a terminal Run's recorded owner PID and
// identity when they still equal pid and identity, and journals the release.
func (store *Store) ReleaseStaleRunOwnerRecord(ctx context.Context, runID string, pid int, identity string, reason string) (bool, error)
```

### Data Models

None changes shape. The Run row keeps its columns; a release sets
`owner_pid` to NULL, `owner_identity` to `''` and `owner_identity_unproven` to
0. `reconcileDebrisSummary` gains `RunOwnerRecordsReleased`
(`runOwnerRecordsReleased`); `reconcileDebrisResult` is unchanged.

### API Contracts

1. API Contract: the mismatch sentinel. `ErrOwnerProcessIdentityMismatch` is
   `fmt.Errorf("%w: live process start identity does not match the recorded
   owner identity", ErrOwnerProcessIdentityUnproven)`, and `proveOwner`
   returns `ownerProcessControlError(pid, "prove owner process identity",
   ErrOwnerProcessIdentityMismatch)` where it built that message inline. The
   message bytes are unchanged; `errors.Is` matches both sentinels. The own
   PID and a PID that is not positive keep returning the bare unproven error;
   every unreadable path keeps `ErrOwnerIdentityUnreadable`.
2. API Contract: the release. `ReleaseStaleRunOwnerRecord` trims its inputs
   and refuses an empty Run ID, identity or reason and a PID below 1. In one
   write transaction it reads the Run's state, refusing a missing Run and a
   non-terminal one (`release stale Run owner record of Run "<id>": Run is
   <state>, not terminal`); it runs `UPDATE runs SET owner_pid = NULL,
   owner_identity = '', owner_identity_unproven = 0 WHERE id = ? AND owner_pid
   = ? AND owner_identity = ?`; with no row affected it returns `false, nil`
   and writes nothing else; with one it appends a `daemon.outcome` Run Event
   from the `daemon` source whose payload is `{"event":
   "stale_run_owner_record_release", "owner_pid": <pid>, "owner_identity":
   <identity>, "reason": <reason>}` and whose summary is `Run released its
   Stale Run Owner Record: <reason>.`, and returns `true, nil`. It changes no
   state, `updated_at` or `completed_at`, and touches no process.
3. API Contract: the candidate. When `InspectTree` returns an error matching
   `ErrOwnerProcessIdentityMismatch` for a selected terminal Run,
   `inspectReconcileProcesses` appends to `processCandidates`, not to
   `preservedCandidates`, an entry with `kind` `runOwnerRecord`, the Run ID
   and outcome, `ownerPid`, `processIds` `[]`, `proof` `terminal Run "<id>"
   outcome "<outcome>" records owner PID <pid>, which now has a different
   start identity: the recorded owner ended`, and `action` `would release the
   Run owner record with --apply` (with `--apply`, `release the Run owner
   record without a signal` until applied). Every other error, and every case
   before `InspectTree`, keeps today's entry. `debrisSummary.processCandidates`
   counts the entry.
4. API Contract: the apply. `applyReconcileProcesses` never passes a
   `runOwnerRecord` entry to the process controller. It opens the Run Database
   once when at least one such entry exists and calls
   `ReleaseStaleRunOwnerRecord` with the Run ID, the owner PID, the inspected
   Run's owner identity and the reason `reconcile proved owner PID <pid> now
   has a different start identity`. On `true` the action becomes `released Run
   owner record; no signal sent` and `runOwnerRecordsReleased` rises by one;
   on `false`, `none; Run owner record already released`; on an error, the
   action is `preserve`, the refusal reason `release Run owner record:
   <error>`, and `debrisSummary.preserved` and `summary.operationalFailures`
   rise by one, as a failed process termination does today.
5. API Contract: the text report. A `runOwnerRecord` entry prints under the
   heading `Run owner record candidate: Run <id> owner PID <pid>` with the
   same four lines as a process candidate (`process-ids: []`, `proof`,
   `action`, `refusal-reason`). The `Debris summary:` line gains
   `run-owner-records-released=<n>` right after `processes-applied=<n>`, so
   every substring the existing tests pin is unchanged.
6. API Contract: unchanged callers. Force Stop, `TerminateAndWait`,
   `SignalProvenProcess`, `SignalProvenProcessGroup`, the Delivery Queue owner
   checks, `doctor` and the staging sweep keep their bytes and behavior.
7. API Contract: documentation. The `reconcile` guide and the Roundfix
   `reconcile` reference describe the `runOwnerRecord` candidate, its release
   and its count; the glossary adds **Stale Run Owner Record** and revises
   **Reconcile Command**.

### Surface Transcripts

None. No command, flag, exit code or stream placement changes. The text
report gains one heading form and one summary field, and the JSON report one
`kind` value and one summary field; the tests of API Contracts 3-5 assert the
exact lines and fields, and every pinned existing line keeps its bytes.

### Invariants

```text
1. The release path never calls the process controller; no signal reaches any process from it.
2. Only a comparable start identity that differs from the recorded one releases a record; unreadable, unsupported, invalid-PID, own-PID and identity-less records stay preserved.
3. The mismatch keeps today's message bytes and still matches ErrOwnerProcessIdentityUnproven, so Force Stop refuses it.
4. A release clears only the PID and identity that were inspected, only for a terminal Run, and journals exactly one event.
5. A dry run writes nothing to the Run Database.
6. The tree of a process that holds a mismatched PID is never enumerated.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 2, Success Metric 1, Invariants 5-6 → the candidate and the text report (API Contracts 3 and 5)
- Goal 2, Core Feature 3, Success Metrics 2 and 6, Invariants 1 and 4 → the release and the apply (API Contracts 2 and 4)
- Goal 3, Success Metric 3, Invariant 2 → the candidate's error test and unchanged preserved entries (API Contract 3)
- Goal 4, Core Feature 1, Success Metrics 4-5, Invariant 3 → the mismatch sentinel and unchanged callers (API Contracts 1 and 6)
- Goal 5, Core Feature 4 → the guide, the reference and the glossary (API Contract 7)
- Success Metric 7 → the PRD's Unreachable Acceptance

## Absence Evidence

1. `ErrOwnerProcessIdentityMismatch` — command: `grep -rn -e IdentityMismatch -e 'ReleaseStale' internal`; result: no match at `34b4e246`; the mismatch is built inline in `proveOwner` on the unproven sentinel.
2. A release of a Run's owner — command: `grep -rn -e 'owner_pid = NULL' internal/store`; result: only `internal/store/delivery.go`, the Delivery Queue owner release; no statement clears a Run's owner.
3. An owner record candidate kind — command: `grep -rn -e runOwnerRecord -e 'owner-records-released' internal docs .agents skills`; result: no match at `34b4e246`.

## ADR Alternatives

- ADR-0313: no rejected alternative is taken; the preserved entry is not kept, the mismatch is not hidden as an absent owner, apply adds no second proof, the Run is not deleted, and no read command releases the record.

## Integration Points

The local process table, read through the store's owner proof, and the Run
Database. No test, Verification command or QA row reads the real
`~/.roundfix`, signals a process it did not start, or reaches the network.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary homes and repositories, and passes no real-time duration (the
Real-Time Budget rule gives a new file a ceiling of zero): the controller is
built with `NewOwnerProcessController()` and its process-table functions are
replaced in the test.

- `internal/store/process_identity_mismatch_test.go` (new):
  `TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity`,
  a fake process table with the rows: a comparable different identity
  (mismatch and unproven, not unreadable, the exact message, and
  `TerminateAndWait` returns the mismatch with no signal); a read failure and
  `ErrOwnerProcessUnsupported` from the identity read (unreadable, not
  mismatch); a foreign platform token (unreadable, not mismatch); the own PID
  and PID 0 (unproven, not mismatch); an owner absent on the second liveness
  read (no error). The table's enumerate and signal functions fail the test
  whenever the proof refused.
- `internal/store/process_identity_mismatch_unix_test.go` (new, `//go:build
  unix`): `TestInspectTreeReportsAProvenMismatchOfALiveProcessWithoutSignallingIt`,
  a live helper from `startOwnerProcessHelper` inspected with a recorded
  identity of the platform's prefix: the mismatch, no process attributed, the
  helper still alive.
- `internal/store/run_owner_record_test.go` (new):
  `TestReleaseStaleRunOwnerRecordClearsOnlyTheRecordedTerminalOwner`: an
  Active Run is refused and kept; a terminal Run with a different PID or
  identity returns `false` and journals nothing; the inspected values clear
  the owner, keep the state, `updated_at` and `completed_at`, and journal one
  event with the payload of API Contract 2; a second call returns `false` and
  journals nothing.
- `internal/cli/reconcile_run_owner_record_test.go` (new), with a fake
  process controller that answers each owner PID with a fixed proof error and
  records every termination request:
  - `TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal`: the text
    dry run prints the heading and four lines of API Contract 5 and the
    summary substring ` processes-applied=0 run-owner-records-released=0
    run-branches-applied=0 `, with no preserved candidate; the JSON dry run
    holds one `runOwnerRecord` entry with `"processIds":[]` and
    `"runOwnerRecordsReleased":0`; the Run Database bytes are unchanged;
    `--apply` releases, asks for no termination and clears the owner; a second
    `--apply` inspects nothing and reports nothing; `Residue` then reads no
    lineage for the Run.
  - `TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale`:
    subtests for an unreadable read failure, an uncomparable token, a bare
    unproven error, `ErrOwnerProcessUnsupported` and a Run recorded without an
    identity; each stays in `preservedCandidates` with today's reason through
    a dry run and `--apply`, no termination is requested, and the record is
    kept.
- Unchanged and run by name: `TestOwnerProcessControllerProveOwnerClassifiesIdentityEvidence`,
  `TestOwnerProcessControllerRefusesMismatchedOwnerIdentity`,
  `TestOwnerProcessControllerProveOwnerRefusesMismatchedIdentity`,
  `TestRunReconcileOffersAndAppliesOwnedProcessTreesAndRunBranches`,
  `TestRunReconcileApplyPreservesProcessWhenTerminationIsUnproven`,
  `TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec`,
  `TestRunForceStopOwnerIdentityUnreadableFlagRequiresUnreadableProof` and
  `TestDoctorReportsProcessResidueUnreadable`.
- No existing test, golden or fixture changes.

## Build Order

1. The `reconcile` guide, the Roundfix `reconcile` reference, the Roundfix
   skill's versions, the Behavior Surface Record, **Stale Run Owner Record**
   and **Reconcile Command** (no dependency: the text states the contracts of
   `_techspec.md`, and its Verification reads only documents).
2. The mismatch sentinel, the release, the candidate, the apply, the text
   report and their tests (depends on: 1, which writes the guide and the
   reference for the surface this step changes, so the code and its text land
   in one chain).
3. QA gate (depends on: 1, 2).

No Spec is listed in `requires`. The cycle AO Specs 0297, 0298 and 0299
deliver first and raise the Roundfix skill's version; this Spec rebases on
them and touches none of their files apart from the skill's version fields.

## Risks & Considerations

- A host whose stored identity token could differ from the live token of the
  same process would release a live owner's record; the platform prefix
  already refuses foreign and legacy tokens, and the path sends no signal, so
  the worst case is a terminal Run that stops naming a process `doctor` could
  have reported as residue.
- On Linux the token is `/proc/<pid>/stat` `starttime`, which proc_pid_stat(5)
  does not mask on a ptrace denial; a denied read fails and stays unreadable.
- A reconcile run inside another PID namespace would read another process at
  the same number; Roundfix records and reconciles Runs from the host's
  Roundfix Home, and the Linux Lane container never reads it.
- The journal event uses `daemon.outcome` like terminal reconciliation, so
  `roundfix events` and the TUI show one more outcome line for the Run.

## Glossary

- adds: **Stale Run Owner Record**
- changes: **Reconcile Command**

## Decisions

- A comparable different start identity is its own sentinel, wrapping the
  unproven one; see ADR-0313.
- The Run owner record is released by compare-and-set and journaled, never by
  a signal and never by deleting the Run; see ADR-0313.
- The candidate lives in `processCandidates` with the `kind`
  `runOwnerRecord`, and the summary count is placed so existing pinned
  substrings keep their bytes.
- `doctor` is left unchanged: after a release it no longer reads the Run.
