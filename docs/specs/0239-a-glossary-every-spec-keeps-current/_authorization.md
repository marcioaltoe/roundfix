---
status: approved
granted: 2026-10-06
action: make the Spec Consistency Check and the Archive Command hold a Spec's Glossary Declaration to the glossary, state in the Baseline's domain clauses that the glossary and the ADRs are the base and that a Spec writes the terms it introduces, make the authoring, QA and Roundfix skills follow it, and add the terms dropped since 2026-10-02 to the glossary
consuming: 0239-a-glossary-every-spec-keeps-current
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/glossary.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - internal/speccheck/coherence.go
  - internal/speccheck/constraints.go
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0239

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário", and of the Baseline source, the guides and
`.roundfixrc.yml` "Autorizar os dois". On 2026-10-06 the maintainer stated the
intent this Spec restores, "quando criei o roundfix a ideia era ter o
context.md e adrs como base do processo de context-driven", and answered the
cycle's questions: the scope "Tudo, na ordem sugerida", of which this Spec is
first and released as v0.50.0; the catch-up of the terms dropped since
2026-10-02 done by editing `CONTEXT.md` directly through `domain-modeling`,
"Sim"; and, for the Governed Paths this Spec declares (the glossary, the
Roundfix and authoring skills, the guides, and the Baseline modules and
clauses), "Concedo". The maintainer also authorized a read-only measurement of
the adopters' `CONTEXT.md` and git log, "Sim, só leitura", for conexus,
fluxus, vortex, tax-poc and oraculum: nothing is written there and none of
their content is sent anywhere.

The governed set was measured with `GovernedPath` on the authoring branch at
`40a7893d`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. No live provider call is
authorized by this record: tests, Verification and QA use temporary
repositories and fake runners.

## Why each governed path is unavoidable

- `internal/speccheck/constraints.go` and `internal/speccheck/coherence.go` —
  the full Spec Consistency Check and its staged form call the new glossary
  detector and list its codes by stage.
- `internal/docscontract/testdata/corpus-golden.json` and
  `internal/spec/archive_layout_characterization_test.go` — the active-corpus
  golden and the characterization that pins it gain the three new codes at
  zero, as every new detector code has.
- `internal/baseline/assets/modules/context-workflow.json` — the two domain
  clauses are reworded in place and the module, rule and guide versions rise.
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md`
  and `docs/agents/setup-context.json` — derived from the module edit by the
  sanctioned regeneration and the Managed Refresh below; the authoring probe
  measured exactly these.
- The canonical files under `.agents/skills/` for `write-prd`,
  `write-techspec`, `write-tasks`, `qa-gate` and `roundfix`, with their
  governed mirrors under `skills/` — the authoring route must call
  `domain-modeling` and write the declaration, the gate must record it, and the
  skill-sync rule requires the Roundfix Skill to describe the new check and
  refusal; an owned skill's content changes only with its version.

## What is not governed

`CONTEXT.md`, `docs/agents/domain.md`, the Go sources and tests under
`internal/speccheck`, `internal/cli`, `internal/docscontract` and
`internal/baseline` other than those listed,
`internal/baseline/testdata/**`, `docs/user-guide/commands/spec.md`,
`docs/user-guide/commands/archive.md`, the skill reference mirrors
`skills/roundfix/references/spec.md`, `skills/roundfix/references/archive.md`
and `skills/write-prd/references/glossary.md`,
`skills/testdata/owned-skill-versions.json` and
`docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`
are ordinary.

## Sanctioned regeneration

The repository-owned commands resolve the derived Baseline files and the skill
mirrors. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: make baseline-digests
```

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to the Source Baseline corpus, to any clause identity or
  enforcement, or to the `domain-modeling`, `grilling` or `grill-with-docs`
  skills, which are upstream content.
- No Jev judge question, call or key read; the judge stays advisory and no
  test or Verification reaches a provider.
- No change to archived Specs, existing QA Reports, `CHANGELOG.md` or the
  `### QA settlement` section of any skill; no write to an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
