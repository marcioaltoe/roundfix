---
status: approved
granted: 2026-10-16
action: run the changed Go packages' tests in a Linux Lane of the Delivery Queue gate and opt this repository in through Project Config; let reconciliation release a delivered Spec's leftovers beneath the paths its Tasks declared or recorded and name why it keeps any other; reclaim merged and unqueued items' Item Binaries at the merge and through the GC Command; and describe each change in the user guides, the Roundfix skill references and the glossary
consuming: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/reconcile.md
  - .agents/skills/roundfix/references/storage.md
  - .roundfixrc.yml
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0279

On 2026-10-16 the maintainer approved cycle AH through AskUserQuestion. The
answers were the authorization "Governed Paths declarados" (option text:
"\"Concedo\" como em AE/AF/AG para os Governed Paths que cada Spec declarar:
módulos da Baseline, registros de versão, make baseline-digests, skills e
testes governados. Makefile, workflows de CI, go.mod e configuração de lint
ficam fora.") and the failure rules "Regras de AE–AG (Recomendado)". Both carry
into cycle AI.

On 2026-10-18, at about 13:30 BRT, at the end of cycle AH, the maintainer chose
the items "Verificação também no Linux, Runs sujos de Spec entregue,
Crescimento da store, Adopters que ficaram fora", chose "Remover agora
(Recomendado)" for the three dirty 0272 worktrees, and chose the mode
"Escrever, entregar e lançar (Recomendado)". The 2026-09-30 skills grant
("considere autorizado a ajustar todas as skills se necessário") also stands.

The governed set was measured on the authoring branch `docs/plan-0279`, based
on `2ef4f080`, with `speccheck.GovernedPath`.
- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status` was identical before and after, and
  the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and the ADRs the plan writes or amends.
- Six of those files are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01, task_02,
  task_03).
- Three Roundfix references are governed:
  - `.agents/skills/roundfix/references/deliver.md` covers `deliver`, whose
    gate gains the Linux Lane (task_01) and whose merged release gains the Item
    Binary reclamation (task_03).
  - `.agents/skills/roundfix/references/reconcile.md` covers `reconcile`, whose
    scope and dirty reason change (task_02).
  - `.agents/skills/roundfix/references/storage.md` covers `gc`, which gains
    the Item Binary sweep and the retention policy (task_03).
- `.roundfixrc.yml` is this repository's Project Config. It gains
  `delivery.linux_lane: true` so the repository opts in (task_01). Nothing else
  in it changes.

## What is not governed

- The production and test sources the Tasks change in `internal/config`,
  `internal/cli` and `internal/worktree`, including every new file.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, the user guides under `docs/user-guide/` and the ADRs.
- Every planning file of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
and the Behavior Surface Record. This declaration records the regeneration
that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No Baseline module changes, so no module version is recorded. An owned skill
edit leaves the derived Baseline artifacts byte-identical, so `make
baseline-digests` has nothing to rewrite. No test is removed or renamed and no
platform set changes, so the coverage record is not re-recorded. The Makefile
does not change, so the Setup Manifest does not move. No record or mirror is
hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `docs/agents/setup-context.json`, `docs/references/coverage-record.json`, or
  the `### QA settlement` section of any skill.
- In `.roundfixrc.yml`, only the `delivery.linux_lane: true` line and its
  comment.
- No change to `internal/delivery/engine.go`, the commit proof of
  reconciliation, Journal Retention, Run Retention, `internal/judge` or any
  judge question.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network or pulls an image;
  - starts a real container, except the QA row that runs one only with an
    image already present locally and a read-only mount;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.

## Maintainer approval of 2026-10-18

The `.roundfixrc.yml` line is expressly authorized: on 2026-10-18 the
maintainer answered through AskUserQuestion "Sim aos dois (Recomendado)"
(option text: "Autorizo a linha `delivery.linux_lane: true` no .roundfixrc.yml,
com emenda datada na autorização, e um `docker pull golang:1.26` único para o
lane rodar de verdade. O lane nunca faz pull sozinho."). The image pull is an
operator step outside every Task; no Task, Verification or QA row pulls an
image.
