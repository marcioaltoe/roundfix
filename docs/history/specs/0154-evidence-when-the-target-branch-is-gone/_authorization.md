---
status: approved
granted: 2026-09-21
action: prove delivered content against the default branch when a Run's target branch is gone, and keep the shipped skill true to what reconcile reports
consuming: 0154-evidence-when-the-target-branch-is-gone
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

# Approved authority for Spec 0154

The skill files ride the standing authorization of 2026-09-18, valid across the
remaining slices of this queue for one purpose: keeping the shipped skill true
to the CLI behavior the slice itself delivers.

The standing authorization of 2026-09-21 covers any governed path a slice
genuinely needs, on the condition that the bounded set is measured against
`internal/speccheck/governed.go` rather than predicted, and recorded here with
its reason. Measured for this Spec: the intersection is the two skill files and
nothing else. `internal/worktree/worktree.go` and `internal/worktree/worktree_test.go`
are ordinary source that no authorization has bounded.

## Why a governed path is unavoidable

`reconcile` reports which Runs are releasable and why, and this Spec changes
what it can prove when the target branch is gone. The shipped skill documents
that reporting, and the repository's hard rule requires a pull request that
changes CLI behavior to ship the skill update with it.

## Approved bounded mutation

State, in the Roundfix skill and its mirror:

- that a Run whose target branch is absent is checked for content evidence
  against the default branch, using the rule already accepted for a missed
  ancestry;
- that a Run is released only on positive content evidence, and that an
  unreachable default branch, an unarchived Spec, or evidence naming another
  Spec each leave it preserved;
- that the preserved reason names the proof that was missing.

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
- No weakening of the deleted-target content proof, which Spec 0125 Core Feature
  3 requires this Spec to retain.
- No release of a Run whose Worktree holds uncommitted changes.
- No change to the branch-naming policy or to repository identity.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
