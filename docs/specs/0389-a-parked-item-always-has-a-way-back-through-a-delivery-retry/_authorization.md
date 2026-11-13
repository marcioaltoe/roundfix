---
status: approved
granted: 2026-11-13
action: requeue through a Delivery Retry an item parked before it has an Item Branch, read the delivery authorization of an archived Spec without a Spec Root directory, return an item whose Pull Request holds an earlier candidate to the repository gate, name the review-stale next action, and document them in the deliver guide, the Roundfix skill's deliver reference, the glossary and ADR-0398
consuming: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0389

On 2026-11-13 fiscus reported, in ai-memory message
`01a1d42a-28f1-72d3-9957-76f106877d2e`, that `roundfix deliver retry` refuses
an item parked before it has an Item Branch, and on 2026-11-12 quaestor
reported, in ai-memory message `01a1d231-03da-79f1-9834-51487a593986`, that a
retry after a Pull Request conflict resolution found no Spec Root and no
pushed head. The operator recorded both as Backlog Entries, which this Spec
adopts
(`references/2026-11-13-a-retry-refuses-an-item-parked-before-it-has-a-branch.md`,
`references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md`).
The maintainer then answered the operator's cycle BE questions with
AskUserQuestion (2026-11-13 ~08:20 BRT). The answers this Spec rests on are,
verbatim:

- Question "Quais trabalhos do backlog do Roundfix eu planejo e entrego agora,
  sem supervisão? (Cada item marcado vira Spec com plano validado e é entregue
  pela fila. Os desmarcados ficam no backlog.)" → selected "Retry da fila
  (Recommended)" — option text: "Uma Spec com os dois defeitos do retry: item
  parado antes de ter branch não aceita retry (Fiscus), e o retry depois de
  resolver conflito no PR não acha a pasta de Specs nem envia o head
  (Quaestor)." This Spec is that item, delivered after Spec 0388.
- Question "Quando publico as releases desses trabalhos?" → "Patch após cada
  Spec (Recommended)" — option text: "v0.143.1 assim que sair a correção do
  aviso; a próxima patch assim que sair o retry. Os adotantes que já trouxeram
  a 0.143.0 (a Argus, por exemplo) recebem logo a correção do texto."
- Question "Quanta folga eu tenho para consertar falhas sem parar? (Na 0386 o
  review usou 3 correções.)" → "Até 3 por Spec + Makefile/CI (Recommended)" —
  option text: "Até 3 correções por Spec, contando as rodadas do review. Fica
  liberado mexer em Makefile, workflows do CI, go.mod e lint quando a correção
  exigir; cada mudança vai descrita no relatório."

The cycle BE rules keep the earlier cycles' grants of owned skills and of
changes to skills that reach the adopters, the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos
que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário").

This Spec's design needs no Makefile, `.github/workflows`, `go.mod`, `go.sum`,
lint or formatter configuration change, and no Baseline module, template or
profile change. The corrections answer would grant a Makefile, workflow,
`go.mod` or lint change if a correction required one; none is listed here,
and any such change must first be added to `paths` on the default branch with
its reason.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-13 on the branch `docs/plan-0389`,
based on `f953f77f`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files, its adopted references and ADR-0398 afterwards, and the
  probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 to 04 changed or
  created, every planning file of this Spec, its adopted references and their
  former paths, ADR-0398 and `docs/agents/setup-context.json`.
- The three paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver` reference, which task_04 changes. The probe also
  named `.agents/skills/domain-modeling/SKILL.md`, which task_01 only reads,
  and `docs/agents/setup-context.json`, which no Task changes; neither is
  granted.

## What is not governed

- Every Go source and test file tasks 01 to 03 create or change under
  `internal/delivery/`, `internal/store/`, `internal/spec/` and
  `internal/cli/`.
- The mirror `skills/roundfix/references/deliver.md` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/deliver.md`, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json` and ADR-0398.
- Every planning file of this Spec and its adopted references.

## Sanctioned regeneration

The repository-owned commands write the skill mirror and the generated
records. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline
  module, template, profile, Setup Snapshot or plan golden, any guide under
  `docs/agents/`, `docs/agents/setup-context.json`,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `deliver`, or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
