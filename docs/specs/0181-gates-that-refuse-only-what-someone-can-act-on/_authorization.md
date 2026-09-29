---
status: approved
granted: 2026-09-29
action: record the paths a Task changed without declaring them in a Daemon-owned Task file section, stop the pre-PR Pull Request row from deciding a qualifying partial, give SC-ADR-RELATED a commit-ancestry horizon, and align the qa-gate, write-tasks and Roundfix skills with that behavior
consuming: 0181-gates-that-refuse-only-what-someone-can-act-on
paths:
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0181

The maintainer approved Wave 4 on 2026-09-29, and this Spec is part of it. The
same day the maintainer expressly authorized edits to the `qa-gate` and
`write-tasks` skills, answering "Autorizar" to a structured question naming the
four exact paths and the sanctioned regeneration. The grant covers only
aligning those skills with the behavior this Spec ships. The Roundfix skill
files ride the standing grant of 2026-09-18 for keeping the shipped skills true
to the CLI. The set was measured with `GovernedPath` on `a626494c`.

## Why each governed path is unavoidable

- `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md` — the skill
  defines the qualifying declared `partial` and tells the QA Agent to record
  every Pull Request journey as environment-blocked. After task_02 its
  definition would be false. Its scope rule must also count a path in a Task's
  `## Recorded paths` section as declared (task_04). `.agents/skills/` is
  canonical and `skills/` its mirror.
- `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md` — the
  skill states that every path a Task edits is declared, and it must say that
  the Daemon records an undeclared path at commit (task_04).
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill states the eligibility policy that `roundfix archive`, `roundfix
  qa-report accept` and QA settlement apply, including that an
  environment-blocked row refuses a `partial` (task_04).

## What is not governed

`internal/daemon/task_engine.go`, `internal/spec/qa.go`,
`internal/agent/spec_prompt.go`, `internal/speccheck/citations.go`, the new
files `internal/spec/recorded_paths.go` and `internal/speccheck/adr_horizon.go`,
and the test files this Spec's Tasks name under `internal/spec/`,
`internal/daemon/`, `internal/agent/` and `internal/speccheck/` are ordinary.
So are `docs/user-guide/commands.md`,
`docs/user-guide/context-driven-development.md` and `CONTEXT.md`. The governed
files `internal/spec/archive.go`, `internal/spec/archive_test.go`,
`internal/speccheck/constraints_characterization_test.go`,
`internal/speccheck/mechanical_test.go`,
`internal/docscontract/testdata/corpus-golden.json`,
`skills/baseline_skill_contract_test.go` and
`docs/references/coverage-record.json` are not touched. No top-level test is
renamed or removed, and the write-tasks phrases that
`TestWriteTasksSkillStatesTheDeclaredPathRules` requires stay in the skill.

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

- No edit to a Task's authored `## Context` by the Daemon, to archived Specs, to
  existing QA Reports, to `skills/_ownership.yml` or to the CI workflow, and no
  renamed or removed top-level test.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate, and the Secondbrain is read only.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
