---
status: approved
granted: 2026-11-24
action: remove the Linux Lane from the Delivery Queue's repository gate with its on-demand command, its doctor check and its two configuration keys, which become deprecated keys, move the verification of pull requests, of pushes to main and of releases to a standard macOS runner, compile the tests for Linux in the contract that compiles them for Windows, and document it in the user guide, the Roundfix skill, the repository's own guide, the glossary and ADR-0452
consuming: 0420-roundfix-stops-verifying-on-linux
paths:
  - .roundfixrc.yml
  - Makefile
  - .github/workflows/ci-verify.yml
  - .github/workflows/release.yml
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/runtime.md
  - skills/roundfix/SKILL.md
  - docs/references/coverage-record.json
  - docs/agents/specific-repository.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0420

The maintainer's decisions this record rests on are copied whole in [`references/sources.md`](references/sources.md).

## Where the grant comes from

Two decisions of the maintainer of 2026-11-24, each delivered as a `DECISÃO DE MARCIO` block to this project's ai-memory inbox by the Praetor, the session that coordinates the software projects and puts their decision requests to the maintainer.

The Spec, its tooling paths, its delivery and its release: ai-memory message `01a20d73-109f-7bb2-a2b6-718854f6b25d`, decided 09:10 BRT (12:10 UTC). The block, verbatim:

```text
DECISÃO DE MARCIO
pedido: PEDIDO de 24/11 09:05 BRT (Linux Lane da 0418)
pergunta: Roundfix, portão de testes em Linux que roda no Docker do mini antes do PR: Marcio não usa o Roundfix em Linux. Como fica, a começar pela 0418 que está parada nele? (Tirar do local, só no CI; tirar de vez; subir a memória por ora)
resposta: "Tirar de vez" — com o contexto que Marcio escreveu antes: "Eu estou quase pedindo para remover esses testes do linux, não uso o roundfix no linux, o quaestor-legacy só roda no mini e o quaestor vai rodar em container. Isso está gerando tanto problema. Talvez só deixar algum tipo de teste no CI antes do merge na main."
alcance: (1) Spec própria no repositório roundfix para remover a verificação Linux, tanto a Linux Lane local quanto o job Linux do CI, ficando só macOS; inclui os caminhos de tooling necessários (Makefile, CI, configuração da Lane), implementação, PR, merge e release patch. (2) A 0418 pode seguir sem o portão Linux: se a Lane não puder ser desligada sem a Spec nova, priorize essa Spec antes de retomar a 0418, dentro das filas novas já autorizadas (01a20d32). Se a remoção afetar outros repositórios que usam o Roundfix (Lane em quaestor, vortex etc.), descreva o efeito no PRD e avise o Praetor.
quando: 24/11/2026 09:10 BRT
```

The eight Governed Paths that decision does not name, and their regeneration: ai-memory message `01a20f32-07af-7b11-a675-a2cbc7a64743`, decided 17:18 BRT (20:18 UTC). The author of the PRD returned the question as a Reserved Question, the Supervisor session sent it as the decision request `01a20d89-cf32`, and the answer came back as this block, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20d89-cf32
pergunta: Roundfix, Spec 0420 (tirar de vez os testes Linux): ela precisa mudar 8 arquivos protegidos que a decisão não citou — os textos da skill do Roundfix que descrevem o portão Linux, o registro de cobertura dos testes e uma frase do guia que chama o CI de "portão Linux". Pode? (Os 8 arquivos; menos a frase do guia; não, só desligar)
resposta: "Os 8 arquivos (Recomendado)"
alcance: o alcance do pedido 01a20d89-cf32, até o release patch: .agents/skills/roundfix/SKILL.md, .agents/skills/roundfix/references/{deliver,setup,implement,runtime}.md, skills/roundfix/SKILL.md, docs/references/coverage-record.json (só pelo comando próprio; só saem testes da Lane, do comando dela, da checagem do doctor e das duas chaves), docs/agents/specific-repository.md (só a frase "CI is the Linux gate"), com as regenerações listadas no pedido; mais .roundfixrc.yml, Makefile, .github/workflows/ci-verify.yml e release.yml (CI para macOS). Sem go.mod, lint, Baseline nem outra skill. Chaves antigas aceitas com aviso.
quando: 24/11/2026 17:18 BRT
```

## What the grant covers

- The operations `implement`, `commit`, `push`, `pull_request`, `merge` and `release`.
- Four Governed Paths the first decision names by their kind ("Makefile, CI, configuração da Lane") and the second decision names one by one:
  - `.roundfixrc.yml`: the line `linux_lane` under `delivery` and the comment lines above it, and nothing else;
  - `.github/workflows/ci-verify.yml`: the runner of the job `Verification gate` and the comments that name it. The workflow's name, the job's name, its timeout, its triggers, its permissions, its steps, their order, their conditions and their commands stay;
  - `.github/workflows/release.yml`: where the steps `Verify gate` and `Verify every contract` run, and what that needs for them to run before any build or publication. The triggers, the permissions, the tag validation, the publication preflight, the cross-compile, the comparison, the npm publication and the GitHub Release keep their steps and their runner;
  - `Makefile`: the comment above `verify-changed` that calls continuous integration the Linux gate, the value of `SUITE_BUDGET_SECONDS` with the comment that states its measurement, and the name of the compile contract in `REPO_CONTRACT_TESTS` when the contract is renamed. No target, recipe or other variable changes.
- Eight Governed Paths only the second decision names:
  - `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/implement.md` and `.agents/skills/roundfix/references/runtime.md`: the sections, table rows and sentences that describe the Linux Lane, its command, its doctor check and its codes, and the list of deprecated keys in the `setup` reference;
  - `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`: the version fields the record command raises, and nothing else;
  - `docs/references/coverage-record.json`, written only by its own command under `## Sanctioned regeneration`. The only test names it may lose are the ones the Task Graph lists, each a test of the Linux Lane, of its command, of its doctor check or of its two keys;
  - `docs/agents/specific-repository.md`: the sentence, outside the Baseline's markers, that calls continuous integration the Linux gate.
- The regeneration by the repository's own commands, declared under `## Sanctioned regeneration`.
- A release at a `patch` Release Ceiling. The Spec is a `type: refactor`, which the release plan classifies by hand; the operator classifies it as a patch, with the first decision as the reason, and the Supervisor session plans it as `v0.153.3`. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

The grant names no number of corrections and no path beyond these twelve. Anything outside this scope needs a new decision request.

## How the governed set was measured

The author of the PRD measured it on 2026-11-24 with `speccheck.GovernedPath`, in a `go test -overlay` probe whose test file lives outside the repository: at `808f9001` over 82 files the change is expected to touch, and again at `cefd4e82`, after Spec 0418 and Spec 0419 merged, over 74 files. The function's source is the same at both commits, and both probes return the same twelve paths. The list of governed and ordinary files is in `references/sources.md`, section 4.

- The twelve paths above are the governed ones this Spec changes.
- These governed paths are not granted, and no Task may change them:
  - `.github/workflows/ci-conventions.yml` and `.github/workflows/secondbrain-sync.yml`, which run no test;
  - `go.mod`, `go.sum` and the lint and formatter configuration;
  - `internal/cli/cli_test.go` and every other historically bounded test or source file;
  - every other reference of the Roundfix skill, and every other owned skill;
  - every Vendored Skill under `.agents/skills` and `skills-lock.json`'s entries;
  - every Baseline module, template, profile, Source Baseline and Setup Snapshot, the starter workflow `internal/baseline/assets/templates/starters/python/verify.yml` included: it names a Linux runner for an adopter's repository and stays;
  - every other guide under `docs/agents/`, every line of `docs/agents/specific-repository.md` inside the Baseline's markers, and `docs/agents/setup-context.json`.

The TechSpec author and the Tasks author measure again every file their artifacts name. When that set holds a Governed Path outside the twelve above, no Run starts: the path needs a new decision of the maintainer, recorded here on the default branch before the consuming delivery.

## What is not governed

- The Go source and test files of the Delivery Queue, of its commands, of the configuration loader, of doctor and of the health checker, and the tests of `internal/verifyselect` that read the workflows and compile the tests for other platforms.
- The mirrors `skills/roundfix/references/deliver.md`, `setup.md`, `implement.md` and `runtime.md`, and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/deliver.md`, `docs/user-guide/commands/doctor.md`, `docs/user-guide/configuration.md`, `README.md`, `CONTEXT.md` and `.secondbrain-export`.
- ADR-0452, the five decision records it supersedes, which are deleted, and ADR-0252, ADR-0285, ADR-0295, ADR-0384, ADR-0385 and ADR-0423, which gain their amendment notices.
- `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json`.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the skill mirror, the skill version fields and the generated records. This declaration records the regeneration that follows the approved edits and adds no source path. The commands are the ones the decision request `01a20d89-cf32` lists and the second decision grants.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

No mirror, version field, digest or record is hand-edited. Every Spec that touches the Roundfix skill raises the same version fields and re-records the same records, as Spec 0418 and Spec 0419 did on 2026-11-24: every command above runs on the tree the Task finds, and the Task writes the next version, never a number.

## Limits

- No change to any of these:
  - `go.mod` or `go.sum`, and the lint, formatter or test-runner configuration;
  - any target, recipe or variable of the Makefile beyond the three places named above;
  - the workflows that validate a pull request title and mirror files to the Secondbrain;
  - `skills-lock.json`, any Vendored Skill, any Baseline module, template, profile, Setup Snapshot or plan golden;
  - any guide under `docs/agents/` other than the one sentence named above, and `docs/agents/setup-context.json`;
  - any Archive Record, and `docs/references/coverage-record.json` beyond its own regeneration command;
  - any skill other than `roundfix`, and any Roundfix reference other than `deliver`, `setup`, `implement` and `runtime`;
  - the `### QA settlement` section of any skill;
  - `dist/npm/platforms.json` and the platforms a release publishes;
  - any accepted ADR other than the one this Spec adds, the five it supersedes and the six that gain their amendment notices.
- No larger runner, no paid plan, no self-hosted runner, and no setting of the repository on GitHub: the required checks of the default branch keep their names.
- No new column, table or stage in the Run Database, no new command or flag, and no new key of Project Config.
- Outside this Spec: every pending message of this project's inbox, and the behavior Spec 0418 and Spec 0419 delivered.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh` or the real `docker`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session, an ACP adapter or a container;
  - reads another repository, except a temporary one the test itself created;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement. The operator tells the Praetor what the removal changes for the other repositories, as the first decision asks.
