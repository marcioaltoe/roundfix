---
status: approved
granted: 2026-10-06
action: retire the-fool and council from the Baseline modules, triggers and Setup Snapshots, keep them out of the asset sync, remove council from the Roundfix-owned bundle and the owned skills that name it, drop the-fool and autoresearch from this repository's lock, recommended list and installed trees, and make baseline update list the retired copies an adopter still holds
consuming: 0241-retire-the-fool-autoresearch-and-council
paths:
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - docs/agents/skill-dispatch.md
  - docs/agents/setup-context.json
  - Makefile
  - .agents/skills/council/SKILL.md
  - .agents/skills/council/assets/synthesis-template.md
  - .agents/skills/council/references/archetypes.md
  - .agents/skills/council/references/debate-protocols.md
  - skills/council/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-idea/references/idea-template.md
  - skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
  - .agents/skills/the-fool/SKILL.md
  - .agents/skills/the-fool/references/cognitive-bias-inventory.md
  - .agents/skills/the-fool/references/dialectic-synthesis.md
  - .agents/skills/the-fool/references/evidence-audit.md
  - .agents/skills/the-fool/references/mode-selection-guide.md
  - .agents/skills/the-fool/references/pre-mortem-analysis.md
  - .agents/skills/the-fool/references/red-team-adversarial.md
  - .agents/skills/the-fool/references/socratic-questioning.md
  - .agents/skills/autoresearch/SKILL.md
  - .agents/skills/autoresearch/references/eval-guide.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0241

On 2026-09-30 the maintainer expressly authorized editing the Baseline source
and its generated guides ("Autorizar os dois") and every skill change
("considere autorizado a ajustar todas as skills se necessário"). On
2026-10-06, asked through AskUserQuestion which skills to retire, the
maintainer answered: "remover: the-fool, autoresearch, council // manter:
grilling, grill-with-docs, write-idea, business-analyst e handoff". Asked
about scope, the maintainer answered "Tudo, na ordem sugerida"; this Spec is
the third item, planned for v0.52.0. Asked whether the adopters' release note
should explain how to remove the local copies, the maintainer answered "Sim,
com instrução". For the Governed Paths this Spec declares, the maintainer
answered "Concedo".

The governed set was measured with `GovernedPath` on the authoring branch at
`40a7893d`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `internal/baseline/assets/modules/context-workflow.json` — the module
  requires and dispatches both retired skills.
- The four `internal/baseline/assets/setups/*.json` snapshots — they list both
  skills, and the asset sync is their only writer.
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, the
  formatter golden `skill-dispatch.md`, `docs/agents/skill-dispatch.md` and
  `docs/agents/setup-context.json` — derived from the catalog and the lock;
  changed only by the sanctioned regeneration and the managed refresh.
- `Makefile` — its `OWNED_SKILLS` list names `council`; `make skills-sync`
  copies every listed tree, so it fails once the tree is deleted. Only that
  line and its comment's count change.
- The `council`, `the-fool` and `autoresearch` trees — deleted, never edited.
- `write-idea` and `write-prd`, both copies — they send the reader to
  `council`; an owned skill's text changes only with its version.
- `.agents/skills/roundfix/SKILL.md`, its mirror and
  `.agents/skills/roundfix/references/baseline.md` — the repository's
  skill-sync rule requires a Pull Request that changes command behavior to
  update the skill that describes it, and `baseline update` gains the
  `Skills retired` report.

## What is not governed

The Go sources and tests under `internal/baseline`, `internal/cli` and
`skills`, `skills/skills.go`, `skills/recommended.txt`, `skills-lock.json`,
`skills/testdata/owned-skill-versions.json`, the parity fixture and the other
derived files under `internal/baseline/testdata/`, the remaining skill mirrors,
`docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`,
the user guides, `README.md` and `CONTEXT.md` are ordinary.

## Sanctioned regeneration

The repository-owned commands resolve the derived Baseline files and the skill
mirrors. This declaration records the regeneration that follows the approved
edit and adds no source path.

```yaml
command: make baseline-digests
```

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, no change to the lint, formatter or
  test-runner configuration or the CI workflows, and no `Makefile` change
  beyond the `OWNED_SKILLS` line and its comment.
- No edit to the content of a vendored skill; `the-fool` and `autoresearch`
  are only deleted. No write to `~/dev/skills` or the upstream repository.
- No change to a Normative Clause, a Source Baseline, a retention transition
  or a parity fixture other than `asset-sync.json`.
- No test, Verification command or QA row reaches the network, reads or
  writes the real `~/.roundfix`, or touches an adopter's repository.
- No change to archived Specs, existing QA Reports, `CHANGELOG.md` or the
  `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
