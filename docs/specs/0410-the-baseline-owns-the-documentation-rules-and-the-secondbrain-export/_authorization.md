---
status: approved
granted: 2026-11-18
action: reword four clauses of the Baseline's docs-layout rule and one of its cross-project rule, add a lesson rule to the `context-workflow` module and an export rule to the `secondbrain` module, and through the Managed Refresh state them in this repository's docs-layout, agent-instructions and Secondbrain guides and Setup Manifest; make `roundfix baseline update` keep the Standard Export List and fill each Amendment List, `roundfix doctor` report departures from the documentation rules, and the Release Plan count the export list and the lessons folder as maintenance-only; add the lesson step to the owned skill `archive-spec` and the steps and descriptions to five references of the Roundfix skill; apply the rules to this repository's export list and decision records; and describe it in the user guide, the glossary and the decision records
consuming: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/corrections.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/release.md
  - .agents/skills/roundfix/references/setup.md
  - docs/agents/agent-instructions.md
  - docs/agents/docs-layout.md
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
  - docs/references/coverage-record.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0410

The maintainer's answers are recorded whole in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The answers this record rests on are, verbatim (2026-11-18, and 2026-11-19 for the last):

- the routing of the request this Spec adopts (AskUserQuestion, about 19:55 BRT): "Spec com etapa de Ideia (Recommended)"; the Backlog Entry about the folders, "Juntar à Spec do Indexus (Recommended)"; the authorization, "Mesmas do BE (Recommended)"; the delivery order, "Corretiva → regras → Indexus (Recommended)";
- when the Idea was authored (about 22:35 BRT): "Ideia 0410 agora (Recommended)";
- the answers of the Idea stage (about 23:20 BRT): "Só a função de cada pasta (Recommended)", "Padrão + variação declarada (Recommended)", "Só as marcadas", "Agente marca, você vê no PR (Recommended)", "Apagar substituída, emendada no topo (Recommended)", "Finding fechado e archive c/ tropeço (Recommended)" and "Agente grava, você vê no PR (Recommended)";
- the approval of the Idea (about 23:55 BRT): "Aprovar e seguir ao PRD (Recommended)", with "O Agente da sessão do update (Recommended)", "Sim, mais docs/lessons/ (Recommended)" and "Sim, na lista padrão (Recommended)";
- the open questions of the PRD and the TechSpec, with the maintainer away, as the operator's brief of 2026-11-19 records it: "Seguir o padrão do autor (Recommended)".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

The request itself names the Baseline as the place of the rules: Roundfix defines, through the Baseline, "the content of documentation files, the folder layout, the required files, what is kept and what is not, and what each project exports to the Secondbrain". The AJ grant above covers the three Baseline modules, the derived Baseline files and the record this Spec declares.

For this Spec the answer grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, with up to five corrections. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: feat`. Its Release Ceiling is `minor` because every adopter receives new Normative Clauses and a file the update writes, the update gains a part of its plan and one JSON member, and the Doctor Command gains two report lines; no command, flag, configuration key or exit code changes or is removed. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge, corrections included.

## How the governed set was measured

The author measured the governed set on 2026-11-19 on the planning clone of the branch `docs/plan-0410`, based on `0ea36244`, with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository. The list of files the Baseline text and its regeneration change came from a prototype on a scratch worktree of the same commit, removed afterwards. [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md) records both. The probe covered every file the prototype changed, every file the TechSpec names, ADR-0438, ADR-0439, ADR-0440 and the controls below.

- The twenty-one paths above are the governed ones this Spec changes:
  - the modules `context-workflow`, `core` and `secondbrain`, the Standard TypeScript profile, whose formatter pin moves, and that profile's formatter goldens of the docs-layout, agent-instructions and Secondbrain guides;
  - this repository's docs-layout, agent-instructions and Secondbrain guides and its Setup Manifest, which the Managed Refresh writes;
  - the canonical and mirrored `SKILL.md` of `archive-spec` and of `roundfix`, the last one for its version fields only;
  - the canonical `baseline`, `corrections`, `modes`, `release` and `setup` references of the Roundfix skill;
  - the coverage record, which is recorded again because the two tests of the history form of the export file are replaced.
- These governed paths are not granted, and no Task may change them:
  - every other owned skill and every other file of the two skills above, except the mirrors that `make skills-sync` writes from the granted files;
  - every Vendored Skill under `.agents/skills`, `domain-modeling` and `knowledge-workspace` included;
  - every other Baseline module, `spec-workflow.json` and `frontend.json` included, every template, every other profile and formatter golden, the Source Baseline corpus and manifests under `internal/baseline/assets/source-baselines/` and every Setup Snapshot;
  - every other guide under `docs/agents/`: `specific-repository.md`, `spec-routing.md`, `skill-dispatch.md` and `autonomous-work.md`;
  - `internal/baseline/plan_test.go`, `internal/cli/cli_test.go`, `internal/cli/baseline_human_test.go` and `skills/baseline_skill_contract_test.go`; new tests go in new files;
  - the controls: the `Makefile`, `.github/workflows/ci-verify.yml`, `go.mod` and `.roundfixrc.yml`.

## What is not governed

- `internal/baseline/module-versions.json`, `internal/baseline/clause_characterization_test.go`, `internal/baseline/retirement_rule_clauses_test.go`, `internal/baseline/repository_rule_placement_test.go`, the three catalog files and the four plan goldens under `internal/baseline/testdata/`.
- The Go files of the change and their tests, under `internal/baseline`, `internal/spec`, `internal/speccheck`, `internal/cli` and `internal/releaseplan`, and `internal/docscontract/secondbrain_export_test.go`.
- The mirrors `skills/roundfix/references/*.md`, `skills/testdata/owned-skill-versions.json` and the new skill contract test.
- `.secondbrain-export`, every file under `docs/adr/` and `docs/lessons/`, `CONTEXT.md`, the user-guide pages and `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the derived Baseline files, this repository's managed guides and Setup Manifest, the skill mirrors and the derived records. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/secondbrain.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/agent-instructions.md
  - docs/agents/docs-layout.md
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/archive-spec/SKILL.md
  - skills/archive-spec/SKILL.md
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

No mirror, version field, digest, golden or record is hand-edited. Once the documentation part of the update exists, the Managed Refresh above also keeps this repository's `.secondbrain-export` and fills an Amendment List; neither file is governed.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`, except for a correction that requires one;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot other than the three modules and the Standard TypeScript profile's pin and three golden guides;
  - any guide under `docs/agents/` other than the three named and `setup-context.json`, and this repository's `AGENTS.md`;
  - any check, code or output of `roundfix spec check`, and any field of the Archive Record or of the Setup Manifest;
  - any existing flag, configuration key or exit code of `roundfix baseline update`, `roundfix doctor`, `roundfix update check` or `roundfix release plan`;
  - `docs/references/skill-coverage.json`, unless a Task records a Coverage Review there;
  - any other skill, any Vendored Skill and any Archive Record;
  - the `### QA settlement` section of any skill;
  - the decision text of any accepted ADR: an existing ADR gains only header fields and, for ADR-0248, ADR-0278 and ADR-0412, a dated notice.
- No change to any file of another repository, and no read of one: the Secondbrain's template, sync and mirrors stay as they are.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`, `rsync` against another repository, `codex` or `claude`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
  - sends a request to a provider, the judge or a reviewer;
  - reads a file under the Spec root;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
