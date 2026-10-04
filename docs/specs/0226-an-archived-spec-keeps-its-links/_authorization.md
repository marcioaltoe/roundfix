---
status: approved
granted: 2026-10-04
action: let roundfix archive rewrite the relative Markdown links that leave a Spec and refuse when one reaches nothing, let the Delivery Queue accept those rewrites in a resumed archive commit, and describe this in the Roundfix Skill and the archive guide
consuming: 0226-an-archived-spec-keeps-its-links
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - internal/spec/archive.go
  - internal/spec/archive_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0226

On 2026-10-04 the maintainer, asked through a structured question, decided to
deliver this cycle's item L, the adjustments the adopters asked for, and
granted the Governed Path `internal/spec/archive.go` by name, answering
"Concedo", to fix the archive's relative links. On 2026-09-30 the maintainer
had said of the skills "considere autorizado a ajustar todas as skills se
necessário". Item L was split into Spec 0223 and this Spec, which owns the
archive's links.

The governed set was measured with `GovernedPath` on the authoring branch at
`6ea9e1e9`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `internal/spec/archive.go` — the archive's move lives there; the link
  planning, refusal, rewrite and the match the Delivery Queue reuses go beside
  it (task_02).
- `internal/spec/archive_test.go` — the replay test
  `TestSpec0058ReplayArchivesDeclaredUnreachableRelease` copies archived Spec
  0058 into an empty temporary repository, where the Spec's three outward
  links reach nothing, so the new refusal fails it; `prepareSpec0058Replay`
  creates those three targets and no assertion changes (task_02).
- `.agents/skills/roundfix/SKILL.md`, its mirror `skills/roundfix/SKILL.md`
  and `.agents/skills/roundfix/references/archive.md` — the rewrite, the
  refusal and the count the command prints, and the version raise the
  skill-sync rule requires for a change of CLI behavior (task_01).

## Grant for the replay test

On 2026-10-04 the maintainer, asked through a structured question for a named
grant of the governed test `internal/spec/archive_test.go` so the replay of
Spec 0058 creates the targets of its three outward links (a nine-line fixture
change, assertions unchanged), answered "Concedo". The grant covers that
fixture change in `prepareSpec0058Replay` and nothing else in the file.

## What is not governed

`internal/cli/archive.go`, `internal/cli/deliver_workflow.go`, the three new
test files, `skills/roundfix/references/archive.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/archive.md`, ADR-0230 and this Spec's own files are
ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows or `.roundfixrc.yml`.
- No change to archived Specs, existing QA Reports, the archive stamp,
  archive eligibility, Park Class names, `CONTEXT.md` or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row reaches GitHub, a provider or the
  network, or writes under the real `~/.roundfix`.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
