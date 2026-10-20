---
status: approved
granted: 2026-10-20
action: make the profile proof's setup deadline the config key defaults.profile_proof_setup_timeout, give a timed-out proof a third attempt after a 5 and 10 second backoff, close an absent disposable Agent Session without a warning, re-record the coverage record after one test rename, and describe the change in the configuration and profiles guides, the Roundfix profiles and runtime references and the glossary
consuming: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/profiles.md
  - .agents/skills/roundfix/references/runtime.md
  - docs/references/coverage-record.json
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0285

On 2026-10-20, at about 06:30 BRT, the maintainer approved cycle AK through
AskUserQuestion. The answers this Spec rests on were:

- next cycle "Duas Specs + um lote (Recomendado)": M (0285, this Spec), then
  N (0286), then O (0287), one minor release each;
- this Spec's design "Retry interno + prazo configurável (Recomendado)"
  (option text: "Um timeout classificado 'temporary' ganha uma nova tentativa
  automática dentro do preflight antes de estacionar, e o prazo de criação da
  sessão passa a ser configurável (documentado em profiles.md). Continua
  provando todas as tuplas.");
- authorization "Sim, mesmas regras do AJ (Recomendado)": the cycle AJ rules
  stand, including "Mexer em Governed Paths declarados" (option text:
  "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
  registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
  pedem autorização específica).").

The 2026-09-30 skills grant ("considere autorizado a ajustar todas as skills
se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0285`, based
on `9448aa64`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before and
  after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  the adopted Inbox Entry, ADR-0298, `docs/references/skill-coverage.json` and
  `.roundfixrc.yml`.
- Six of those files are governed. Five are changed and listed above.
  `.roundfixrc.yml` is governed too, but no Task changes it, so it is not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01 and
  task_04).
- `.agents/skills/roundfix/references/profiles.md` documents the new key
  (task_01) and the attempts (task_04), as the maintainer's option text asks.
- `.agents/skills/roundfix/references/runtime.md` says today that a
  timed-out proof "is retried once", which this Spec changes (task_04).
- `docs/references/coverage-record.json` lists every test name; task_02
  renames `TestProfileProofReportsASecondTimeoutAsTemporary` to
  `TestProfileProofReportsAThirdTimeoutAsTemporary`, so the record is
  re-recorded with its own command.

## What is not governed

- The production and test sources in `internal/config`, `internal/agent` and
  `internal/cli`.
- The skill reference mirrors `skills/roundfix/references/profiles.md` and
  `skills/roundfix/references/runtime.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, the user guides `docs/user-guide/configuration.md` and
  `docs/user-guide/commands/profiles.md`, and ADR-0298.
- Every planning file of this Spec, including its adopted reference.

## Sanctioned regeneration

The repository-owned commands resolve the coverage record, the skill mirror,
the skill version and the Behavior Surface Record. This declaration records
the regeneration that follows the approved edits and adds no source path.

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

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

No Baseline module changes, so no module version is recorded and `make
baseline-digests` has nothing to rewrite. No command surface is added, so
`docs/references/skill-coverage.json` does not change. No record or mirror is
hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, the `roundfix init`
  template or the `### QA settlement` section of any skill.
- No change to the classification of proof failures other than the setup
  timeout, the cleanup deadline, the Fallback Chain probe after Run creation,
  any Park Class or `roundfix deliver retry`.
- Every configured tuple is still proven.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts a real ACP agent, adapter or acpx;
  - reaches the network, a provider, the judge or a reviewer;
  - needs Docker;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
