---
spec: 0022-cleanup-robustness
prd: _prd.md
created: 2026-07-07
---

# Cleanup Robustness — Technical Spec

## Executive Summary

Two surgical changes: worktree removals gain `--force`, and the
Clean-outcome call sites stop treating cleanup errors as Run failures. The
trade-off accepted with `--force` is discarding worktree-local modifications
on removal — correct here by definition, because every removal site runs
only after integration succeeded or after the branch was verified to carry
no commits beyond its base; anything else keeps the worktree today and
continues to.

## System Architecture

Two modules touched, none added:

- `internal/worktree` — `CleanupClean` and `CleanupTask` run
  `git worktree remove --force`; the force-stop reap and preflight sweep
  removal paths get the same flag (all removal call sites audited).
- `internal/cli` — the post-integration cleanup call sites (`implement`,
  `resolve`/`watch` clean path) downgrade a cleanup error from
  fail-the-Run to warn-and-continue: stderr line shaped like
  `roundfix: Run Worktree cleanup failed; kept <path>: <reason>`, one
  Daemon-source Run Event, outcome and exit code untouched. Failures
  before integration keep today's behavior.

Docs-only counterpart in the roundfix and write-tasks skills.

## Implementation Design

### API Contracts

- `git worktree remove --force` everywhere Roundfix removes its own
  worktrees; branch deletion unchanged.
- Cleanup failure after a Clean integration: warning + journal, `Clean`
  outcome and exit `0` preserved; the kept worktree remains reapable by
  `stop --force` and the preflight sweep.
- No flag, config, report, or exit-code changes.

## Coverage Map

- Core Feature 1 → `internal/worktree` removal call sites
- Core Feature 2 → CLI clean-path cleanup handling + Run Event
- Core Feature 3 → roundfix SKILL.md Settle section; write-tasks skill
  authoring rules

## Testing Approach

- `internal/worktree` test: removal succeeds with untracked debris (a dirty
  file and a nested directory) in the worktree.
- `internal/cli` test: a failing cleanup after Clean integration leaves the
  stdout report and exit code byte-identical, emits the warning line, and
  journals the event.

## Build Order

1. `--force` removals plus the warn-and-continue clean path, with worktree
   and CLI tests.
2. Docs and skills: Settle recovery guidance (edit the task Verification,
   re-run settle; skip-verify rejected), write-tasks authoring rules,
   `make skills-sync` (depends on: 1).

## Risks & Considerations

- `--force` discards worktree-local changes at removal time; every call site
  runs post-integration or on verified-empty branches, so nothing
  integrable is lost. Sites that keep worktrees on non-integrated outcomes
  are untouched.

## Decisions

- Reject `settle --skip-verify` (field-report recommendation): verification
  is the only gate; the sanctioned recovery is fixing the task's
  Verification and re-running Settle, which re-reads the task file.
