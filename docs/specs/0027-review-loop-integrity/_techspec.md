---
spec: 0027-review-loop-integrity
prd: _prd.md
created: 2026-07-14
---

# Review Loop Integrity — Technical Spec

## Executive Summary

Review Runs move from Run Worktree isolation to direct execution in the user's checkout, guarded by a new Branch Integrity Preflight. The primary trade-off this design accepts: we give up worktree isolation's tolerance for dirty checkouts and its crash containment in exchange for eliminating the stranded-work and Integration Pending failure modes entirely — the fix delta is always applied to, committed on, and pushed from the branch the pull request actually serves. The change is mostly seam redirection, not new machinery: the Daemon's engine already takes its working directory as a plan field, the Final Push already runs from the user root, and fetch already runs without a worktree. New capability concentrates in three places: a deterministic preflight (worktree/lock enumeration plus fast-forward auto-integration), a Merge-Ready confirmation grace window ending in a new `CleanUnverified` terminal state, and a Review Source write surface extended from "resolve thread" to "resolve, reply, and comment" so outcomes become auditable on GitHub.

## System Architecture

Existing modules extended, in glossary vocabulary:

- **`internal/cli`** — owns command preflight blocks, flag parsing, and reports. Gains the Branch Integrity Preflight step (wired into fetch, resolve, and watch before any Run is created), the `--skip-branch-integrity` bypass with its audit comment, the removal of the review-run worktree path, and the per-Run vs cumulative report split.
- **`internal/preflight`** — stays git-inspection-only; the Branch Integrity Preflight composes at the CLI layer so `preflight` does not grow dependencies on `store` and `worktree`.
- **`internal/worktree`** — gains a read-only enumeration helper: all `roundfix/run-*` branches and their worktrees whose base is a given branch, with ahead-commit detection (reuses the `for-each-ref` and `worktree list --porcelain` parsing already present for terminal pruning).
- **`internal/store`** — gains the `CleanUnverified` terminal state and its journal/terminal-set registration. Lock semantics are unchanged; the preflight consumes the existing active-run lookups.
- **`internal/watch`** — `confirmMergeReady` treats a missing check as pending within a configurable grace window instead of returning ready-with-warning immediately; exhausting the window yields a distinct confirmation result the CLI maps to `CleanUnverified`.
- **`internal/daemon`** — the engine's `CyclePlan.GitRoot` is fed the user root instead of the worktree path; batch-settlement source resolution extends from "resolve resolved/invalid threads" to per-issue propagation with Outcome Comments and failure reasons.
- **`internal/rounds`** — issue artifact frontmatter gains a terminal-reason field, threaded through `SetIssueStatus`.
- **`internal/reviewsource` / `coderabbit`** — the `GitHubClient` boundary gains thread-reply and PR-comment mutations next to the existing resolve mutation, all through the same `gh` runner.
- **`skills/roundfix`** — SKILL.md and the OpenAI manifest updated to the new contract; `skills.Check` phrase anchors updated in the same change.

## Implementation Design

### Interfaces

Branch Integrity Preflight (new, `internal/cli`, composing `worktree` + `store`):

```go
// PendingRunWork is one roundfix/run-* branch with commits ahead of base.
type PendingRunWork struct {
    Branch       string // roundfix/run-<id>
    WorktreePath string // empty if the worktree is gone
    AheadCommits int
    FastForward  bool // base can ff to this branch tip
}

type BranchIntegrityReport struct {
    Pending    []PendingRunWork
    Integrated []PendingRunWork // auto-integrated this preflight
    ActiveRun  *store.Run       // run holding a lock on the target
}

func inspectBranchIntegrity(ctx, runner gitcmd.Runner, st *store.Store,
    git preflight.GitState, target lockTargetKeys) (BranchIntegrityReport, error)
```

Review Source write surface (extends `coderabbit.GitHubClient`):

```go
type GitHubClient interface {
    // ...existing reads + ResolveReviewThread(ctx, threadID) error
    ReplyToReviewThread(ctx context.Context, threadID, body string) error
    CommentOnPullRequest(ctx context.Context, pr int, body string) error
}
```

Merge-Ready confirmation (changed shape, `internal/watch`):

```go
type confirmResult struct {
    ready      bool
    unverified bool // grace window exhausted with no check observed
    timedOut   bool
}
// CheckMissing now loops like CheckPending until req.CheckGracePeriod elapses.
```

Issue settlement with reason (`internal/rounds`):

```go
// SetIssueStatus rewrites status, duplicate_of, and terminal_reason.
func SetIssueStatus(path, status, duplicateOf, terminalReason string) error
```

### Data Models

- **Run state**: new `StateCleanUnverified` in `internal/store`, added to the terminal set, journal state list, and both resolve/watch exit mappings. No schema migration — states are strings.
- **Issue artifact frontmatter**: new optional `terminal_reason` string (empty for resolved issues), mirrored on `rounds.Issue`. Written when the Daemon settles an issue `failed` (failed step, command, exit status, diagnostics path) and when triage yields `invalid`/`duplicated` (the reason/canonical reference used for the Outcome Comment).
- **Run row**: review Runs no longer set `WorkDir` to a worktree path; it records the user root. No column changes.
- **Config**: new `watch.check_grace_period` (duration, default 5m) in the `Watch` group, following the existing overlay/default/`DefaultConfigYAML` pattern. The PRD's open-question default (reuse the quiet period) is refined here: 30s is shorter than a Review Source's typical re-analysis startup, so a dedicated key with a 5m default keeps the semantics honest.

### API Contracts

Public CLI surface changes (breaking-change review applies):

- **New flag** `--skip-branch-integrity` (bool) on fetch, resolve, and watch. Using it requires the audit comment to publish; a publish failure fails the command at preflight with exit 2 and no side effects beyond the attempted comment.
- **New terminal outcome** `CleanUnverified` with **exit code 3** (0 = Clean/OK, 1 = run failed/unresolved, 2 = preflight, 130 = SIGINT are unchanged). Watch maps it to 3; Clean stays 0.
- **Removed for review Runs**: the Integration Pending outcome and kept-worktree messaging. `IntegrationPending` remains for implement.
- **Report format**: the summary line splits into two labeled lines — this Run's counts (from the Run's own selection/accumulated issues) and the pull request's cumulative counts (from the existing disk scan). Failed and unresolved issue lines gain a one-line reason suffix sourced from `terminal_reason`.
- **Preflight failure text** follows the existing `printPreflightFailure` convention: reason, "No side effects" statement, and the exact unblocking command per pending item (`git merge --ff-only <branch>` equivalent or `roundfix stop --run-id <id> [--force]`).

## Coverage Map

- Goal 1, Stories 1–3 → Branch Integrity Preflight (`internal/cli` + `worktree` enumeration + `store` lookups), checkout execution (CyclePlan.GitRoot seam), ADR-0042/0045.
- Goal 2, Story 5 → Merge-Ready grace window (`internal/watch`), `StateCleanUnverified` (`internal/store`), exit mapping (`internal/cli`), ADR-0043.
- Goal 3, Story 6 → per-issue propagation + Outcome Comments (`internal/daemon` engine batch settlement, `coderabbit` write surface).
- Goal 4, Story 8 → report split (`internal/cli` report rendering).
- Goal 5, Story 7 → `terminal_reason` persistence (`internal/rounds`, engine settle paths).
- Story 4 → bypass flag + audit comment (`internal/cli`, `coderabbit.CommentOnPullRequest`).
- Core Feature 10 → skill/glossary sync (`skills/roundfix`, `CONTEXT.md`, `skills.Check` anchors).

## Integration Points

- **GitHub via `gh`** (existing adapter pattern in `coderabbit`): two new GraphQL/REST mutations — `addPullRequestReviewThreadReply` for Outcome Comments on threads, and a PR-level comment for the bypass audit. Idempotency: each comment body carries a stable marker line (run id + issue fingerprint or `bypass` tag); before posting, the already-fetched thread data is checked for the marker and posting is skipped on match.
- **Git** (existing `gitcmd.Runner`): branch enumeration, merge-base/ancestry checks, and `merge --ff-only` for auto-integration — porcelain-only, consistent with ADR-0024.
- **Review Source check** (existing `CheckSource` seam): unchanged interface; only the missing-state policy changes.

## Testing Approach

All existing seams; no new test seams are introduced.

- **Unit (table tests)**: `confirmMergeReady` grace-window matrix (missing→appears, missing→exhausted, pending→timeout) via the existing fake clock/sleeper/CheckFunc; branch-integrity report construction against a fake `gitcmd.Runner` script; `rounds.SetIssueStatus` reason round-trip; comment idempotency marker matching; exit-code mapping including `CleanUnverified`.
- **Integration (buffer-captured CLI runs, `Run(args, &stdout, &stderr) int`)**: fetch/resolve/watch preflight refusal output and exit 2 with pending worktrees and with an active lock; bypass path with a fake GitHubClient asserting the audit comment body and the fail-on-publish-failure rule; resolve cycle in a temp git repo asserting commits land on the user branch with no `roundfix/run-*` branch created; report rendering with distinct per-Run and cumulative fixtures; dirty-tracked-tree refusal (ADR-0045).
- **Engine tests** (existing daemon test style): per-issue propagation invoked at batch settlement with statuses resolved/invalid/failed/duplicated, asserting resolve-vs-reply-vs-leave-open per the PRD decision table and that no thread resolves before Verification passes.
- **Skill check**: `roundfix skills check` phrase anchors updated with the SKILL.md edits in the same step, keeping `make verify` green.

## Build Order

1. **Store state + exit mapping groundwork**: add `StateCleanUnverified` to states, terminal set, journal list; extend watch/resolve exit mapping with exit 3 (mapping is dead until step 6 produces the state).
2. **Issue artifact terminal reason**: `terminal_reason` frontmatter field, `SetIssueStatus` signature change, update existing callers passing empty reasons.
3. **Review Source write surface**: `ReplyToReviewThread` + `CommentOnPullRequest` on `GitHubClient` and `GHClient`, with idempotency markers and fakes.
4. **Branch Integrity Preflight** (depends on: 3): enumeration helper in `worktree`, report construction, ff auto-integration, active-lock lookup, `--skip-branch-integrity` flag, audit comment via 3, preflight failure texts; wired into fetch, resolve, watch.
5. **Checkout execution for resolve/watch** (depends on: 4): pass the user root as `CyclePlan.GitRoot`; delete the review-run worktree creation, integration, Integration Pending, and kept/cleanup messaging paths; add the clean-tracked-tree preflight check (ADR-0045); failed-batch report states that dirty paths are Agent work.
6. **Merge-Ready grace window + CleanUnverified** (depends on: 1): `check_grace_period` config key, `confirmMergeReady` missing-as-pending loop, `unverified` result threading through watch to the new outcome, report wording naming the next action.
7. **Per-issue propagation + Outcome Comments** (depends on: 2, 3): extend batch settlement to propagate every settled issue individually — resolve after comment for invalid/duplicated, reply-and-leave-open for failed, plain resolve for resolved — and post run-end comments for still-unresolved issues; journal each propagation with the Review Issue reference.
8. **Terminal reasons from the engine** (depends on: 2): settle paths pass the failed step/command/exit status/diagnostics path into `SetIssueStatus`; triage reasons captured for invalid/duplicated.
9. **Report split** (depends on: 1, 8): per-Run counts from the Run's own issue set, cumulative counts from the disk scan, labeled separately; reason suffixes on failed/unresolved lines.
10. **Skill, glossary, and docs sync** (depends on: 4, 5, 6, 7, 9): SKILL.md review sections, OpenAI manifest, `skills.Check` anchors, `CONTEXT.md` Run Worktree/Fetch Run/Integration Pending definitions scoped to spec Runs, command help text.

## Risks & Considerations

- **Failed batches dirty the user checkout.** Mitigated by ADR-0045 (clean tree at start makes all dirt Agent work) and explicit report wording; the user decides to keep, fix, or discard. Never auto-reset.
- **Comment volume on large Rounds.** Only non-resolved outcomes comment (PR #17's round would have posted 2 comments, not 50). The idempotency marker prevents retry duplication.
- **Grace window extends watch wall-clock** by up to `check_grace_period` when the Review Source genuinely never checks. Bounded, configurable, and strictly better than a false Clean.
- **Auto ff-integration mutates the user branch during preflight.** It only fast-forwards — no merge commits, no conflict resolution — and each integration is reported; anything non-ff blocks instead.
- **Exit code 3 is new public API.** Scripts checking `== 0` keep working (CleanUnverified is not success); scripts checking `!= 0` now distinguish it from failure by value.
- **Lock semantics for orphaned owners** stay as-is until spec 0028; the preflight's active-run refusal names the force stop form for dead-owner cases.

## Decisions

- Review Runs execute in the user checkout; the worktree path is deleted for resolve/watch, not flagged off. See ADR-0042.
- Resolve/watch preflight requires a clean tracked working tree; untracked files allowed. See ADR-0045.
- `CleanUnverified` gets exit code 3; existing codes 0/1/2/130 unchanged. See ADR-0043 for the outcome itself.
- Outcome propagation happens per issue at batch settlement — the earliest point compatible with "no thread resolves before Verification passes"; this satisfies the PRD's never-later-than-batch bound without a mid-batch artifact watcher.
- Branch Integrity Preflight lives in `internal/cli`, composing `worktree` and `store`, keeping `internal/preflight` dependency-free.
- Grace period is a dedicated `watch.check_grace_period` key (default 5m), refining the PRD's reuse-the-quiet-period default because 30s undercuts real Review Source startup latency.
- The bypass flag is `--skip-branch-integrity` — self-describing, not an overloaded `--force` (PRD open question closed).
