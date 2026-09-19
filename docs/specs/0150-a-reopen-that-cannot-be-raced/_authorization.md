---
status: approved
granted: 2026-09-18
action: keep the shipped Roundfix skill true to the recheck refusal this Spec delivers
consuming: 0150-a-reopen-that-cannot-be-raced
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0149

On 2026-09-18 the maintainer granted a standing authorization for the Roundfix
skill and its distributed mirror, valid across the remaining slices of this
queue, for one purpose: keeping the shipped skill true to the CLI behavior the
slice itself delivers. This record consumes that authorization for Spec 0149 and
bounds it to this Spec's change.

## Why a governed path is unavoidable

The Spec adds a command to the public CLI surface, and this repository's hard
rule requires a pull request that changes CLI behavior to ship the skill update
with it. The skill is the shipped description of that surface, so a new command
that the skill does not name is a skill that is no longer true.

## What is not governed

Every other path this Spec touches was classified through the governance probe
rather than estimated: `internal/cli`, `internal/spec`, their tests and
`docs/user-guide` are ordinary source that no authorization has bounded.

## Approved bounded mutation

State, in the Roundfix skill and its mirror:

- that a settled QA gate whose dependencies are no longer completed is cleared
  with the reopen command, never by editing the Task file;
- that the command refuses when the gate is not settled or is not stale;
- that the prior QA Report and Result are retained, and the invalidation is
  recorded in the QA Task.

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
- No change to the skill's version, or to any behavior it documents beyond this
  Spec's own.
- No change to when the loader raises its stale-gate refusal.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
# Approved authority for Spec 0150

On 2026-09-18 the maintainer granted a standing authorization for the Roundfix
skill and its distributed mirror, valid across the remaining slices of this
queue, for one purpose: keeping the shipped skill true to the CLI behavior the
slice itself delivers. This record consumes that authorization for Spec 0150 and
bounds it to this Spec's change.

## Why a governed path is unavoidable

The Spec adds a refusal that exits 2. Exit codes are public API in this
repository, and its hard rule requires a pull request that changes CLI behavior
to ship the skill update with it.

## What is not governed

Every other path this Spec touches was classified through the governance probe
rather than estimated: `internal/cli`, `internal/spec` and their tests are
ordinary source that no authorization has bounded.

## Approved bounded mutation

State, in the Roundfix skill and its mirror:

- that the reopen command rechecks the gate's staleness immediately before it
  writes, and refuses if the gate changed since preflight.

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
- No change to the skill's version, or to any behavior it documents beyond this
  Spec's own.
- No change to the shared help-token handling, which this Spec excludes.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
