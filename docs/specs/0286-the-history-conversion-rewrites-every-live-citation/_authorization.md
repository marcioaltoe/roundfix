---
status: approved
granted: 2026-10-20
action: make the history conversion rewrite every citation it lists in every Citation Source, name and date a retired ADR's lineage line by its repository's identifier and supersession day, leave pending Task targets alone and resolve broken links by file name, extend SC-HISTORY-RETAINED to report a removed history citation without a Git Locator, and describe the change in the baseline, history and spec guides, the Roundfix skill references and the glossary
consuming: 0286-the-history-conversion-rewrites-every-live-citation
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/spec.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0286

On 2026-10-20, at about 06:30 BRT, the maintainer approved cycle AK through
AskUserQuestion. The answers this Spec rests on were:

- the cycle: "Duas Specs + um lote (Recomendado)", Specs 0285 → 0286 → 0287,
  one minor each;
- this Spec's scope: "Corrigir os 4 + o check acusa (Recomendado)" (option
  text: "Reescrever como Git Locator toda citação listada, inclusive caminho
  entre crases e relativo. Não tocar link para arquivo que uma Task pendente
  declara em creates:. Resolver link relativo quebrado pelo nome do arquivo.
  Derivar o ID da ADR com nome por data do nome do arquivo ou de
  superseded_by, com a data da substituição. O spec check passa a acusar
  citação viva quebrada para docs/history. A nota sobre --no-history vai para
  a documentação.");
- the authorization: "Sim, mesmas regras do AJ (Recomendado)". It carries the
  cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados" (option
  text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
  registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
  pedem autorização específica).").

The 2026-09-30 skills grant ("considere autorizado a ajustar todas as skills
se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0286`, based
on `9448aa64`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before and
  after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  the adopted Inbox Entry, the three ADRs the plan writes or amends,
  `docs/references/skill-coverage.json`,
  `internal/docscontract/testdata/corpus-golden.json`, `.roundfixrc.yml` and
  the `Makefile`.
- Eight of those files are governed. Five are changed and listed above.
  `internal/docscontract/testdata/corpus-golden.json`, `.roundfixrc.yml` and
  the `Makefile` are governed too, but no Task changes them, so they are not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content, and the canonical
  `SKILL.md` still says the update "never rewrites citing files" (task_04).
- `.agents/skills/roundfix/references/baseline.md` covers `baseline update`,
  whose history section now rewrites every listed citation, names and dates
  the lineage line differently, keeps pending Task targets, resolves broken
  links by file name and gains the `--no-history` note (task_04).
- `.agents/skills/roundfix/references/archive.md` covers `history sanitize`,
  which shares the same conversion (task_04).
- `.agents/skills/roundfix/references/spec.md` covers `spec check`, whose
  `SC-HISTORY-RETAINED` gains the removed-citation finding (task_04).

## What is not governed

- The production and test sources the Tasks change in `internal/spec`,
  `internal/speccheck` and `internal/cli`.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, the user guides `docs/user-guide/commands/baseline.md`,
  `history.md` and `spec.md`, and the ADRs.
- Every planning file of this Spec, including its adopted reference.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirror, the skill version and
the Behavior Surface Record. This declaration records the regeneration that
follows the approved edits and adds no source path.

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
platform set changes, so the coverage record is not re-recorded. No command
surface is added, so `docs/references/skill-coverage.json` does not change.
`SC-HISTORY-RETAINED` keeps its count of 0 on the characterized corpus, so
`internal/docscontract/testdata/corpus-golden.json` does not change. No record
or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`,
  `internal/docscontract/testdata/corpus-golden.json` or the `### QA
  settlement` section of any skill.
- No change to the Archive Relink, `record-links`, Stale History Citations,
  `superseded` records, the Plan Digest's fields or the JSON field set.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag push or deployment. Verification stays Daemon-owned, and
  Task status stays Daemon-written.
