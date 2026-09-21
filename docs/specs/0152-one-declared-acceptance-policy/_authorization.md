---
status: approved
granted: 2026-09-21
action: apply one declared-acceptance eligibility policy across archive and the derived QA Verification, and keep the shipped skill true to it
consuming: 0152-one-declared-acceptance-policy
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

# Approved authority for Spec 0152

Two grants are consumed here.

The skill files ride the standing authorization of 2026-09-18, valid across the
remaining slices of this queue for one purpose: keeping the shipped skill true
to the CLI behavior the slice itself delivers.

`internal/spec/archive.go` was granted on 2026-09-21, bounded to calling a
single shared eligibility decision instead of carrying its own copy of the rule.

## Why each governed path is unavoidable

Archive is one of the two places that decide whether the newest QA Report is
acceptable, and it holds the policy the repository documented. Leaving it
untouched would mean adding a second authoritative implementation elsewhere and
holding the two in sync with a test — the arrangement that produced the defect
this Spec repairs.

The derived Verification command is public contract surface the skill describes,
and this Spec changes what it accepts.

## What is not governed

Measured by file through the governance probe, not inferred from a directory:
`internal/spec/task.go`, `internal/spec/qa.go`, `internal/spec/task_test.go`,
`internal/spec/qa_test.go` and `internal/cli/settle.go` are ordinary source that
no authorization has bounded.

## Approved bounded mutation

In `internal/spec/archive.go`: replace the inline verdict judgement with a call
to the single shared eligibility decision. Archive's own behavior does not
change — the same reports are accepted and the same reports are refused, for the
same reasons.

In the Roundfix skill and its mirror, state:

- that one eligibility policy decides whether the newest QA Report is
  acceptable, and that settlement and archive both apply it;
- that a `partial` verdict whose blocked rows are declared unreachable settles
  the terminal `qa` Task;
- that `fail`, an undeclared partial, a missing or unparseable report, and a
  `pass` carrying blocked rows each still refuse.

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
- No change to which QA Report is newest, or to what makes a blocked row
  unreachable.
- No change to the skill's version, or to any behavior it documents beyond this
  Spec's own.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
