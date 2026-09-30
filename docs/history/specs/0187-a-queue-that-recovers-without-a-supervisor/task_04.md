---
task: task_04
spec: 0187-a-queue-that-recovers-without-a-supervisor
status: completed
type: backend
complexity: high
---

# Task 04: A Task commit is authorized by the grant it ran under

## Overview

The QA gate's mechanical authorization audit reads each governed Task commit's grant at `mechanicalAuthorizingRevision`. That is `merge-base <delivery target> <commit>^1` in `internal/speccheck/mechanical.go`, the point where the item branch forked from main. On 2026-09-29, Spec 0181's grant was widened on main (#277) and cherry-picked onto the item branch before task_07 ran. The audit still read the grant at the fork, `6784210b`, and refused. Only a rebase fixed it.

ADR-0178 keeps the fork-point read first. It also accepts the grant recorded in the commit's parent when the delivery target's current tip holds a byte-identical record at the same path. An identical record elsewhere in history never matches, so a revoked or narrowed grant is never revived. This Task implements that rule. It also describes the rule in `docs/user-guide/commands.md` and the Roundfix skill.

## Requirements

1. MUST keep the fork-point read in `detectMechanicalAuthPaths` as the first choice. When that read is not `granted`, or a changed governed path is outside its bounded paths, it MUST try the parent's grant:
   - read the blob of the authorization record at `<commit>^1`;
   - find a commit in `git log --format=%H <delivery target> -- <record path>` whose blob at that path is identical;
   - on a match, read the grant at that commit through `readMechanicalAuthorization`, and rerun the bounded-path check against it. The read's revision is that commit.
2. MUST keep the fork-point result when the parent has no record, or when no delivery-target commit holds an identical one.
3. MUST keep the self-approval check first, keep `!sameRepository` behavior unchanged, and keep every existing finding's text.
4. MUST describe the authorizing revision rule in `docs/user-guide/commands.md` and the QA section of `.agents/skills/roundfix/SKILL.md` with the phrase `the grant the Task ran under`, then run `make skills-sync`.
5. MUST rename or remove no top-level test, change no exported function signature, leave the governed `internal/speccheck/mechanical_test.go` untouched, and put the new tests in `internal/speccheck/mechanical_grant_ran_under_test.go` over temporary Git repositories.

## Subtasks

- [ ] Accept the grant the Task ran under when the delivery target holds it.
- [ ] Describe the rule in the guide and the skill, then sync the mirror.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria
- [ ] A parent grant that matches an older record but not the delivery target's current tip, because the grant was narrowed since, is not accepted.

- [ ] Main widens a grant after the item forked, and the item cherry-picks the widened record before a Task commit changes the newly bounded path. The commit is granted, and the read's revision is main's widening commit.
- [ ] A parent record that differs from every version on the delivery target refuses the path, as today.
- [ ] A Task commit that edits the authorization record is still refused as self-approval.
- [ ] A commit already covered at the fork point keeps the fork-point revision.
- [ ] The guide, the skill and its mirror carry `the grant the Task ran under`, and `make skills-sync-check` passes.

## Context

- interface: `internal/speccheck/mechanical.go`
- creates: `internal/speccheck/mechanical_grant_ran_under_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- instruction: `docs/adr/0178-a-task-commit-is-authorized-by-the-grant-it-ran-under.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAuthPathsAcceptTheGrantTheTaskRanUnder|TestAuthPathsRefuseAParentGrantMainNeverHeld|TestAuthPathsStillRefuseSelfApproval|TestAuthPathsKeepTheForkPointRevisionWhenItCovers)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAuthPathsAcceptTheGrantTheTaskRanUnder TestAuthPathsRefuseAParentGrantMainNeverHeld TestAuthPathsStillRefuseSelfApproval TestAuthPathsKeepTheForkPointRevisionWhenItCovers; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the four named tests exists, so the command fails.
- `for pair in "docs/user-guide/commands.md|the grant the Task ran under" ".agents/skills/roundfix/SKILL.md|the grant the Task ran under" "skills/roundfix/SKILL.md|the grant the Task ran under"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && make skills-sync-check` — expected: exit 0; before this Task the phrase `the grant the Task ran under` is in neither the guide nor the skill, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 4; Core Feature 4; Success Metric 4
- [_techspec.md](_techspec.md) — Grant ran under; API Contracts 1-3; Testing Approach 4; Testing Approach 5; Build Order 4
- [references/2026-09-30-a-grant-widened-mid-delivery-needs-a-rebase.md](references/2026-09-30-a-grant-widened-mid-delivery-needs-a-rebase.md)
- ADR-0178; ADR-0081; ADR-0149; ADR-0160

## Result

Implemented the authorizing-grant fallback in the mechanical authorization
audit. The audit still reads the fork-point grant first and keeps that read
when it covers the commit. For an uncovered governed path or a non-granted
read, it now compares the authorization blob in the Task commit's parent with
the blob currently held at the same delivery-target path. On an exact match it
finds the latest delivery-target commit that established those bytes, reads
the grant through `readMechanicalAuthorization`, and reruns the existing
bounded-path audit. Self-approval is checked before this fallback, and
separate-repository authorization keeps its existing path.

Documented the same rule in the command guide and the canonical Roundfix
skill, including the phrase `the grant the Task ran under`. Ran `make
skills-sync` successfully to regenerate the distributed skill mirror. The
repository-required `make baseline-digests` also succeeded and reported that
the derived artifacts already matched, with no changes.

Focused checks:

- `GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run
  '^TestAuthPathsAcceptTheGrantTheTaskRanUnder$' ./internal/speccheck` — exit 0.
- `GOCACHE=/tmp/roundfix-task04-gocache go test -count=1
  ./internal/speccheck` — exit 0 after the final test edit.
- `GOCACHE=/tmp/roundfix-task04-gocache go vet ./internal/speccheck` — exit 0.
- `cmp -s .agents/skills/roundfix/SKILL.md
  skills/roundfix/SKILL.md` — exit 0 after regeneration.
- `git diff --check` — exit 0.

Acceptance evidence from
`internal/speccheck/mechanical_grant_ran_under_test.go`:

- `TestAuthPathsRefuseAParentGrantMainNeverHeld/older_main_grant_was_narrowed`
  proves that an older byte-identical grant is not revived after main narrows
  the current record; the read keeps the fork-point revision.
- `TestAuthPathsAcceptTheGrantTheTaskRanUnder` makes the item branch diverge,
  cherry-picks main's wider grant, advances main with an unrelated commit, and
  proves the Task commit is granted with the widening commit as its revision.
- `TestAuthPathsRefuseAParentGrantMainNeverHeld/record_exists_only_on_item_branch`
  proves that an item-only parent record does not authorize the path. The
  separate `parent_has_no_record` case also proves that absence keeps the
  fork-point result.
- `TestAuthPathsStillRefuseSelfApproval` proves that a Task commit editing the
  authorization record retains the existing self-approval finding text and
  does not switch away from the fork-point read.
- `TestAuthPathsKeepTheForkPointRevisionWhenItCovers` proves that an already
  covered commit keeps the fork-point revision.
- The guide, canonical skill, and generated mirror each contain `the grant the
  Task ran under`; `make skills-sync` and the byte-for-byte mirror comparison
  succeeded. The Daemon still owns the declared `make skills-sync-check` run.

## Carry-forward provenance

- Source Run: `run_20260930T102010Z_231eaa2e314ed8d8`
- Source commit: `520d51741ba285843e0d954b387bf6206283584e`
