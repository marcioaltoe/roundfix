---
status: approved
granted: 2026-10-22
action: make the Linux Lane script mark a failed package-selection go list with its exit code and exit 4, make the lane name go list and that code as a Lane Change Failure below 128 and as a Lane Container Failure with its signal from 128, keep every other lane outcome byte for byte, and say so in the deliver guide, the Roundfix deliver reference and the glossary
consuming: 0303-the-linux-lane-names-the-step-that-failed
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
---

# Approved authority for Spec 0303

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion; the
questions of about 15:25 BRT were answered by about 18:30 BRT. The answers
this Spec rests on were:

- scope (multi-select, all four) "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is one of the adopter fixes of batch 1, the
  operator's bullet "0303 the Linux Lane names the step that failed";
- order "Correções, follow-up, Python (Recomendado)" (option text: "Segue a
  ordem que você definiu para o follow-up: primeiro o que trava outros
  projetos (Light Tier no Fluxus), depois o bom funcionamento, depois as
  features.");
- authorization: the same as cycle AO, that is the Governed Paths each Spec
  declares, and still not this repository's Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

This Spec's design is the operator's reading within that authorization: the
lane script marks a failed `go list` explicitly with its exit code, the cause
names `go list`, the Park Class follows the measured causes, and `go test`
failures keep their bytes (ADR-0316).

The cycle AO rules that answer carries were approved on 2026-10-22, at about
00:55 BRT ("the same rules as cycle AN"), which carry the cycle AN approval of
2026-10-21, at about 01:00 BRT, the cycle AM approval of 2026-10-20, at about
18:15–18:35 BRT, the cycle AL approval of 2026-10-20, at about 13:10 BRT
("Mesmas regras do AK (Recomendado)"), and the cycle AJ grant of 2026-10-19,
at about 16:20 BRT, "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and "Merge squash + release + push". The
2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row starts a container, calls a real Docker daemon, pulls an
image, reads or writes the real `~/.roundfix`, reaches the network or reads
another repository.

The governed set was measured on the authoring branch `docs/plan-0303`, based
on `42f523c2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference, the Backlog Entry path it names, ADR-0316, ADR-0294,
  `docs/references/skill-coverage.json` and `CONTEXT.md`.
- Three of those files are governed, and all three are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` describes the Linux Lane's
  outcomes and gains the `go list` cause (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).

## What is not governed

- The production and test sources task_02 changes or creates in
  `internal/cli` and `internal/delivery`.
- The skill reference mirror `skills/roundfix/references/deliver.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guide `docs/user-guide/commands/deliver.md`.
- Every planning file of this Spec, its adopted reference, ADR-0316 and
  ADR-0294.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
and the Behavior Surface Record after task_01's edits. This declaration
records the regeneration that follows the approved edits and adds no source
path.

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
No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json` or the `### QA
  settlement` section of any skill.
- No change to the lane's package selection, container arguments, mounts,
  environment, skips, image resolution, `roundfix deliver lane`'s exit-code
  mapping or the Park Class mapping.
- No test, Verification command or QA row:
  - starts a container or calls a real Docker daemon;
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
