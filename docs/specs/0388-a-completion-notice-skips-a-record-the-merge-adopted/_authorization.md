---
status: approved
granted: 2026-11-13
action: rewrite the closure rule of the Completion notices section of the roundfix skill's release reference so that a Backlog Entry or Finding the merge adopted sends no completion notice, restore the guarantees the section lost, pin each guarantee with a contract-test phrase and an executable check of the section's Git commands, and raise and record the roundfix skill's version with the repository's record command
consuming: 0388-a-completion-notice-skips-a-record-the-merge-adopted
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/release.md
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

# Approved authority for Spec 0388

The Spec adopts the Backlog Entry
`docs/backlog/2026-11-13-a-completion-notice-reads-an-adopted-record-as-closed.md`,
which records ai-memory message `01a1d483` of 2026-11-13; reported by argus.

The maintainer answered the operator's questions through AskUserQuestion on
2026-11-13 at about 08:20 BRT:

- Scope. Question: "Quais trabalhos do backlog do Roundfix eu planejo e
  entrego agora, sem supervisão? (Cada item marcado vira Spec com plano
  validado e é entregue pela fila. Os desmarcados ficam no backlog.)" —
  selected "Aviso de conclusão (Recommended)", option text: "Spec corretiva
  da regra publicada na v0.143.0: a regra lê como fechado um registro que uma
  Spec adotou e pode mandar o aviso cedo e de novo (achado da Argus; a falha
  foi minha na rodada 2 da 0386)."
- Releases. Question: "Quando publico as releases desses trabalhos?" —
  answer: "Patch após cada Spec (Recommended)", option text: "v0.143.1 assim
  que sair a correção do aviso; a próxima patch assim que sair o retry. Os
  adotantes que já trouxeram a 0.143.0 (a Argus, por exemplo) recebem logo a
  correção do texto."
- Corrections. Question: "Quanta folga eu tenho para consertar falhas sem
  parar? (Na 0386 o review usou 3 correções.)" — answer: "Até 3 por Spec +
  Makefile/CI (Recommended)", option text: "Até 3 correções por Spec,
  contando as rodadas do review. Fica liberado mexer em Makefile, workflows
  do CI, go.mod e lint quando a correção exigir; cada mudança vai descrita no
  relatório."

These answers keep the earlier cycle grants of every owned skill: the
2026-09-30 answer "considere autorizado a ajustar todas as skills se
necessário" and the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths
declarados".

The corrections answer makes the Makefile, `.github/workflows`, `go.mod`,
`go.sum` and the lint and formatter configuration grantable when a correction
needs them. This design needs none of them, and none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch`; the
maintainer's answer releases it as v0.143.1 right after its merge, and that
release may change only the three version-only files in `release_paths`. No
Task performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-13 on the branch `docs/plan-0388`,
based on `f953f77f`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file task_01 declares under `interface:` or `creates:`,
  `docs/references/behavior-surfaces.json`, every planning file of this Spec,
  its adopted reference and that reference's former path.
- The three paths above are governed: the `roundfix` skill's `SKILL.md` pair
  and the canonical release reference.

## What is not governed

- `skills/roundfix/references/release.md`, `skills/roundfix_modes_test.go`
  and `skills/testdata/owned-skill-versions.json`.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the skill mirror, the raised version
fields and the owned-skill version record. This declaration records the
regeneration that follows the approved edit and adds no source path.

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

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill, any Baseline module or derived
  Baseline file, any guide under `docs/agents/` or `docs/user-guide/`,
  `CONTEXT.md`, any ADR, any Archive Record, any Go production file, any file
  under `internal/`, any skill other than `roundfix`, any roundfix reference
  other than `references/release.md`, any section of it other than
  `## Completion notices`, or the `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal.
- No Task sends an ai-memory message, a completion notice or a Spec Event
  Announcement; the operator does.
