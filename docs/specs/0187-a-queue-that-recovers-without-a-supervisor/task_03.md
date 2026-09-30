---
task: task_03
spec: 0187-a-queue-that-recovers-without-a-supervisor
status: pending
type: backend
complexity: medium
---

# Task 03: An item started by an older queue owner records a warning

## Overview

A Delivery Queue owner is the `roundfix` binary that ran `deliver start`. On 2026-09-29 the Wave 5 owner was built at `c98b1641`. Its items started from a main that held Spec 0181, so the owner kept defects 0181 had already fixed, which caused four parks. Nothing said the owner was older than its items' main. This Task reuses the Auditor Staleness evidence. When an item starts, revalidation compares the owner's build commit with the item's starting main. When the build predates Roundfix source changes on that main, revalidation records a warning that `deliver status` prints. The queue never stops for it.

## Requirements

1. MUST add `OwnerWarning string` to `delivery.Revalidation` in `internal/delivery/engine.go`.
2. MUST make `commandDeliveryWorkflow.Revalidate` in `internal/cli/deliver_revalidate.go` run the owner check only once, at item start, before the item's first Run, when the item worktree's `HEAD` is exactly the starting main the item branch was created from. It records the result in the item's warning and never recomputes it on a retry. It MUST call `spec.ResolveAuditorEvidence(ctx, workDir, <starting main>, app.Auditor())` after the strict findings. It sets `OwnerWarning` to `owner-older-than-main: owner build <build[:12]> predates starting main <HEAD[:12]>` only when all of these hold:
   - `SelfAudit` is true;
   - `Ancestry` is `app.AncestryOlder`;
   - `git diff --name-only <build> <starting main> -- cmd internal go.mod go.sum` in `workDir` lists at least one path.

   A Git error in that diff MUST leave the warning empty and never fail revalidation.
3. MUST make `advanceItem` join a non-empty `OwnerWarning` into `item.Warning`, after any `premise-changed` text and separated by `; `. It MUST log it the way `premise-changed` is logged and never park or stop for it.
4. MUST rename or remove no top-level test and change no exported function signature. It MUST put the new tests in `internal/cli/deliver_owner_staleness_test.go`, over temporary Git repositories that set and restore `app.BuildCommit`, and in `internal/delivery/owner_staleness_test.go`, over the engine's existing fakes.
5. MUST describe the owner warning `deliver status` prints in `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md` with the phrase `owner-older-than-main`, then run `make skills-sync`.

## Subtasks

- [ ] Compute the owner warning during revalidation.
- [ ] Record, log and keep running in the engine.
- [ ] Describe the behavior in the guide and the skill, then sync the mirror.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria
- [ ] A retry of an item that already has commits does not recompute the owner check, and the recorded warning stays unchanged.

- [ ] A starting main that changed `internal/` after the build commit yields the warning, with both commits named.
- [ ] A starting main that changed only `docs/`, a build commit equal to the starting main, and a build commit absent from the repository each yield no warning.
- [ ] The engine records a `premise-changed` warning and the owner warning together, and the item still advances to `running`.
- [ ] The guide, the skill and its mirror carry `owner-older-than-main`, and `make skills-sync-check` passes.

## Context

- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_revalidate.go`
- creates: `internal/cli/deliver_owner_staleness_test.go`
- creates: `internal/delivery/owner_staleness_test.go`
- instruction: `internal/spec/auditor_evidence.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRevalidateWarnsWhenTheOwnerPredatesSourceOnMain|TestRevalidateStaysQuietForDocsOnlyMain|TestRevalidateStaysQuietForACurrentOwner|TestRevalidateStaysQuietWithoutTheBuildCommit)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRevalidateWarnsWhenTheOwnerPredatesSourceOnMain TestRevalidateStaysQuietForDocsOnlyMain TestRevalidateStaysQuietForACurrentOwner TestRevalidateStaysQuietWithoutTheBuildCommit; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the four named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^TestAdvanceItemRecordsTheOwnerWarningAndKeepsRunning$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestAdvanceItemRecordsTheOwnerWarningAndKeepsRunning"` — expected: exit 0; before this Task the named test does not exist, so the command fails.
- `for pair in "docs/user-guide/commands.md|owner-older-than-main" ".agents/skills/roundfix/SKILL.md|owner-older-than-main" "skills/roundfix/SKILL.md|owner-older-than-main"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && make skills-sync-check` — expected: exit 0; before this Task the phrase `owner-older-than-main` is in neither the guide nor the skill, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; Core Feature 3; Success Metric 3
- [_techspec.md](_techspec.md) — Owner warning; API Contract 1; Testing Approach 3; Build Order 3
- [references/2026-09-30-a-queue-owner-older-than-its-starting-main-runs-silently.md](references/2026-09-30-a-queue-owner-older-than-its-starting-main-runs-silently.md)

## Result
