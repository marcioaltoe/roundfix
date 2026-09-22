---
status: approved
granted: 2026-09-21
action: run the configured Codex pre-PR reviewer over the current candidate and record its evidence, and keep the shipped skill true to it
consuming: 0153-a-reviewer-the-workflow-runs
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - internal/cli/cli_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0153

Two standing grants are consumed here.

The skill files ride the authorization of 2026-09-18, valid across the remaining
slices of this queue for one purpose: keeping the shipped skill true to the CLI
behavior the slice itself delivers.

`internal/cli/cli_test.go` rides the authorization of 2026-09-21, which covers
any governed path a slice genuinely needs, on the condition that the bounded set
is measured rather than predicted and recorded here with its reason.

## Why each governed path is unavoidable

A new command changes the public CLI surface, and this repository's hard rule
requires a pull request that changes CLI behavior to ship the skill update with
it. `internal/cli/cli_test.go` is where the command surface's own tests live, so
registering a command there is what proves it reachable.

## What is not governed

Computed as the intersection of this Spec's changed paths with the literal set
in `internal/speccheck/governed.go`. Everything else the Spec touches —
`internal/cli/review.go` and its test, `internal/cli/cli.go`,
`internal/config/config.go`, `internal/agent/agent.go` and
`docs/user-guide/commands.md` — is ordinary source that no authorization has
bounded.

## Approved bounded mutation

Register the review command on the public surface in `internal/cli/cli_test.go`
alongside the commands already covered there.

State, in the Roundfix skill and its mirror:

- that the configured Codex pre-PR reviewer is run by the workflow over the
  current candidate, and its evidence names the repository, base and head;
- that explicit `none` performs no reviewer call and no readiness probe and
  records a configured omission;
- that a runtime failure, a timeout or unreadable output blocks the selected
  mode and never becomes a pass or an omission;
- that `claude` and `coderabbit` are valid policy values this command refuses to
  execute for now, rather than reviewing nothing.

After the canonical edit, regenerate the distributed mirror with
`make skills-sync`. If any derived pin changes as a result, rewrite it only with
`make baseline-digests`.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No action, operation or path beyond those above.
- No change to how the Pre-PR Review Policy is resolved, defaulted, or reported
  by `doctor`.
- No gating of publication or merge on the evidence this Spec records.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No paid API use: every test stubs `agent.Runner` rather than starting a real
  reviewer session.
- No release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
