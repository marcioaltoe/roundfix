---
status: approved
granted: 2026-09-18
action: keep the shipped Roundfix skill true to the supersession command this Spec delivers
consuming: 0151-a-supersession-the-archive-can-see
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - internal/spec/archive.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0151

On 2026-09-18 the maintainer granted a standing authorization for the Roundfix
skill and its distributed mirror, valid across the remaining slices of this
queue, for one purpose: keeping the shipped skill true to the CLI behavior the
slice itself delivers. This record consumes that authorization for Spec 0151 and
bounds it to this Spec's change.

## Why a governed path is unavoidable

The Spec adds a command and changes what archive accepts as evidence. Both are
public CLI behavior, and this repository's hard rule requires a pull request
that changes CLI behavior to ship the skill update with it.

## The archive decision path

On 2026-09-19 the maintainer granted `internal/spec/archive.go` for this Spec,
bounded to giving `Archive` a second accepted proof.

The first authorization record claimed the whole of `internal/spec` was
unbounded. That was wrong at file granularity: the governance probe was run on
`spec.go`, `errors.go` and `task.go` — the files the work was predicted to touch
— and not on `archive.go`, which a prior authorization had bounded and which
therefore stays governed under ADR-0130. The QA gate caught it as QA-AUTH-PATHS.

The path cannot be avoided. `Archive` is where the evidence decision is made;
moving the decision to `internal/cli/archive.go` would duplicate it in the CLI
while the library function kept refusing.

## What is not governed

Measured through the governance probe, not estimated: `internal/cli/archive.go`,
`internal/spec/spec.go`, `internal/spec/task.go`, their tests and
`docs/user-guide/commands.md` are ordinary source that no authorization has
bounded.

## Approved bounded mutation

State, in the Roundfix skill and its mirror:

- that a Spec whose content another Spec delivered records that as a
  supersession rather than by editing its own status or writing a note;
- that archive accepts a recorded supersession in place of completed Tasks and a
  passing gate, and keeps every other precondition;
- that a supersession naming an unknown Spec, naming itself, or applied to a
  Spec that already carries one is refused.

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
- No change to archive's evidence rules for a Spec that has a Task Graph.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
