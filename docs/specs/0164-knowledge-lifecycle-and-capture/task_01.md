---
task: task_01
spec: 0164-knowledge-lifecycle-and-capture
status: completed
type: docs
complexity: low
---

# Task 01: Supersede ADR-0123 with ADR-0163

## Overview

ADR-0123 retires an orphan Review Artifact on local Git reachability, and Task 02 replaces that rule; the decision must be revised before the behaviour it governs changes. The proposed ADR-0152 already describes the replacement and is absorbed rather than left pending.

## Requirements

1. MUST add `docs/adr/0163-review-artifact-retirement-reads-recorded-evidence.md` with the lifecycle front matter of `docs/agents/docs-layout.md`, `status: accepted`, stating that an orphan Review Artifact retires on its recorded outcome, that a recorded squash merge with its merge commit is a valid receipt, that a missing outcome is an explicit unknown and retained, that no local Git object, ref or ancestry decides it, and that relocating `docs/specs/_reviews/` is independent of liveness.
2. MUST carry over, in ADR-0163, that the Review Artifact root resolver never resolves into the history root, with the literal phrase `never resolves into the history root`.
3. MUST end ADR-0163 with the literal sentence start `Supersedes ADR-0123 and ADR-0152`, followed by the date and that both are archived under `docs/history/adr/`.
4. MUST move ADR-0123 and ADR-0152 to `docs/history/adr/` with `git mv`, setting `status: superseded`, `superseded_by: ADR-0163` and a new `updated_at`, and keeping `created_at` and every body byte.
5. MUST cite no accepted ADR other than ADR-0123 and ADR-0152 in ADR-0163, so no active Spec gains an unlisted related decision.

## Subtasks

- [ ] Write ADR-0163.
- [ ] Move and mark ADR-0123 and ADR-0152 as superseded.

## Acceptance Criteria

- [ ] ADR-0163 is accepted and names both superseded decisions.
- [ ] ADR-0123 and ADR-0152 live under `docs/history/adr/`, superseded by ADR-0163, with their bodies unchanged.
- [ ] Nothing remains at the two former `docs/adr/` paths.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- creates: `docs/adr/0163-review-artifact-retirement-reads-recorded-evidence.md`
- creates: `docs/history/adr/0123-a-review-retires-on-local-reachability-and-never-writes-into-history.md`
- creates: `docs/history/adr/0152-review-artifact-retirement-uses-stable-evidence.md`

## Verification

- `new=docs/adr/0163-review-artifact-retirement-reads-recorded-evidence.md; old=docs/history/adr/0123-a-review-retires-on-local-reachability-and-never-writes-into-history.md; proposed=docs/history/adr/0152-review-artifact-retirement-uses-stable-evidence.md; grep -q "^status: accepted" "$new" || exit 1; grep -q "Supersedes ADR-0123 and ADR-0152" "$new" || exit 1; grep -q "never resolves into the history root" "$new" || exit 1; for f in "$old" "$proposed"; do grep -q "^status: superseded" "$f" || exit 1; grep -q "^superseded_by: ADR-0163" "$f" || exit 1; done; test ! -e docs/adr/0123-a-review-retires-on-local-reachability-and-never-writes-into-history.md || exit 1; test ! -e docs/adr/0152-review-artifact-retirement-uses-stable-evidence.md` — expected: exit 0; before this Task ADR-0163 does not exist and both records are still in `docs/adr/`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The decision
- `docs/history/findings/2026-08-14-a-review-retires-on-whatever-the-object-store-happens-to-hold.md`

## Result

Implementation:

- Added accepted ADR-0163 defining recorded-outcome retirement, valid recorded squash-merge receipts, explicit unknown retention, Git-independent liveness, independent `_reviews` relocation, and the live-root resolver boundary.
- Moved ADR-0123 and ADR-0152 to `docs/history/adr/`, changing only their lifecycle status, `updated_at`, and `superseded_by` metadata.

Focused checks:

- `rtk git diff --check`: passed.
- `rtk git diff --cached --find-renames --summary`: reported both ADR moves as 100% renames.
- Body SHA-256 comparisons against the pre-change ADRs: matched for ADR-0123 (`2d6cb85c35497686f4c9becfcef2d9fd261fa915761247252e214e27c03ef5c9`) and ADR-0152 (`757976b007fd9140ee9c59019962453b92bbb680a85d11327ec15026819fca2b`).
- Lifecycle front-matter scan: ADR-0163 is `accepted`; both archived ADRs are `superseded` with `superseded_by: ADR-0163` and new `updated_at` values.
- ADR reference scan: ADR-0163 names only ADR-0123 and ADR-0152.
- Former-path inventory under `docs/adr/`: no ADR-0123 or ADR-0152 files found.

Acceptance evidence:

- ADR-0163 is accepted and names both superseded decisions, including the required `never resolves into the history root` wording and terminal supersession sentence.
- Both prior ADRs are under `docs/history/adr/`; their body hashes match the original bodies.
- Neither former `docs/adr/` path remains.

The Task's declared `## Verification` commands were not run; they remain for the Daemon's verification pass.

## Carry-forward provenance

- Source Run: `run_20260924T223414Z_54e6f85ae98b6550`
- Source commit: `db8fa762f06530d3d1eed4464e36a7d4ccb9ea37`
