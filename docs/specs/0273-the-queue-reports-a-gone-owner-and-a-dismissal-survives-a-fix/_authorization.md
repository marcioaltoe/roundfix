---
status: approved
granted: 2026-10-16
action: make roundfix deliver status name a recorded Queue Owner proven not running, make the deliver status, resume, retry and stop refusals run in an item worktree name the checkout that owns the Delivery Queue, and make roundfix review dispose --dismiss accept a dismissal after a fix that leaves the finding's anchor unchanged, with the deliver and review user guides, the Roundfix skill's deliver and review references and the glossary describing them
consuming: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
paths:
  - .agents/skills/roundfix/SKILL.md
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

# Approved authority for Spec 0273

On 2026-10-16 the maintainer approved cycle AH through AskUserQuestion: scope
"A: history só com records, B: ergonomia da fila, C: guardas de autoria, D:
remetente do F6", order "B → A → C → D (Recomendado)", authorizations
"Governed Paths declarados" (the option text: "\"Concedo\" como em AE/AF/AG
para os Governed Paths que cada Spec declarar: módulos da Baseline, registros
de versão, make baseline-digests, skills e testes governados. Makefile,
workflows de CI, go.mod e configuração de lint ficam fora."), "Merge e
releases", "Merge da PR #545", "Avisos por push"; failures "Regras de AE–AG
(Recomendado)". For this Spec, item B, the decided scope is: "deliver status
avisa quando o dono da fila morreu e indica deliver resume; review dispose
aceita dismiss depois de um fix commitado (head revisado é ancestral, âncora
intacta); a recusa do retry nomeia o checkout principal." On 2026-09-30 the
maintainer had granted the skills ("considere autorizado a ajustar todas as
skills se necessário").

The governed set was measured on the authoring branch at `9e723b82` with
`GovernedPath`, through a `go test -overlay` probe in `internal/speccheck`
that wrote nothing to the repository, against every file the Tasks declare
(`instruction:`, `interface:` and `creates:`) and every planning file. Four
declared edits are governed; `.agents/skills/domain-modeling/SKILL.md` is
governed too but only read, as an `instruction:` of task_04.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` covers the `command deliver
  status`, `resume`, `retry` and `stop` surfaces in the Skill Coverage Map,
  whose source `internal/cli/deliver.go` task_01 and task_02 change; task_04
  describes the `Owner:` line and the item worktree refusal there.
- `.agents/skills/roundfix/references/review.md` covers the `command review
  dispose` surface, whose source `internal/cli/review.go` task_03 changes;
  task_04 replaces its old dismissal rule.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  skill's version fields, raised with its content.

## What is not governed

The Go sources and tests under `internal/cli` and `internal/store`; the
mirrors `skills/roundfix/references/deliver.md` and
`skills/roundfix/references/review.md`;
`skills/testdata/owned-skill-versions.json`;
`docs/references/behavior-surfaces.json`; `CONTEXT.md`;
`docs/user-guide/commands/deliver.md`; `docs/user-guide/commands/review.md`;
and the ADR
`docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill version and
the Behavior Surface Record. This declaration records the regeneration that
follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
records the raised skill version, and
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
re-records the fingerprints of the `deliver` and `review` user guides and of
the skill. No Baseline module changes, so no module version is recorded.
`make baseline-digests` was not run while authoring; Spec 0271 measured on
2026-10-13 that it rewrites nothing after a Roundfix skill reference edit, so
it is not declared, and the QA scope audit confirms that no derived Baseline
file changed. No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module, the review
  conventions, the Run Database schema, or the `### QA settlement` section of
  any skill.
- No change to the `deliver` or `review dispose` help text, to the review
  correction proof or to the Reviewer Lineage.
- No test, Verification command or QA row reads or writes the real
  `~/.roundfix`, reaches the network, reads another repository, or signals a
  process it did not start.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
