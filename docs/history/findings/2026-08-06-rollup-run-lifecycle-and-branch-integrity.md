---
status: done
created_at: 2026-08-06
updated_at: 2026-09-08
kind: rollup
absorbed_by: 0125-repository-identity-and-run-branch-policy
members:
  - 2026-09-08-run-branches-ignore-the-selected-prefix.md
  - 2026-08-06-the-detach-tests-leak-the-process-they-prove-survives.md
  - 2026-08-06-three-gigabytes-of-event-journal-inside-the-retention-window.md
  - 2026-07-16-vortex-pr87-detached-watch-notification.md
  - 2026-07-17-global-run-storage-sanitation-and-compaction.md
  - 2026-07-27-owner-identity-forks-ps-and-fails-closed-under-load.md
  - 2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs.md
  - 2026-07-30-failed-qa-runs-accumulate-unreleasable-run-branches.md
  - 2026-07-30-run-termination-does-not-reach-the-acpx-child.md
  - 2026-08-02-a-spec-cycle-leaves-branches-and-worktrees-nobody-audits.md
  - 2026-08-04-branch-integrity-preflight-prescribes-a-remedy-that-reintroduces-superseded-work.md
  - 2026-08-04-watch-derives-a-review-head-it-never-checks-is-reachable.md
  - 2026-08-05-preflight-prescribes-integrating-a-superseded-run-branch.md
  - 2026-08-06-six-parallel-runs-on-one-machine-show-the-seams.md
---

# Run lifecycle and branch integrity — every created resource needs one terminal disposition (2026-08-06)

The Run findings show one lifecycle spread across process trees, Run Branches,
Task and Run Worktrees, refs, artifacts, notifications, and database storage.
Failures recur when one surface decides terminal state without disposing or
classifying the resources created on another.

## Consolidated learning

- Owner identity and termination must cover the real process tree without
  depending on a new process at the moment the host is already constrained.
- Failed and superseded Run Branches need explicit classifications;
  preflight must not prescribe integrating work that reconciliation can prove
  superseded.
- Branch and review-head checks must prove reachability against the intended
  Pull Request, not infer it from names or the global refs namespace.
- Reconciliation, GC, and notifications must report actionable terminal state
  across repositories while preserving evidence needed to audit cleanup.

## Live edge

Specs 0039, 0055, 0059, 0066, and 0068 closed important parts of the lifecycle.
The rollup remains `pending` because parallel Runs still expose cross-surface
ownership and classification seams that no single terminal audit covers.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.

## Addendum — 2026-09-08 — Current triage

Existing deliveries cover owner identity (0055), storage lifecycle (0059),
teardown (0066), close audit (0068), journal economics (0081), carry-forward
and disposition (0092, 0118), process-tree fixtures/termination (0103), and the
deleted-target content proof delivered with 0097 in commit `8235a850`.

On 2026-09-08, `TestInspectTerminalRunSafeWhenTargetDeletedAfterSquashMerge`
passed with a real Git fixture, as did the deleted-target nonmatching-content,
unknown-default, and active-Run cases. This closes the two triaged Secondbrain
captures dated 2026-09-02 about the observed deleted-target reconciliation path.
It does not prove a squash against a target that remains present: current
`internal/worktree/worktree.go` still uses ancestry then QA supersession there.
An archived Spec alone is not evidence that arbitrary Run content is disposable.

Current `internal/config/config.go` derives repository identity from the
checkout path, and `internal/store/store.go` fixes the internal branch prefix
to `roundfix/run-`. Route those gaps and the narrow reconciliation residuals to
provisional P6, repository identity and Run branch policy. Historical load
failures remain provisional P5 investigations. Keep this Rollup active to
license its 13 archived members.

## Addendum — 2026-09-08 — Branch naming decision settled

The maintainer removed the personal-prefix requirement and permits the existing
Roundfix Run/Task namespace. Purpose-based work branches are now the canonical
policy, so a configurable Run prefix is no longer required to resolve that
observation. Spec 0125 retains the separate repository-identity and reconciliation
residuals; the original namespace-conflict observation remains historical evidence.

The retired prefix-conflict Finding is absorbed here as an additional archived
member. Its runtime suggestion is deferred because the policy changed; this
Rollup remains active for the separate residuals above and now licenses 14
members.

## Addendum — 2026-09-08 — Complete implementation routing

The maintainer selected the residual work for the queue. Its primary owner is
[0125-repository-identity-and-run-branch-policy](../../specs/0125-repository-identity-and-run-branch-policy/_prd.md).
- [0124-verification-capacity-and-measured-economics](../../specs/0124-verification-capacity-and-measured-economics/_prd.md): host-load investigations.
- [0127-durable-unattended-spec-workflow](../../specs/0127-durable-unattended-spec-workflow/_prd.md): durable queue ownership and duration enforcement.

`done` records complete routing to implementation Specs, not a passing repair or
QA verdict. This Rollup remains in the active findings directory because its
archived members still name this basename as their absorption license. Those
licenses and original observations are preserved; retirement waits for the
durable replacement contract in 0120. Shipped mechanisms remain regression
obligations rather than duplicate implementation Tasks.

## Addendum — 2026-09-08 — Routing document removed

The maintainer requested removal of `docs/workflow/` and its routing documents.
The earlier citation remains a dated historical observation; its original bytes
can be read at Git revision `6b8ea48725cbca13974eee0b400b3482202874f6`.
The current primary and secondary Spec owners remain those in the complete-triage
addendum above; removing the old plan does not reopen or erase the members.

## Addendum — 2026-09-08 — Archived after complete routing

The maintainer requires terminal Findings and Rollups to leave the active
family directory. Every member now points directly to an existing active or
archived Spec, and this Rollup has its own direct Spec absorber. The earlier
statements retaining this file as an active license root are superseded by
this completed routing migration. Original observations and prior pointers
remain recorded; archival does not claim implementation of pending Specs.
