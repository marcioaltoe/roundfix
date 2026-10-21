---
status: approved
granted: 2026-10-21
action: make the Archive Command rebase each unproven entry of the Archive Record it writes exactly as it rebases the outcome, make the Record Link Repair of roundfix history sanitize and roundfix baseline update and the SC-HISTORY-RETAINED fix line point a promoted target at its docs/references copy, and state both in the archive, history, baseline and spec user guides, the Roundfix archive, baseline and spec references and the glossary
consuming: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
paths:
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0295

On 2026-10-21, at about 01:00 BRT, the maintainer approved cycle AN through
AskUserQuestion. The answers this Spec rests on were:

- scope "V → W → X (Recomendado)" (option text: "Três Specs: v0.105 a
  v0.107. Y fica no backlog como entrada aberta; são só 3 linhas de ruído no
  reconcile e não custam nada.");
- for this Spec, "Archive + alvo promovido (Recomendado)" (option text: "O
  archive reescreve os `satisfied-by` (Git Locator para arquivo removido,
  `docs/references/` para arquivo promovido), e o sanitize passa a apontar
  arquivos promovidos para `docs/references/`. A exigência da tag fica como
  está; ela é a âncora do Git Locator.");
- authorization: the same rules as cycle AM, the Governed Paths each Spec
  declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AM rules that answer carries ("the same rules as cycle AL",
2026-10-20, about 18:15–18:35 BRT) carry the cycle AL rules ("Mesmas regras
do AK (Recomendado)", 2026-10-20, about 13:10 BRT), the cycle AK rules ("Sim,
mesmas regras do AJ (Recomendado)", 2026-10-20, about 06:30 BRT), and those
carry the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
squash + release + push". The 2026-09-30 skills grant ("considere autorizado
a ajustar todas as skills se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, changes
`.roundfixrc.yml`, reaches the network or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0295`, based
on `956066a7`, with `speccheck.GovernedPath`, and measured again on
`8fb2da1f` after the base moved.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted Inbox Entry copy, ADR-0308, `CONTEXT.md`,
  `docs/references/skill-coverage.json`, `internal/spec/archive.go`,
  `internal/spec/archive_test.go`, `internal/spec/history_sanitize.go`,
  `internal/spec/record_links_test.go` and `.roundfixrc.yml`.
- Eight of those files are governed. Five are changed and listed above.
  `internal/spec/archive.go` and `internal/spec/archive_test.go` are
  governed, but no Task changes either: this Spec leaves `archive.go` to Spec
  0293 and adds its tests in new files. `.roundfixrc.yml` is governed, but no
  Task changes it. None of the three is bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/archive.md` covers `archive` and
  `history sanitize`, whose record text and repair change (task_03).
- `.agents/skills/roundfix/references/baseline.md` covers the history section
  of `baseline update`, whose repair changes (task_03).
- `.agents/skills/roundfix/references/spec.md` covers `SC-HISTORY-RETAINED`,
  whose fix line changes (task_03).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_03).

## What is not governed

- The production and test sources the Tasks change in `internal/spec`,
  `internal/speccheck` and `internal/cli`.
- The user guides `docs/user-guide/commands/archive.md`,
  `docs/user-guide/commands/history.md`,
  `docs/user-guide/commands/baseline.md` and
  `docs/user-guide/commands/spec.md`, and `CONTEXT.md`.
- The skill reference mirrors under `skills/roundfix/references/`,
  `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec, its adopted Inbox Entry copy and
  ADR-0308.

## Sanctioned regeneration

The repository-owned commands resolve the Behavior Surface Record after
task_01 and task_02 edit the `archive`, `history`, `baseline` and `spec`
guides, and the skill mirror and the skill version after task_03 edits the
Roundfix references. This declaration records the regeneration that follows
the approved edits and adds no source path.

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

No Baseline module changes, so no module version is recorded. An owned skill
edit leaves the derived Baseline artifacts byte-identical, so `make
baseline-digests` has nothing to rewrite. No test is removed or renamed and no
platform set changes, so the coverage record is not re-recorded. No command
surface is added, so `docs/references/skill-coverage.json` does not change.
No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, `internal/config/` or
  the `### QA settlement` section of any skill.
- No change to `internal/spec/archive.go`, `internal/spec/repository_path.go`,
  `internal/spec/cancel.go`, `internal/cli/archive.go`,
  `internal/cli/cancel.go` or the containment logic Spec 0293 adds.
- No change to what the Archive Command decides: the QA eligibility, the QA
  Archive Override, the Archive Relink, the promotions, the refusals and
  their order, the record schema and fields, and the Delivery Queue's
  exact-archive and Pull Request Record Commit proofs stay as they are.
  Every Archive Record already written stays byte-identical unless an
  operator applies the Record Link Repair.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment in this repository. Verification stays
  Daemon-owned, and Task status stays Daemon-written.
