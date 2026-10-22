---
status: approved
granted: 2026-10-22
action: make the Delivery Queue's check wait the User Config key delivery.check_timeout, a Go duration greater than 0 and at most 24h with a default of 60 minutes, ignored with a warning in Project Config, make the checks-timeout park name the key and the value it waited, keep each wait's own deadline, and say so in the configuration and deliver guides, the Roundfix deliver reference and the glossary
consuming: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
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

# Approved authority for Spec 0298

On 2026-10-22, at about 00:55 BRT, the maintainer approved cycle AO through
AskUserQuestion. The answers this Spec rests on were:

- scope, all four options in this order: "Z: portão e corretiva", "Espera de
  checks", "Retomar item arquivado", "Y: reconcile PID", one Spec and one
  minor each (v0.108–v0.111): this Spec is "Espera de checks" and delivers
  second;
- check wait "Config com padrão 60m (Recomendado)" (option text: "Nova chave
  de User Config (ex.: `delivery.check_timeout`) com padrão 60 min, igual ao
  timeout atual do job ci-verify. O operador ajusta sem release.");
- docscontract performance "Backlog por enquanto (Recomendado)": not part of
  this Spec;
- authorization: the same rules as cycle AN, that is the Governed Paths each
  Spec declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AN rules that answer carries were approved on 2026-10-21, at about
01:00 BRT ("the same rules as cycle AM"), which carry the cycle AM approval of
2026-10-20, at about 18:15–18:35 BRT ("the same rules as cycle AL"), the cycle
AL approval of 2026-10-20, at about 13:10 BRT ("Mesmas regras do AK
(Recomendado)"), and the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
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

The governed set was measured on the authoring branch `docs/plan-0298`, based
on `34b4e246`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference, ADR-0311, `docs/references/skill-coverage.json` and
  `CONTEXT.md`.
- Three of those files are governed, and all three are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` covers every `deliver`
  surface, whose `checks-timeout` park text and check wait change (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).

## What is not governed

- The production and test sources task_02 changes in `internal/config`,
  `internal/delivery` and `internal/cli`.
- The skill reference mirror `skills/roundfix/references/deliver.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, `docs/user-guide/configuration.md` and
  `docs/user-guide/commands/deliver.md`.
- Every planning file of this Spec, its adopted reference and ADR-0311.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill version and
the Behavior Surface Record after task_01's edits. This declaration records
the regeneration that follows the approved edits and adds no source path.

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
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, `internal/store/`,
  `internal/delivery/park_class.go` or the `### QA settlement` section of any
  skill.
- No change to the check poll interval, the re-run rules, the merge policy
  refusal, the queue deadline, the retry limit or `runs.slot_wait`.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
