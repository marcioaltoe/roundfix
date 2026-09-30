---
status: approved
granted: 2026-09-30
action: correct every statement in the Roundfix-owned skills that contradicts the shipped CLI, Daemon or checker, add the missing command sections to the Roundfix skill, and raise the version of each changed skill
consuming: 0192-owned-skills-that-describe-the-product-as-it-is
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - skills/write-prd/references/prd-template.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - skills/write-techspec/references/techspec-template.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
  - .agents/skills/archive-spec/SKILL.md
  - skills/archive-spec/SKILL.md
  - .agents/skills/setup-context-driven/SKILL.md
  - skills/setup-context-driven/SKILL.md
  - .agents/skills/implement-spec/SKILL.md
  - skills/implement-spec/SKILL.md
  - .agents/skills/council/SKILL.md
  - skills/council/SKILL.md
  - .agents/skills/council/references/archetypes.md
  - .agents/skills/write-idea/SKILL.md
  - skills/write-idea/SKILL.md
  - .agents/skills/write-idea/references/opportunity-scan.md
  - .agents/skills/business-analyst/SKILL.md
  - skills/business-analyst/SKILL.md
  - .agents/skills/brainstorming/SKILL.md
  - skills/brainstorming/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0192

On 2026-09-30 the maintainer reopened the Jev front and approved the program
"Três ondas", whose Wave 7 includes bringing the owned skills up to date. In
the same conversation the maintainer stated: "Sobre a autorização, considere
autorizado a ajustar todas as skills se necessário." That is the express
authorization for the skill paths above. It covers aligning skill text with
shipped behavior; it does not cover restructuring a skill, which Spec 0194
owns.

The set was measured with `GovernedPath` on `9e439dbb`, through a test overlay
that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the skill
  lacks five commands the CLI lists, omits the `agent-selection` event
  category, and describes a delivery order without the pre-PR review
  (task_01).
- The `write-prd`, `write-techspec` and `write-tasks` skill files and their
  templates — they do not teach the authorization record's `operations` and
  `paths: []`, the `## Unreachable Acceptance` section or the four coverage
  units, and `write-tasks` states an old delivery order and denies a checker
  capability that exists (task_03).
- The `qa-gate`, `archive-spec`, `setup-context-driven`, `implement-spec`,
  `council` and `write-idea` skill files and the two reference files — each
  carries a statement the Daemon, a command or the docs layout guide
  contradicts (task_04).
- The `business-analyst` and `brainstorming` skill files, with `write-idea`,
  `write-prd` and `council` above — each shows a question form the
  structured-question clause of the Baseline forbids (task_03, task_04).

`.agents/skills/` is canonical and `skills/` is its mirror.

## What is not governed

The mirrors `skills/write-tasks/references/task-template.md`,
`skills/council/references/archetypes.md` and
`skills/write-idea/references/opportunity-scan.md` are ordinary. So are the new
tests `internal/docscontract/command_documentation_test.go` and
`internal/docscontract/user_guide_contract_test.go`, the guides
`docs/user-guide/commands.md`, `docs/user-guide/usage.md`,
`docs/user-guide/configuration.md`,
`docs/user-guide/context-driven-development.md` and
`docs/user-guide/run-database-lifecycle.md`, and `README.md`.

The governed files `skills/baseline_skill_contract_test.go`,
`skills/settlement_guidance_repocontract_test.go`,
`internal/docscontract/publicdocs_test.go` and `internal/cli/cli_test.go` are
read by this Spec's Verification and are not edited.

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

- No edit to a Baseline module, a generated guide under `docs/agents/`, the
  Makefile, a lint, formatter or test-runner configuration, a CI workflow,
  `go.mod`, `skills/_ownership.yml`, or any governed test.
- No change to the `### QA settlement` section of any skill.
- No change to CLI behavior, help text or exit codes, and no renamed or removed
  top-level test.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
