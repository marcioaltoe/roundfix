---
status: approved
granted: 2026-09-18
action: keep the shipped Roundfix skill true to the release planner's accepted tag spellings and its ambiguity refusal
consuming: 0147-a-planner-that-reads-both-tag-spellings
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

# Approved authority for Spec 0147

On 2026-09-18 the maintainer granted a standing authorization for the Roundfix
skill and its distributed mirror, valid across the remaining slices of this
queue, for one purpose: keeping the shipped skill true to the CLI behavior the
slice itself delivers. This record consumes that authorization for Spec 0147 and
bounds it to this Spec's change.

## Why a governed path is unavoidable

The Spec teaches the release planner a second accepted tag spelling and adds a
typed preflight refusal when the highest version exists under both. Both are
public CLI behavior, and the repository's hard rule requires a pull request that
changes CLI behavior to ship the skill update with it.

The planner, the command surface, its tests and the user guide are ordinary
source that no authorization has bounded.

## Approved bounded mutation

State, in the Roundfix skill and its mirror:

- that the Release Plan Command accepts a stable tag written with or without the
  `v` prefix;
- that when the highest reachable version exists under both spellings, the
  command refuses in preflight and names the selector that resolves it;
- that the proposed version keeps the spelling of the tag it was selected from.

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
- No change to this repository's own release workflow, its tag trigger, or the
  `v`-prefixed tags it publishes.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
