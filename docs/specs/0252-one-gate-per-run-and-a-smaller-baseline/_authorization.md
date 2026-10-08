---
status: approved
granted: 2026-10-08
action: scope the incremental and repository Verification to work outside a Run so a Daemon-assigned Task runs focused tests, record this repository's Verification and runtime decisions with the values CI and the Project Config use, make the runtime header defer to Agent Selection Profiles, state the review, tracker and local-research rules once, make the docs-layout and Secondbrain guides conditional in the root block, and align implement-task, setup-context-driven, the user guide and the glossary
consuming: 0252-one-gate-per-run-and-a-smaller-baseline
paths:
  - .agents/skills/implement-task/SKILL.md
  - .agents/skills/setup-context-driven/SKILL.md
  - docs/agents/agent-instructions.md
  - docs/agents/autonomous-work.md
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - docs/agents/specific-repository.md
  - internal/baseline/assets/decisions.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/issue-tracker.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/templates/guides/autonomous-work.md
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/root/context-workflow.md
  - internal/baseline/assets/templates/root/secondbrain.md
  - internal/baseline/plan_test.go
  - internal/cli/baseline_human_test.go
  - skills/baseline_skill_contract_test.go
  - skills/implement-task/SKILL.md
  - skills/setup-context-driven/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0252

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources, guides and
`.roundfixrc.yml` the maintainer said "Autorizar os dois". The standing answer
for the Governed Paths each Spec declares is "Concedo".

On 2026-10-08 the maintainer asked: "Atualizar e verificar o baseline de regras
para que tenhamos o máximo de performance e resultado com o roundfix e
congruente com as mudanças". The maintainer approved the route "Auditoria e
depois Spec (Recommended)". The audit's Specs are delivered without further
consultation within the bound of four implementation Tasks, with grants for
the governed paths each Spec declares. The audit is committed as
`docs/references/2026-10-08-baseline-audit.md`. This Spec is its Spec A.

Two operator decisions, taken within that grant on 2026-10-08, shape it:

- Finding B03: `runtime.backend` and `runtime.design` are not retired, because
  of the adopter impact. Their values and the guide wording are made current
  with the Project Config (`gpt-6.1-sol`, `opus`, the Light Tier).
- Verification: Task agents must not run the full suite. The Baseline guidance
  is aligned with `implement-task`: the Daemon owns Verification, and agents run
  focused tests only. The recorded repository Verification command matches CI
  (`make verify`). Every claim of time saved stays tied to the audit's measured
  numbers.

No live provider call is authorized by this record. Tests, Verification and QA
use the embedded catalog, temporary repositories and this repository's own
files. No adopter repository is read or written.

The governed set was measured with `GovernedPath` on the authoring branch at
`55de2a73`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository. It ran against every file the Tasks
declare, and 30 of them are governed. The derived files among them were
measured by an authoring prototype in a scratch clone, which applied the four
Tasks in order and ran the sanctioned commands after each.

## Why each governed path is unavoidable

- `internal/baseline/assets/modules/core.json`, `spec-workflow.json`,
  `autonomous-work.json`, `secondbrain.json` and `context-workflow.json` hold
  the clauses, rule, guide and root-block versions this Spec rewords, merges or
  removes. The Module Version Record step writes their version lines.
- `internal/baseline/assets/decisions.json` holds the gate default, the
  incremental suggestion and the two runtime defaults.
- `internal/baseline/assets/templates/guides/autonomous-work.md`,
  `templates/root/context-workflow.md`, `templates/root/secondbrain.md` and
  `templates/index.json` hold the runtime header, the two root sentences and
  their template versions.
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` and
  the six Standard TypeScript formatter goldens are rewritten by
  `make baseline-digests` from the module and template edits. The prototype
  measured exactly these.
- `docs/agents/agent-instructions.md`, `autonomous-work.md`, `secondbrain.md`,
  `spec-routing.md` and `setup-context.json` are rewritten by this
  repository's Managed Refresh. task_02's recorded-decision change also
  rewrites `agent-instructions.md`, `autonomous-work.md` and
  `setup-context.json`.
- `docs/agents/specific-repository.md` loses its duplicate skill-sync rule
  (finding B19).
- `.agents/skills/implement-task/SKILL.md` and
  `.agents/skills/setup-context-driven/SKILL.md`, with their mirrors
  `skills/implement-task/SKILL.md` and `skills/setup-context-driven/SKILL.md`:
  the handoff rule and the suggested values. An owned skill's content changes
  only with its version, and the record command raises both version fields.
- `internal/baseline/plan_test.go`, `internal/cli/baseline_human_test.go` and
  `skills/baseline_skill_contract_test.go` pin the removed status clause, the
  old defaults and the removed ownership tail. Each changes only to the new
  contract.

## What is not governed

These paths are ordinary: `AGENTS.md`, `CONTEXT.md`,
`docs/agents/issue-tracker.md`, `docs/user-guide/context-driven-development.md`,
`docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`,
`docs/references/2026-10-08-baseline-audit.md`,
`internal/baseline/module-versions.json`, `internal/baseline/testdata/**`,
`skills/testdata/owned-skill-versions.json`, the new and changed tests under
`internal/baseline` other than `plan_test.go`, and
`internal/cli/baseline_incremental_verification_test.go` and
`internal/cli/baseline_update_test.go`.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files and the
skill mirrors. This declaration records the regeneration that follows the
approved edits and adds no source path.

```yaml
command: make baseline-digests
```

The Module Version Record step writes each changed module's version line and
the record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/modules/spec-workflow.json
```

```yaml
command: make skills-sync
```

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
renders this repository's guides, root blocks and Setup Manifest, and a second
refresh reports no file change. No digest pin, golden or generated guide is
hand-edited.

## Recorded decision change

`roundfix baseline update` never changes a recorded decision. task_02
therefore uses the automation path. It runs `roundfix baseline plan --profile
go-cli-tui --decision preservation.mode=preservation` with every decision the
Setup Manifest records, changing only `verification.gate` to `make verify`,
`verification.incremental` to `make verify-changed`, `runtime.backend` to
`codex gpt-6.1-sol high` and `runtime.design` to `claude opus high`. It then
runs `roundfix baseline apply` with that plan and its digest. The prototype
showed this writes only the Setup Manifest and the guides that render those
values.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml` or the CI
  workflows.
- No change to the Source Baseline corpus, manifest, accounting or index,
  `internal/baseline/assets/retention/**`, any production Go file, or any
  module other than the five named.
- No change to the `### QA settlement` section of any skill, to archived
  Specs, existing Archive Records or `CHANGELOG.md`.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or reads or writes an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
