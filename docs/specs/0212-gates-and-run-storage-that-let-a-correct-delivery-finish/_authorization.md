---
status: approved
granted: 2026-09-30
action: omit QA evidence and upstream-managed skill copies from the pre-PR review's diff with a measured bound and a fuller record, skip compiled source in the failed-pass QA import, add a deletes Task Context kind checked at settlement, release a merged and archived Spec's Runs whose only unrepresented work is Spec-directory or declared, and describe this in the Roundfix, qa-gate and write-tasks skills and the command guides
consuming: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/roundfix/references/reconcile.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - skills/roundfix/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0212

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, under the broad autonomy granted on 2026-09-29 for authoring,
corrective work and delivery through merge, and said of the skills
"considere autorizado a ajustar todas as skills se necessário". On 2026-10-01
the maintainer extended the program to the next cycle, "Tudo, de A a F", one
minor release per wave, with this Spec as part of wave A: a delivery queue
that finishes without operator intervention.

The governed set was measured with `GovernedPath` on the authoring branch at
`30cd147c`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. The probe also found
`internal/speccheck/coherence.go` governed; the design keeps it untouched by
checking a deleted path at settlement instead of adding a Spec Consistency
Check finding.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, and the
  repository's skill-sync rule requires a Pull Request that changes CLI
  behavior to ship the skill update.
- `.agents/skills/roundfix/references/review.md` — the omitted paths, the
  bound and the record fields.
- `.agents/skills/roundfix/references/reconcile.md` — the merged Spec's
  release rule.
- `.agents/skills/qa-gate/SKILL.md` and `skills/qa-gate/SKILL.md` — the QA
  Agent keeps runnable evidence out of compiled source, under a new heading;
  the `### QA settlement` section stays byte-identical.
- `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md` and
  `skills/write-tasks/SKILL.md` — the `deletes` Context kind.

## What is not governed

The Go sources and tests under `internal/cli`, `internal/daemon`,
`internal/spec`, `internal/speccheck` and `internal/worktree` that the Tasks
declare, `skills/roundfix/references/review.md`,
`skills/roundfix/references/reconcile.md`,
`skills/write-tasks/references/task-template.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/review.md`, `docs/user-guide/commands/reconcile.md`
and ADR-0212 are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edits and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows, `.roundfixrc.yml`
  or `internal/speccheck/coherence.go`.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  `CONTEXT.md` or the `### QA settlement` section of any skill.
- No test, Verification command or QA row reaches a provider, GitHub or the
  network, or writes under the real `~/.roundfix`; no QA row runs
  `roundfix reconcile --apply` outside a disposable repository.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
