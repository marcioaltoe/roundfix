---
status: approved
granted: 2026-10-20
action: make the Delivery Queue record a Linux Lane skipped for an environment reason on the item and in the Pull Request body, add the roundfix doctor linux-lane check, make an item refused only by the Active Implement Run ceiling wait for a Run slot up to runs.slot_wait, and describe both in the user guide, the Roundfix deliver, setup and runtime references and the glossary
consuming: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/references/setup.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0291

On 2026-10-20, at about 18:15–18:35 BRT, the maintainer approved cycle AM
through AskUserQuestion. The answers this Spec rests on were:

- scope "R → S → T (Recomendado)" (option text: "Três Specs, um minor cada
  (v0.101–v0.103), nesta ordem. R primeiro porque já afeta o oraculum. Mesmas
  regras do AL.");
- minor items "Distribuir + U (Recomendado)" (option text: "'runs causes' vai
  em T (mesmo resolvedor e storage). 'Esperar vaga' vai em S (gate e
  preflight da fila). O endurecimento de symlink vira a Spec U, no fim
  (v0.104.0)."), then "Rodar gc; T vira runs causes (Recomendado)";
- for this Spec's wait, "Espera com limite" (option text: "Espera até um
  prazo configurável (padrão 2h) e depois estaciona como hoje.");
- authorization: the same rules as cycle AL — the Governed Paths each Spec
  declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AL authorization "Mesmas regras do AK (Recomendado)" carries the
cycle AJ rules approved on 2026-10-19, at about 16:20 BRT, among them "Mexer
em Governed Paths declarados" (option text: "'Concedo' para os caminhos que
cada Spec declarar (módulos Baseline, registros, skills, testes governados;
fora: Makefile, CI, go.mod, lint, que pedem autorização específica).") and
"Merge squash + release + push". The 2026-09-30 skills grant ("considere
autorizado a ajustar todas as skills se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row runs real Docker, reads or writes the real `~/.roundfix`, or
reaches the network; the maintainer runs `roundfix doctor` against the real
Docker daemon himself after the release.

The governed set was measured on the authoring branch `docs/plan-0291`, based
on `70c78d15`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its two adopted Backlog Entries, ADR-0303, ADR-0304,
  `docs/references/skill-coverage.json` and `.roundfixrc.yml`.
- Six of those files are governed. Five are changed and listed above.
  `.roundfixrc.yml` is governed too, but no Task changes it, so it is not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_02, task_04).
- `.agents/skills/roundfix/references/deliver.md` covers `roundfix deliver`,
  whose gate now records a Lane Environment Skip (task_02) and whose queue now
  waits for a Run slot (task_04).
- `.agents/skills/roundfix/references/setup.md` covers `roundfix doctor`'s
  findings table, which gains the four `DR-LANE-*` codes (task_02).
- `.agents/skills/roundfix/references/runtime.md` lists the checks `roundfix
  doctor` runs, which gain `linux-lane` (task_02).

## What is not governed

- The production and test sources the Tasks change in `internal/cli`,
  `internal/delivery`, `internal/config` and `internal/store`.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guide pages `docs/user-guide/commands/deliver.md`,
  `docs/user-guide/commands/doctor.md` and `docs/user-guide/configuration.md`.
- Every planning file of this Spec, its adopted references, ADR-0303 and
  ADR-0304.

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
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, the Run Database schema
  or the `### QA settlement` section of any skill.
- No change to the Linux Lane's console and log notices, skip reasons and
  outcomes, `roundfix deliver lane`'s exit codes, `roundfix implement`'s
  output and exit codes, or the store's ceiling check.
- No test, Verification command or QA row:
  - runs real Docker or a real `roundfix doctor` outside a test;
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge, a reviewer or GitHub;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
