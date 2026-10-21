---
status: approved
granted: 2026-10-21
action: make the Delivery Queue's archive commit subject and pull request title keep the slug while they fit 100 characters and otherwise name the Spec by its Spec Number, the title followed by the Spec title cut at a word so the title plus a squash merge's pull request number fits, keep every reader of those messages working, and say so in the deliver guide, the Roundfix deliver reference, the archive-spec skill and the glossary
consuming: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0296

On 2026-10-21, at about 01:00 BRT, the maintainer approved cycle AN through
AskUserQuestion. The answers this Spec rests on were:

- scope "V → W → X (Recomendado)" (option text: "Três Specs: v0.105 a
  v0.107. Y fica no backlog como entrada aberta; são só 3 linhas de ruído no
  reconcile e não custam nada."): X, this Spec, delivers third;
- this Spec's design "Slug se couber (Recomendado)" (option text: "Continua
  `docs: archive <slug>` e `feat: deliver <slug>` quando cabem em 100 (quase
  sempre). Acima disso, usa o número: `docs: archive the 0290 Spec` e `feat:
  deliver 0290` + título da Spec cortado na palavra. Nada muda nos adotantes
  com slugs curtos.");
- authorization: the same rules as cycle AM, that is the Governed Paths each
  Spec declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AM rules that answer carries were approved on 2026-10-20, at about
18:15–18:35 BRT ("the same rules as cycle AL"), which carry the cycle AL
approval of 2026-10-20, at about 13:10 BRT ("Mesmas regras do AK
(Recomendado)") and the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push". The 2026-09-30 skills
grant ("considere autorizado a ajustar todas as skills se necessário") also
stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, reaches the network,
reads another repository or runs `roundfix deliver` against a real forge.

The governed set was measured on the authoring branch `docs/plan-0296`, based
on `8fb2da1f`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference, ADR-0309, `docs/references/skill-coverage.json` and
  `CONTEXT.md`.
- Five of those files are governed, and all five are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` covers every `deliver`
  surface, whose archive subject and pull request title change (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).
- `.agents/skills/archive-spec/SKILL.md` states the subject the Delivery
  Queue writes for the archive commit, which changes, and with
  `skills/archive-spec/SKILL.md` carries that skill's version fields
  (task_01).

## What is not governed

- The production and test sources task_02 changes in `internal/cli`.
- The skill reference mirror `skills/roundfix/references/deliver.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guide `docs/user-guide/commands/deliver.md`.
- Every planning file of this Spec, its adopted reference and ADR-0309.

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
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, `internal/delivery/`,
  `internal/spec/`, `internal/releaseplan/` or the `### QA settlement`
  section of any skill.
- No change to the Pull Request Record Commit subject, the derived merge
  message, the Pull Request body, the archive-commit proof, the hook refusal
  or the pull request lookup.
- The queue never passes `--no-verify` and never renames a pull request.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
