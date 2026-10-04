---
status: approved
granted: 2026-09-30
action: keep the verdict of a read-only review turn that ended after the session refused a permission, give its findings ids, name a runtime's "Prompt is too long" answer in the review reason, bound the Claude review prompt by estimated tokens against half of its context window, and describe this in the Roundfix Skill and the review command guide
consuming: 0222-a-review-through-claude-that-keeps-its-verdict
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

# Approved authority for Spec 0222

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-04, answering a structured question, the
maintainer authorized delivering this Spec and up to about five real reviews
through `claude / opus` to reproduce the defect and prove the fix, with the
rest of the work through Codex.

The governed set was measured with `GovernedPath` on the authoring branch at
`6ea9e1e9`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/review.md` — the review command's
  reference, which states the classification after a refused permission, the
  prompt-too-long reason and the Claude prompt bound.

## What is not governed

The Go sources and tests under `internal/agent` and `internal/cli` that the
Tasks declare, `skills/roundfix/references/review.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/review.md` and ADR-0227 are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirror. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows or `.roundfixrc.yml`.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  `CONTEXT.md` or the `### QA settlement` section of any skill.
- No test or Verification command reaches a provider, GitHub or the network,
  or writes under the real `~/.roundfix`; tests use fake runners, a fake acpx
  and disposable homes. The QA gate may run at most one real review through
  `claude / opus` in a disposable clone, inside the maintainer's budget of
  about five; three were spent during authoring.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
