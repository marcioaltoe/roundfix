---
status: approved
granted: 2026-11-20
action: reword the two clauses of the Baseline's Secondbrain module that name qmd so that the Secondbrain search and the advisory duplicate check use the Indexus MCP tool and no other search tool; through the Managed Refresh and the sanctioned regeneration, state it in this repository's Secondbrain guide, its Setup Manifest and the derived Baseline files; reword the same check in the owned skill write-idea, in both copies; and record the decision in a new ADR and in an amendment notice of ADR-0330
consuming: 0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd
paths:
  - .agents/skills/write-idea/SKILL.md
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/write-idea/SKILL.md
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

# Approved authority for Spec 0416

The messages and the maintainer's decisions this record rests on are copied whole in [`references/2026-11-20-messages-and-decisions.md`](references/2026-11-20-messages-and-decisions.md).

## Where the grant comes from

Two decisions of the maintainer of 2026-11-20, each delivered as a `DECISÃO DE MARCIO` block to this project's ai-memory inbox by the maintainer's assistant.

- The change itself, ai-memory message `01a1f952-79e3-7c81-8459-626be6c3e375`, decided 11:33 BRT. The maintainer's answer, verbatim: "Agora, os testes devem ser com o objetivo inicial de ser superiores ao que o qmd/secondbrain entrega no que os projetos/jarvis usam. Se for, já podemos usar somente o indexus e eliminar o qmd no mini e macbook." Its recorded scope names this repository's part: "trocar toda referência ao qmd por Indexus, sem reserva: guia docs/agents/secondbrain.md e write-idea do baseline (Roundfix)", authorized "sem novo pedido" once the Indexus is shown superior. The author of the PRD read this block in the inbox.
- Preparing it now, ai-memory message `01a1f955-c546-7ae2-b4ed-786f9bc9be79`, decided 11:44 BRT. The operator's brief records it as: answer, option 1 "Pular a Ideia"; scope, "this Spec — PRD, TechSpec, Tasks, implementation in a PR; the release waits for the Secondbrain's notice that the Indexus is superior to qmd (DECISÃO 01a1f952)". The block's `resposta` reads: opção 1, "Pular a Ideia"; its `alcance`: "a Spec do Roundfix que troca o qmd pelo Indexus no módulo secondbrain do baseline (docs/agents/secondbrain.md) e no write-idea: PRD, TechSpec, Tasks, implementação num PR. O release espera o aviso do Secondbrain de que o Indexus é superior ao qmd (DECISÃO 01a1f952)." The session that popped the block copied it whole into section 5 of the reference file.

The first decision names the two things this Spec changes, the Secondbrain guide of the Baseline and `write-idea`. The second lets the work reach an open Pull Request before the test's result. Together they grant `implement`, `commit`, `push` and `pull_request`.

## The widening of 2026-11-20

The Indexus test ended and the maintainer decided the replacement goes ahead now, without waiting for the wiki's access list. ai-memory message `01a1f993-7830-7362-b8dc-49813553eab8`, popped by the Supervisor session at 12:29 BRT (15:29 UTC), copied whole in section 10 of the reference file. Its answer, verbatim: "Superior; remover já". Its scope names this repository's part: "Roundfix faz o merge da 0416 e solta a versão de correção (release patch) com o guia e o write-idea só com Indexus". It also records that the maintainer accepts that, until the wiki's classification ends, project Agents go without the wiki, which lives in the personal acervo.

This change, on the default branch before the consuming delivery, widens the grant to `merge` and `release` at a `patch` Release Ceiling, with the three version-only release paths. The governed paths, the sanctioned regeneration and the limits below are unchanged.

## How the governed set was measured

The author measured it on 2026-11-20 at `851024fa` with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` whose test file lived outside the repository, over the sixteen files a prototype on a scratch clone changed and over the candidates below. [`references/2026-11-20-characterization.md`](references/2026-11-20-characterization.md), sections 3 and 4, records the prototype and the probe.

- The seven paths above are the governed ones this Spec changes:
  - the Secondbrain module of the Baseline;
  - the Standard TypeScript profile, whose formatter pin moves, and that profile's formatter golden of the Secondbrain guide;
  - this repository's Secondbrain guide and Setup Manifest, which the Managed Refresh writes;
  - the canonical and mirrored `SKILL.md` of `write-idea`.
- These governed paths are not granted, and no Task may change them:
  - the Source Baseline corpus, manifests and index under `internal/baseline/assets/source-baselines/`, which keep the wording they recorded;
  - every Vendored Skill under `.agents/skills`, `knowledge-workspace` included, and `skills-lock.json`'s entries;
  - every other Baseline module, template, profile, formatter golden and Setup Snapshot;
  - every other owned skill and every other file of `write-idea`;
  - every guide under `docs/agents/` but `secondbrain.md` and `setup-context.json`, `specific-repository.md` included;
  - `.roundfixrc.yml`, the `Makefile`, `.github/workflows` and `go.mod`.

The Tasks author measures again every file the Tasks declare and, when the set differs, amends this list before a Run starts.

## What is not governed

- `internal/baseline/module-versions.json`, the three catalog files and the four plan goldens under `internal/baseline/testdata/`.
- `skills/testdata/owned-skill-versions.json`.
- ADR-0330, the new ADR, `.secondbrain-export` and `CONTEXT.md`.
- `internal/baseline/clause_characterization_test.go` and a new test file under `internal/baseline/`, `skills/` or `internal/docscontract/`.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the module's version line, the derived Baseline files, this repository's managed guide and Setup Manifest, the skill mirror and the skill version record. This declaration records the regeneration that follows the approved edits and adds no source path. The commands are the ones the prototype ran; the Tasks author confirms the outputs on the tree the Tasks run on.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/write-idea/SKILL.md
  - skills/write-idea/SKILL.md
```

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/secondbrain.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
  - docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md
  - .secondbrain-export
```

No mirror, version field, digest, golden or record is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot other than `secondbrain.json` and the Standard TypeScript profile's pin and golden Secondbrain guide;
  - any clause of the Secondbrain module other than the two that name qmd, which keep their ids and their force;
  - any command, flag, exit code or configuration key;
  - any other skill, any Vendored Skill and any Archive Record;
  - the `### QA settlement` section of any skill;
  - any accepted ADR other than the one this Spec adds and ADR-0330, which gains its amendment notice.
- No change to any file of another repository, the Secondbrain included, and no removal of qmd from any machine.
- No test or Task Verification:
  - reaches the network, calls the Indexus, ai-memory or Exa, or runs `qmd`;
  - reads the Secondbrain or any file outside the repository;
  - starts an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
  - reads a file under the Spec root.
- No Task sends an ai-memory message or a Spec Event Announcement; the operator does.
