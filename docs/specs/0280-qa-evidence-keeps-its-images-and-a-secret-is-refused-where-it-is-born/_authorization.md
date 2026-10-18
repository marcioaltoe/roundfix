---
status: approved
granted: 2026-10-16
action: accept PNG, JPEG and SVG QA evidence and refuse any other binary at the QA Report commit; run the Secret Scan over the lines each Task commit adds before the Task settles; let a Delivery Retry of a secret-found item accept a Secret Rewrite and print the rewrite that works; and describe each change in the user guides, the Roundfix skill references, the qa-gate and implement-task skills and the glossary
consuming: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
paths:
  - .agents/skills/implement-task/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/implement.md
  - skills/implement-task/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0280

On 2026-10-16 the maintainer approved the cycle through AskUserQuestion. The
answers were the authorization "Governed Paths declarados" (option text:
"\"Concedo\" como em AE/AF/AG para os Governed Paths que cada Spec declarar:
módulos da Baseline, registros de versão, make baseline-digests, skills e
testes governados. Makefile, workflows de CI, go.mod e configuração de lint
ficam fora.") and the failure rules "Regras de AE–AG (Recomendado)". The
2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
necessário") also stands.

On 2026-10-18, at about 14:00 BRT, the maintainer placed this work as "Spec I
depois da 0279 (Recomendado)" and chose, for PNG evidence, "Aceitar imagens
(Recomendado)" (option text: "O detector aceita PNG, JPEG e SVG como evidência
legível e recusa só binário de build ou executável. O Daemon recusa no commit
do relatório qualquer binário que não seja imagem."). For `secret-found` the
choices were "Retry aceita head reescrito", "Scan em cada Task" and
"Documentar fixture" (option texts: "Para secret-found, o retry aceita um head
reescrito que preserva a árvore final, exceto as linhas apontadas, e a
próxima ação impressa explica o caminho real."; "O Secret Scan roda no commit
de cada Task, como já roda no do relatório de QA, e recusa a fixture quando
ela nasce."; "Documentar no guia e na skill como escrever um valor de fixture
que não case a regra, sem afrouxar o detector.").

The governed set was measured on the authoring branch `docs/plan-0280`, based
on `2ef4f080`, with `speccheck.GovernedPath`.
- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status` was identical before and after, and
  the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and the ADRs the plan writes.
- Nine of those files are governed. Eight are changed and listed above.
  `internal/speccheck/mechanical_test.go` is governed too, but task_01 only
  reads its fixtures as an `instruction:` and never changes it, so it is not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01, task_02,
  task_03).
- Two Roundfix references are governed:
  - `.agents/skills/roundfix/references/implement.md` covers `implement`,
    whose QA Report commit gains the non-image binary refusal (task_01) and
    whose Task settlement gains the Secret Scan and the fixture guidance
    (task_02).
  - `.agents/skills/roundfix/references/deliver.md` covers `deliver retry`,
    which gains the Secret Rewrite and the new `secret-found` next action
    (task_03).
- `.agents/skills/qa-gate/SKILL.md` and `skills/qa-gate/SKILL.md` are the
  `qa-gate` skill and its mirror. They gain the rule that evidence files are
  text or images (task_01), with the version fields they carry.
- `.agents/skills/implement-task/SKILL.md` and `skills/implement-task/SKILL.md`
  are the `implement-task` skill and its mirror. They gain the `Fixture
  values` rule (task_02), with the version fields they carry.

## What is not governed

- The production and test sources the Tasks change or create in
  `internal/speccheck`, `internal/daemon`, `internal/secretscan`,
  `internal/delivery` and `internal/cli`, including every new file.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, the user guides under `docs/user-guide/commands/` and the
  ADRs.
- Every planning file of this Spec, including its adopted references.

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
baseline-digests` has nothing to rewrite. No test is removed or renamed and
no platform set changes, so the coverage record is not re-recorded. No record
or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, any Secret Pattern in
  `internal/secretscan/patterns.go`, or the `### QA settlement` section of any
  skill.
- No change to the Final Push or Delivery Queue push scans, the retry of any
  other blocker, the Verification retry policy, `internal/judge`, any judge
  question or threshold, or the Jev Ceiling.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - writes a literal value that matches a Secret Pattern.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
