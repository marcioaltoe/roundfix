---
status: approved
granted: 2026-10-20
action: make roundfix runs causes read each failed Verification attempt's diagnostic under the Artifact Root the Run wrote to, through the one resolver the storage report, GC sanitation and Run Retention use, keep its symbolic-link and containment checks, and describe the change in the runs guide, the Roundfix runs reference and the glossary
consuming: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runs.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0292

On 2026-10-20, between about 18:15 and 18:35 BRT, the maintainer approved
cycle AM through AskUserQuestion. The answers this Spec rests on were:

- scope "R → S → T (Recomendado)" (option text: "Três Specs, um minor cada
  (v0.101–v0.103), nesta ordem. R primeiro porque já afeta o oraculum. Mesmas
  regras do AL.");
- minor items "Distribuir + U (Recomendado)" (option text: "'runs causes' vai
  em T (mesmo resolvedor e storage). 'Esperar vaga' vai em S (gate e
  preflight da fila). O endurecimento de symlink vira a Spec U, no fim
  (v0.104.0)."), then "Rodar gc; T vira runs causes (Recomendado)": this Spec
  is only the `runs causes` Backlog Entry;
- authorization: the same rules as cycle AL, that is the Governed Paths each
  Spec declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AL answer "Mesmas regras do AK (Recomendado)" carries the cycle AJ
rules approved on 2026-10-19, at about 16:20 BRT, among them "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
squash + release + push". The 2026-09-30 skills grant ("considere autorizado a
ajustar todas as skills se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`.

The governed set was measured on the authoring branch `docs/plan-0292`, based
on `70c78d15`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted Backlog Entry at both paths, ADR-0301, ADR-0214,
  `docs/references/corrective-causes-measurement.md`,
  `docs/references/skill-coverage.json` and `.roundfixrc.yml`.
- Four of those files are governed. Three are changed and listed above.
  `.roundfixrc.yml` is governed too, but no Task changes it, so it is not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_02).
- `.agents/skills/roundfix/references/runs.md` covers `runs causes`, whose
  diagnostic read changes (task_02).

## What is not governed

- The production and test sources the Tasks change in `internal/runcause` and
  `internal/cli`.
- The skill reference mirror `skills/roundfix/references/runs.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guide `docs/user-guide/commands/runs.md`.
- Every planning file of this Spec and its adopted reference.

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
surface is added, so `docs/references/skill-coverage.json` does not change. No
record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, `internal/config/`,
  `internal/store/` or the `### QA settlement` section of any skill.
- No change to the resolver, to the storage report, GC sanitation or Run
  Retention, to the signature table, to the `runs causes` flags, help, exit
  codes or JSON fields, or to any other reader of a recorded Artifact
  Directory.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
