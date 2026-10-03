---
status: approved
granted: 2026-09-30
action: refuse a non-Markdown file that names an active Spec's directory and resolve another Spec's Task Context through the archive, resume a Delivery Retry at review after a post-archive correction, continue an existing item branch on deliver start, add the QA Archive Override Delivery Convention, report a Verification command the shell cannot parse as malformed and name uncommitted Verification sources, and describe this in the Roundfix Skill and the command guides
consuming: 0219-a-delivery-that-survives-archive-requeue-and-review
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0219

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, under the broad autonomy granted on 2026-09-29 for authoring,
corrective work and delivery through merge, and said of the skills
"considere autorizado a ajustar todas as skills se necessário". On 2026-10-01
the maintainer extended the program ("Tudo, de A a F"), and on 2026-10-03
answered "Continue" to the proposed next cycle G, H and I, of which this Spec
is G: a delivery that survives archive, requeue and review.

The governed set was measured with `GovernedPath` on the authoring branch at
`305211a3`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. The probe also reported
`internal/speccheck/constraints.go`, `internal/speccheck/coherence.go` and
`internal/spec/archive.go` governed; the design keeps all three untouched: the
new Spec Consistency findings run from the reference and Verification
detectors, and the archive refusal lives in the Archive Command.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, and the
  repository's skill-sync rule requires a Pull Request that changes CLI
  behavior to ship the skill update.
- `.agents/skills/roundfix/references/archive.md` — the archive refusal.
- `.agents/skills/roundfix/references/spec.md` — the two new Spec Consistency
  findings, the `malformed` verdict and the uncommitted-source line.
- `.agents/skills/roundfix/references/deliver.md` — the retry after a
  post-archive correction and the continued item branch.
- `.agents/skills/roundfix/references/review.md` — the fifth Delivery
  Convention.

## What is not governed

The Go sources and tests under `internal/cli`, `internal/daemon`,
`internal/delivery` and `internal/speccheck` that the Tasks declare,
`skills/roundfix/references/archive.md`, `skills/roundfix/references/spec.md`,
`skills/roundfix/references/deliver.md`, `skills/roundfix/references/review.md`,
`skills/testdata/owned-skill-versions.json`, the command guides under
`docs/user-guide/commands/` and ADR-0223 and ADR-0226 are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edits and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows, `.roundfixrc.yml`,
  `internal/speccheck/constraints.go`, `internal/speccheck/coherence.go` or
  `internal/spec/archive.go`.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  the Baseline sources, `CONTEXT.md` or the `### QA settlement` section of any
  skill.
- No test, Verification command or QA row reaches a provider, GitHub or the
  network, or writes under the real `~/.roundfix`; every queue, retry and
  archive scenario runs in a disposable repository and home.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
