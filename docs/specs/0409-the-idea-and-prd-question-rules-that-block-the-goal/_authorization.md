---
status: approved
granted: 2026-11-18
action: add to the Baseline's Spec workflow module a rule with two Normative Clauses under which an author of brainstorming, an Idea, a PRD or a TechSpec asks the user whether a rule that blocks the goal is kept, changed or retired, and through the Managed Refresh state them in this repository's Spec routing guide and Setup Manifest; add the step to the owned skills brainstorming, write-idea, write-prd and write-techspec, the record section to the Idea and PRD templates and the waiting rule to the Roundfix skill's modes reference; and describe it in the user guide, the glossary and the decision records
consuming: 0409-the-idea-and-prd-question-rules-that-block-the-goal
paths:
  - .agents/skills/brainstorming/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-idea/references/idea-template.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-techspec/SKILL.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/brainstorming/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
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

# Approved authority for Spec 0409

The maintainer's answers are recorded whole in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The answers this record rests on are, verbatim (2026-11-18):

- the routing of the Backlog Entry this Spec adopts: when, "Spec após a release"; authorization, "Mesmas do BE (Recommended)";
- the Idea stage (AskUserQuestion, about 19:55 BRT): "Ideia agora (Recommended)", for the Spec that the question names as authorized "Mesmas do BE";
- the delivery order (the same call): "Corretiva → regras → Indexus (Recommended)";
- the Idea stage's six answers (AskUserQuestion, about 20:30 BRT): "Cláusula, passo e registro (Recommended)", "Entra, só ao contradizer ADR (Recommended)", "A própria Spec (Recommended)", "A Spec espera (Recommended)", "Sim, vira proposta ao Roundfix (Recommended)" and "Na própria regra (Recommended)";
- the approval of the Idea (AskUserQuestion, about 21:05 BRT): "Aprovar e seguir ao PRD (Recommended)", "Três: manter, mudar, retirar (Recommended)" and "Segue dentro da regra (Recommended)".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

The maintainer's request itself names the Baseline as the place of the rule: "Isso tem que existir também no baseline." The AJ grant above covers the Baseline module and the derived Baseline files this Spec declares.

For this Spec the answer grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, with up to five corrections. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: feat`. Its Release Ceiling is `minor` because every adopter of the Spec workflow receives two new Normative Clauses with its next Baseline update; no command, flag, configuration key or exit code changes. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge, corrections included.

## How the governed set was measured

The author measured the governed set on 2026-11-19 on the planning clone of the branch `docs/plan-0409`, based on `348962d2`, with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository. The file list came from a prototype of the whole change on a scratch worktree of the same commit, where the edits and the six regeneration commands below had run; the worktree was removed afterwards. [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md) records both. The probe covered every file the prototype changed, the two new test files the TechSpec names, ADR-0436, ADR-0437 and the controls below.

- The nineteen paths above are the governed ones this Spec changes:
  - the Spec workflow module, the Standard TypeScript profile, whose formatter pin moves, and that profile's formatter golden of the Spec routing guide;
  - this repository's Spec routing guide and Setup Manifest, which the Managed Refresh writes;
  - the canonical and mirrored `SKILL.md` of `brainstorming`, `write-idea`, `write-prd`, `write-techspec` and `roundfix`, the last one for its version fields only;
  - the canonical Idea template, the canonical and mirrored PRD template, and the canonical `modes` reference of the Roundfix skill.
- These governed paths are not granted, and no Task may change them:
  - every other owned skill and every other file of the five skills above, the TechSpec template in both copies included, except the mirrors that `make skills-sync` writes from the granted files;
  - every Vendored Skill under `.agents/skills`;
  - every other Baseline module, `core.json` and `context-workflow.json` included, every template, every other profile and formatter golden, the Source Baseline corpus and manifests under `internal/baseline/assets/source-baselines/` and every Setup Snapshot;
  - every guide under `docs/agents/` but `spec-routing.md` and `setup-context.json`: `agent-instructions.md`, `docs-layout.md`, `specific-repository.md`, `skill-dispatch.md`, `secondbrain.md` and `autonomous-work.md`;
  - `internal/baseline/plan_test.go` and `skills/baseline_skill_contract_test.go`; new tests go in new files;
  - the controls: the `Makefile`, `.github/workflows/ci-verify.yml`, `go.mod` and `.roundfixrc.yml`.

## What is not governed

- `internal/baseline/module-versions.json`, `internal/baseline/clause_characterization_test.go`, the new `internal/baseline/rule_challenge_clauses_test.go`, the three catalog files and the four plan goldens under `internal/baseline/testdata/`.
- The mirrors `skills/write-idea/references/idea-template.md` and `skills/roundfix/references/modes.md`, `skills/testdata/owned-skill-versions.json` and the new `skills/rule_challenge_authoring_test.go`.
- `CONTEXT.md`, `docs/user-guide/context-driven-development.md`, `docs/references/behavior-surfaces.json`, ADR-0436 and ADR-0437.
- `docs/agents/domain.md`, which no Task changes.
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
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/brainstorming/SKILL.md
  - skills/brainstorming/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
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
  - any command, flag, configuration key, exit code or check of `roundfix spec check`;
  - `docs/references/coverage-record.json` and `docs/references/skill-coverage.json`;
  - any other skill, any Vendored Skill and any Archive Record;
  - the `### QA settlement` section of any skill;
  - any accepted ADR other than the two this Spec adds.
- No change to any file of another repository, and no read of one.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`, `codex` or `claude`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
  - sends a request to a provider, the judge or a reviewer;
  - reads a file under the Spec root;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
