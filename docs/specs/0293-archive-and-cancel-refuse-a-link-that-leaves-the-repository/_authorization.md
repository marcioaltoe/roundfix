---
status: approved
granted: 2026-10-20
action: make the Archive Command and the Cancel Command refuse, before writing anything, every path and destination whose resolved target lies outside the repository, with the refusal text each step already prints, while a link that stays inside, a checkout reached through a link and an external Spec Root keep working, and state the refusal in the archive and cancel user guides, the Roundfix archive reference and ADR-0269
consuming: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
paths:
  - internal/spec/archive.go
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0293

On 2026-10-20, at about 18:15–18:35 BRT, the maintainer approved cycle AM
through AskUserQuestion. The answers this Spec rests on were:

- scope "R → S → T (Recomendado)" (option text: "Três Specs, um minor cada
  (v0.101–v0.103), nesta ordem. R primeiro porque já afeta o oraculum. Mesmas
  regras do AL.");
- minor items "Distribuir + U (Recomendado)" (option text: "'runs causes' vai
  em T (mesmo resolvedor e storage). 'Esperar vaga' vai em S (gate e
  preflight da fila). O endurecimento de symlink vira a Spec U, no fim
  (v0.104.0)."), then "Rodar gc; T vira runs causes (Recomendado)";
- for this Spec, "Recusar no archive/cancel (Recomendado)" (option text:
  "`archive` e `cancel` recusam esse caminho com a mesma mensagem da 0288,
  nomeando o caminho. Hoje ele é aceito, então a Spec declara a quebra.");
- authorization: the same rules as cycle AL, the Governed Paths each Spec
  declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AL rules that answer carries ("Mesmas regras do AK (Recomendado)",
2026-10-20, about 13:10 BRT) carry the cycle AK rules ("Sim, mesmas regras do
AJ (Recomendado)", 2026-10-20, about 06:30 BRT), and those carry the cycle AJ
grant of 2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths declarados"
(option text: "'Concedo' para os caminhos que cada Spec declarar (módulos
Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
lint, que pedem autorização específica).") and "Merge squash + release +
push". The 2026-09-30 skills grant ("considere autorizado a ajustar todas as
skills se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, changes
`.roundfixrc.yml`, reaches the network or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0293`, based
on `70c78d15`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted Backlog Entry, ADR-0306, ADR-0269, `CONTEXT.md`,
  `docs/references/skill-coverage.json` and `.roundfixrc.yml`.
- Six of those files are governed. Four are changed and listed above.
  `internal/spec/archive_test.go` is governed too, but task_01 only runs
  tests from it; `.roundfixrc.yml` is governed, but no Task changes it. Neither
  is bounded.

## Why each governed path is unavoidable

- `internal/spec/archive.go` holds `Archive`, whose promotion destination is
  joined to the repository root and written without a containment test
  (task_01).
- `.agents/skills/roundfix/references/archive.md` covers `archive` and
  `cancel`, whose refusals change (task_03).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_03).

## What is not governed

- The other production and test sources the Tasks change in `internal/spec`
  and `internal/cli`.
- The user guides `docs/user-guide/commands/archive.md` and
  `docs/user-guide/commands/cancel.md`, and ADR-0269.
- The skill reference mirror `skills/roundfix/references/archive.md`,
  `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec, its adopted Backlog Entry and ADR-0306.

## Sanctioned regeneration

The repository-owned commands resolve the Behavior Surface Record after
task_02 adds a paragraph to the `archive` and `cancel` guides, and the skill
mirror and the skill version after task_03 edits the Roundfix `archive`
reference. This declaration records the regeneration that follows the
approved edits and adds no source path.

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
- No change to what the Archive Command or the Cancel Command decide for a
  path whose target stays inside the repository: the QA eligibility, the QA
  Archive Override, the existing refusals and their order, the Archive Record
  schema and fields, the Archive Relink and the Legacy Link Rebase stay as
  they are, and every Archive Record already written stays byte-identical.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
