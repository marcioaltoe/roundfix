---
status: approved
granted: 2026-09-30
action: split the Roundfix skill into an entry file plus one reference per command family and the command reference into an index plus one file per command, moving bytes without rewording; make every contract that pinned either document read the entry file with its companions; and teach the write-tasks skill and template that a Task declares the one command file it changes
consuming: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/events.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/profiles.md
  - .agents/skills/roundfix/references/reconcile.md
  - .agents/skills/roundfix/references/release.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/roundfix/references/review-runs.md
  - .agents/skills/roundfix/references/runs.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/references/settle.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/roundfix/references/spec-delivery.md
  - .agents/skills/roundfix/references/stop.md
  - .agents/skills/roundfix/references/storage.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - internal/docscontract/publicdocs_test.go
  - internal/cli/cli_test.go
  - internal/cli/baseline_documentation_contract_test.go
  - skills/baseline_skill_contract_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0194

On 2026-09-30 the maintainer reopened the work on the skills and said:
"considere autorizado a ajustar todas as skills se necessário". In the same
message he stated that architecture, workflow and structures may change
freely. He then approved the program "Três ondas", whose Wave 7 includes
skills that are up to date and smaller. This Spec is the last of that wave.
The four governed test files ride the standing grant of 2026-09-21 for
governed paths a slice genuinely needs. The set was measured with
`GovernedPath` on the authoring branch at `9e439dbb`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the skill
  is the document being split (task_02). `.agents/skills/` is canonical and
  `skills/` its mirror.
- The eighteen files under `.agents/skills/roundfix/references/` — each
  receives the sections of one command family (task_02). Every path under
  `.agents/skills/` is governed. Their mirror copies under
  `skills/roundfix/references/` are ordinary files written by
  `make skills-sync`, so they are not listed.
- `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md` — the authoring
  rule and the Task template must say that a Task declares the one command
  file it changes (task_04).
- `internal/docscontract/publicdocs_test.go`, `internal/cli/cli_test.go`,
  `internal/cli/baseline_documentation_contract_test.go`,
  `skills/baseline_skill_contract_test.go` — each pins text of the skill or of
  the command reference by reading the single file. After the split the text
  lives in a companion file, so each pin must read the entry file with its
  companions, or it fails or silently checks less (task_01).

## What is not governed

The new package `internal/mdtree`, `skills/skills.go`, `skills/skills_test.go`,
the new test files this Spec's Tasks name under `skills/`,
`internal/docscontract/` and `internal/speccheck/`, the mirror copies under
`skills/roundfix/references/`, `skills/write-tasks/references/task-template.md`,
`docs/user-guide/commands.md`, the new files under `docs/user-guide/commands/`
and `docs/user-guide/usage.md` are ordinary. The governed files
`skills/owned_skill_edit_repocontract_test.go`,
`internal/speccheck/governed.go`, `docs/references/coverage-record.json` and
the `Makefile` are not touched, and no setup-owned guide under `docs/agents/`
or Baseline asset is edited.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved skill edits and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No word of the Roundfix skill or of the command reference is changed, added
  or removed, except the two index blocks and the relative link targets of the
  command reference.
- No change to the `Makefile`, to lint, formatter or test-runner
  configuration, to CI workflows, or to `go.mod`.
- No top-level test is renamed or removed, and no archived Spec, existing QA
  Report or Baseline asset is edited.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
