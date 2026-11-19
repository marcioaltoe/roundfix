---
status: approved
granted: 2026-11-19
action: add the Authoring Question Preference to Project Config and report it, with the Delegated Decisions of the active Specs, in the status view of the Follow-up Command; teach the Clarification Record and its check the Delegated Decision; add to the Baseline's Spec workflow module a rule that says what a session does with an open authoring question when nobody is present, and through the Managed Refresh state it in this repository's Spec routing guide and Setup Manifest; add the steps to the Roundfix skill's modes, follow-up and spec references and to the owned skills write-prd, write-techspec and write-tasks; and describe it in the user guide, the glossary and the decision records
consuming: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/follow-up.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/clarification-record.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
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
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0412

The maintainer's request and answers are recorded whole in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md). The answers this record rests on are, verbatim (2026-11-19):

- the request, in his own words (about 00:15 BRT): "Antes de iniciar o fluxo, perguntar se pode seguir o padrão de autorização decidido pelo autor ou se deve esperar; seria bacana no projeto ter uma configuração no .roundfixrc.yml com a preferencia por projeto e o default ser a spec para e espera o autor.";
- how the request is recorded (about 00:20 BRT): "Spec 0412 com Ideia agora", whose option text reads "abre a Spec já, com etapa de Ideia, entregue depois da 0410";
- the open questions of this cycle (the same call): "Seguir o padrão do autor (Recommended)";
- the Idea stage's four answers (about 01:10 BRT): "Duas Specs (Recommended)", "PRD, TechSpec e Tasks (Recommended)", "Config responde; sem ela, pergunta (Recommended)" and "Até o fim, com relato (Recommended)";
- the approval of the Idea (about 02:15 BRT): "Aprovar e seguir ao PRD (Recommended)", "O que a sessão foi informada (Recommended)", "Sem teto, com contagem (Recommended)" and "Enquanto a Spec estiver ativa (Recommended)".
- the authorization of Specs 0412 and 0413 (later the same day): "Mesmas do BE para as duas (Recommended)", whose option text reads "implement, commit, push, PR, merge e release minor depois de cada uma; até 5 correções; Makefile/CI/go.mod/lint só se a correção exigir".

## Where the grant of the operations comes from

The grant is the maintainer's answer of 2026-11-19 to the authorization of Specs 0412 and 0413: "Mesmas do BE para as duas (Recommended)", with the option text "implement, commit, push, PR, merge e release minor depois de cada uma; até 5 correções; Makefile/CI/go.mod/lint só se a correção exigir". It grants this Spec `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, a `minor` Release Ceiling and up to five corrections, the same authorization "Mesmas do BE (Recommended)" that he gave on 2026-11-18 for the Specs then planned in this cycle.

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: feat`. Its Release Ceiling is `minor` because it adds one Project Config key, two fields and two blocks of lines to the status view, and one rule that every adopter of the Spec workflow receives with its next Baseline update; no command, flag or exit code is removed or renamed. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge, corrections included.

## How the governed set was measured

The author measured the governed set on 2026-11-19 on the planning clone of the branch `docs/plan-0412`, based on `894710c3`, with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository. The probe covered every file the TechSpec names, the ones a regeneration command writes included, and the controls below. The files a Baseline module edit regenerates are the ones the prototype of a new rule in the same module measured for the Spec delivered just before this one, and the prototype of this Spec's own rule, recorded in [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md), confirmed them.

- The nineteen paths above are the governed ones this Spec changes:
  - the Spec workflow module, the Standard TypeScript profile, whose formatter pin moves, and that profile's formatter golden of the Spec routing guide;
  - this repository's Spec routing guide and Setup Manifest, which the Managed Refresh writes;
  - the canonical and mirrored `SKILL.md` of `write-prd`, `write-techspec`, `write-tasks` and `roundfix`, the last one for its version fields only;
  - the canonical `modes`, `follow-up` and `spec` references of the Roundfix skill, the canonical Clarification Record reference of `write-prd`, and the canonical and mirrored PRD template.
- These governed paths are not granted, and no Task may change them:
  - `.roundfixrc.yml`: this Spec adds the key to the product and does not declare it for this repository;
  - every other owned skill and every other file of the four skills above, `write-idea` and the TechSpec template in both copies included, except the mirrors that `make skills-sync` writes from the granted files;
  - every Vendored Skill under `.agents/skills`;
  - every other Baseline module, every template, every other profile and formatter golden, the Source Baseline corpus and manifests under `internal/baseline/assets/source-baselines/` and every Setup Snapshot;
  - every guide under `docs/agents/` but `spec-routing.md` and `setup-context.json`: `agent-instructions.md`, `autonomous-work.md`, `docs-layout.md`, `specific-repository.md`, `skill-dispatch.md` and `secondbrain.md`;
  - the detector table `internal/speccheck/coherence.go`, the corpus pins `internal/docscontract/testdata/corpus-golden.json` and `internal/spec/archive_layout_characterization_test.go`, and `docs/references/coverage-record.json`: the design adds no check code and removes or renames no test;
  - the controls: the `Makefile`, `.github/workflows/ci-verify.yml` and `go.mod`.

## What is not governed

- Under `internal/config/`, `internal/cli/` and `internal/speccheck/`: the files the TechSpec names, `internal/speccheck/clarification_record.go` included, and their tests.
- `internal/baseline/module-versions.json`, `internal/baseline/clause_characterization_test.go`, the new clause test, the three catalog files and the plan goldens under `internal/baseline/testdata/`.
- The mirrors `skills/roundfix/references/modes.md`, `skills/roundfix/references/follow-up.md`, `skills/roundfix/references/spec.md` and `skills/write-prd/references/clarification-record.md`, `skills/testdata/owned-skill-versions.json` and the new skill test.
- `CONTEXT.md`, the user guide pages `configuration.md`, `commands/follow-up.md`, `commands/spec.md` and `context-driven-development.md`, `docs/references/behavior-surfaces.json`, ADR-0442 and ADR-0280.
- `.secondbrain-export`, which the Managed Refresh of the last implementation Task may rewrite between its markers on a tree that holds Spec 0410, and the files the Tasks stage added: `internal/cli/follow_up_status_delegated.go` and its test and `skills/delegated_decision_record_test.go`, measured not governed with the same probe on 2026-11-19.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the derived Baseline files, this repository's managed guide and Setup Manifest, the skill mirrors and the derived records. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/spec-workflow.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/spec-routing.md
  - docs/agents/setup-context.json
  - docs/adr/0280-a-guide-task-records-its-own-surface-fingerprint-a-prd-records-who-decided-and-the-glossary-is-held-to-its-own-rule.md
  - .secondbrain-export
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/write-techspec/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest, golden or record is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`, except for a correction that requires one;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot other than `spec-workflow.json` and the Standard TypeScript profile's pin and golden Spec routing guide;
  - any guide under `docs/agents/` other than `spec-routing.md` and `setup-context.json`, and this repository's `AGENTS.md`;
  - the stop-and-ask clause about questions, the clauses of the Rule Challenge and the clauses of the Idea stage, which keep their ids, force and text;
  - any command, flag or exit code, and any check code of `roundfix spec check`;
  - `docs/references/coverage-record.json` and `docs/references/skill-coverage.json`;
  - any other skill, `write-idea` included, any Vendored Skill and any Archive Record;
  - the `### QA settlement` section of any skill;
  - any accepted ADR other than ADR-0442, which this Spec adds, and ADR-0280, which gains its amendment notice.
- No change to any file of another repository, and no read of one.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`, `codex` or `claude`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
  - sends a request to a provider, the judge or a reviewer;
  - reads a file under the Spec root;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
