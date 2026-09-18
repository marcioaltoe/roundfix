---
status: approved
granted: 2026-09-18
action: describe the pre-Pull-Request review policy check in the Roundfix skill that ships with the CLI
consuming: 0143-a-repository-says-who-reviews-before-the-pull-request
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

# Approved authority for Spec 0143

The Spec adds a readiness check to the Doctor Command, which is public CLI
behavior. The repository's hard rule requires a pull request that changes CLI
behavior to ship the Roundfix skill update with it. Asked to approve this record
with its exact paths and operations, the maintainer authorized both copies on
2026-09-18, after pre-PR review named the stale skill.

## Why a governed path is unavoidable

The Roundfix skill is how an operating Agent learns what the CLI does. It is a
Roundfix-owned Skill, so changing its text needs an express grant. Every Run
reads the canonical copy under `.agents/skills/`, and the binary embeds the
distributed mirror under `skills/`, so both are bounded here.

The configuration schema, the Doctor check itself and the user guide are
ordinary source that no authorization has bounded, so they need no grant.

## Approved bounded mutation

Describe, in the Roundfix skill and its mirror:

- that the Doctor Command reports the resolved pre-Pull-Request review provider
  and the configuration layer that supplied it;
- that `none` is reported as review disabled by configuration;
- that the check invokes no provider and mutates nothing.

After the canonical edit, regenerate the distributed mirror with
`make skills-sync`. The mirror is already a bounded path. If any derived pin
changes as a result, rewrite it only with `make baseline-digests`.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No action, operation or path beyond those above.
- No change to the skill's version, its owned-skill contract, or any behavior it
  documents beyond the new check.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No edit to this repository's own `.roundfixrc.yml`.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
