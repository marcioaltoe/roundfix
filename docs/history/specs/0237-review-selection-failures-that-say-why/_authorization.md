---
status: approved
granted: 2026-10-06
action: let a blocked pre-PR review name the protocol step that failed and the adapter's message, retry a failure before the prompt once on the same selection, and describe it in the review guide and the Roundfix Skill's review reference
consuming: 0237-review-selection-failures-that-say-why
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/review.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0237

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-06 the maintainer answered the cycle's
questions through AskUserQuestion: the scope "Y1, Y2 e Y3", of which this Spec
is Y3, delivered third; and, for the Governed Paths this Spec declares,
"Concedo". No live provider call is authorized by this record: authoring,
tests, Verification and QA use fake runners, a fake acpx and a fake ACP
adapter under a temporary home.

The governed set was measured with `GovernedPath` on the authoring branch at
`bde616d1`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/review.md` — the repository's skill-sync
  rule requires a Pull Request that changes CLI behavior to ship the skill
  update, and this reference describes the review's reason, record and
  fallback.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields by the record command.

## What is not governed

The Go sources and tests under `internal/agent` and `internal/cli`,
`docs/adr/0242-a-blocked-review-names-the-protocol-step-and-retries-once-before-the-prompt.md`,
`docs/user-guide/commands/review.md`, the skill reference mirror
`skills/roundfix/references/review.md`, and
`skills/testdata/owned-skill-versions.json` are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to the daemon's selection fallback, its reason codes, the review
  provider policy or the review profiles.
- No test, Verification command or QA row reaches a provider, reads a
  credential, or reads or writes the real `~/.roundfix` or `~/.acpx`.
- No change to archived Specs, existing QA Reports, `CONTEXT.md`,
  `CHANGELOG.md` or the `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
