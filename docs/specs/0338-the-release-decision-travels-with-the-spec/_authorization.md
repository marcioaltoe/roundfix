---
status: approved
granted: 2026-11-01
action: make the PRD's Delivery Type reach the Archive Record and the Delivery Queue's pull request title and be checked by the Spec Check, let an approved authorization record that grants release carry a Release Ceiling and release paths that the release plan reads to answer approval_required within every ceiling, treat Backlog Entries as planning evidence, reword the Baseline's release and authorization clauses and regenerate their derived files and this repository's managed guides, and teach the write-prd, write-techspec, write-tasks and Roundfix skills both rules
consuming: 0338-the-release-decision-travels-with-the-spec
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/release.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/speccheck/constraints.go
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0338

On 2026-11-01, at about 22:20 BRT, the maintainer answered the operator's
cycle AU questions with AskUserQuestion, in the `/roundfix --unattended`
mode after the cycle AT follow-up. The answers this Spec rests on are,
verbatim:

- Scope: "3 Specs (Recommended)" — option text: "0338 release: tipo do PR
  segue a Spec + teto de release no _authorization.md (os dois mexem no
  release plan e no write-*); 0339 apply aprovado apaga as aposentadas (nova
  ADR, context7-cli/exa-web-search ficam); 0340 links de evidência da QA
  existem no commit. Fila nessa ordem."
- PR type: "Campo no _prd.md (Recommended)" — option text: "Frontmatter
  `type` no _prd.md, preenchido pelo write-prd/write-techspec (rota de
  correção → fix ou refactor, feature → feat). A fila usa no título; sem o
  campo continua feat. O release plan lê o tipo do commit, como hoje."
- Release ceiling (asked earlier the same day, in the Backlog Entry): "Teto de
  impacto (Recommended)" and "Próximo ciclo (Recommended)".
- Release of the cycle: "Uma no fim, v0.129.0 (Recommended)" — option text:
  "Depois das três Specs, com aviso aos adotantes; o apply que apaga
  aposentadas e o teto de release chegam juntos pela Baseline."
- Authorization: "Mesmas do AT (Recommended)" — option text: "Push e PR dos
  planos, entrega e merge pela fila, Makefile/workflows/go.mod/lint quando o
  desenho pedir, 'Consertar e seguir' em falha técnica, mudança da Baseline
  (agent-instructions, release, spec-routing, skills de autoria) que chega aos
  adotantes, release v0.129.0 no fim e aviso aos adotantes."
- Start: "Já (Recommended)".

Those answers keep the cycle AT rules: the answer of 2026-11-01 "Mesmas do AS
(Recommended)", the cycle AS answers of 2026-10-31 "Consertar e seguir
(Recommended)" (option text: "Diagnostico, emendo o plano ou crio uma Spec
bloqueadora e entrego, sem teto fixo de corretivas. Só pulo a Spec se o
conserto exigir uma decisão de produto ou arquitetura que contradiga o PRD
aprovado; aí ela sai do ciclo e a release vai com as demais, ou espera se for
a 0333.") and the Tooling answer "Makefile, .github/workflows, go.mod /
go.sum, Config de lint/formatter", and the answer of 2026-10-30 "Mesmas regras
do AR (Recommended)" (option text: "Push de branches, PR dos planos, entrega e
merge pela Delivery Queue, eu escolho o ponto de release e aviso os adotantes
por ai-memory."). Through them they keep the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos
que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário"). The Authorization answer names the
Baseline change this Spec makes: "mudança da Baseline (agent-instructions,
release, spec-routing, skills de autoria) que chega aos adotantes".

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release` and carries no `release_ceiling`: the v0.129.0
release is the operator's, after Specs 0338, 0339 and 0340 merge, as the
Release of the cycle answer states.

The governed set was measured on 2026-11-01 on the branch `docs/plan-0338`,
based on `dd3b2b62`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:` or `interface:`,
  every planning file of this Spec, its adopted references and ADR-0349.
- The 23 paths above are governed: the Roundfix skill's `SKILL.md` pair and
  four canonical references, the `write-prd`, `write-techspec` and
  `write-tasks` skill files that change, the two Baseline modules, the
  Standard TypeScript goldens and profile, the three managed guides and Setup
  Manifest, and the Spec Check's `constraints.go`.

## What is not governed

- `internal/spec/delivery_type.go`, `internal/spec/delivery_type_test.go`,
  `internal/spec/archive_record.go`, `internal/speccheck/delivery_type.go`,
  `internal/speccheck/delivery_type_test.go`,
  `internal/cli/deliver_workflow.go`,
  `internal/cli/deliver_queue_subject_test.go`.
- `internal/authorization/authorization.go`,
  `internal/authorization/release_ceiling_test.go`,
  `internal/releaseplan/model.go`, `internal/releaseplan/build.go`,
  `internal/releaseplan/classifier.go`,
  `internal/releaseplan/classifier_test.go`,
  `internal/releaseplan/release_grant.go`,
  `internal/releaseplan/release_grant_test.go`,
  `internal/cli/releaseplan_command.go`,
  `internal/cli/releaseplan_git_source.go`,
  `internal/cli/releaseplan_release_ceiling_test.go`, `internal/cli/cli.go`.
- `internal/baseline/module-versions.json`,
  `internal/baseline/release_clause_test.go`, the catalog files and the plan
  characterization goldens under `internal/baseline/testdata/`.
- The mirrors `skills/roundfix/references/*.md`,
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `docs/references/behavior-surfaces.json`,
  `CONTEXT.md`, ADR-0309 and ADR-0349.
- Every planning file of this Spec and its adopted references.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files, the skill
mirrors and the generated records. This declaration records the regeneration
that follows the approved edits and adds no source path.

The Module Version Record step writes the two modules' version lines and the
record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/spec-workflow.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/agent-instructions.md
  - docs/agents/spec-routing.md
  - docs/agents/setup-context.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-techspec/SKILL.md
  - skills/write-tasks/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

The Managed Refresh runs with `HOME` at a temporary directory, and a second
refresh reports no file change. No digest pin, golden, generated guide,
mirror or record is hand-edited. No top-level test is removed or renamed, so
the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module other than `core` and `spec-workflow`, the Source Baseline corpus,
  manifest, accounting or index, `internal/baseline/assets/retention/**`, the
  Setup Snapshots, any skill other than `roundfix`, `write-prd`,
  `write-techspec` and `write-tasks`, `docs/references/skill-coverage.json`,
  any existing Archive Record, or the `### QA settlement` section of any
  skill.
- No Task performs a release, creates a tag or edits `CHANGELOG.md`, a
  version file or a package manifest.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
