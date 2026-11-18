---
status: approved
granted: 2026-11-18
action: make CLAUDE.md a symbolic link to AGENTS.md in every Baseline plan where no content is lost, report every other shape in roundfix baseline update and roundfix doctor, state the rule in the agent-instructions guide and refresh this repository's managed guide, and describe it in the Roundfix and setup-context-driven skills, the user guide and the glossary
consuming: 0406-claude-md-is-always-a-link-to-agents-md
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/setup-context-driven/SKILL.md
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/templates/guides/agent-instructions.md
  - internal/baseline/assets/templates/index.json
  - skills/roundfix/SKILL.md
  - skills/setup-context-driven/SKILL.md
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

# Approved authority for Spec 0406

The maintainer's decision is recorded whole in [`references/2026-11-18-maintainer-decision.md`](references/2026-11-18-maintainer-decision.md). The answer this record rests on is, verbatim (Supervisor session, 2026-11-18 about 06:30 BRT): authorization for Spec 0406 → "Mesmas do BE (Recommended)".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

For this Spec the answer grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, in the single release at the end of the queue, with up to five corrections. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: fix`. Its Release Ceiling is `minor`, not `patch`, because the design adds public output: the `baseline-claude-link` line of `roundfix doctor` and two warnings of the Baseline plans. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

## How the governed set was measured

The author measured the governed set on 2026-11-18 on a scratch worktree of the branch `docs/plan-0406`, based on `eb888135`, with `speccheck.GovernedPath`, removed afterwards.

- The measurement used a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository. The probe file never existed in the tree.
- It covered 58 paths: every file the prototype's sanctioned regeneration rewrote, every source, test, guide, skill, record and document file the TechSpec names or expects the Tasks to create, this Spec's planning files, ADR-0432, and three controls.
- The 14 paths above are the governed ones this Spec's design changes:
  - the `core` module, the `agent-instructions.md` guide template and `templates/index.json`;
  - the Standard TypeScript profile's golden pin and the formatter golden it renders from the changed guide;
  - this repository's managed `docs/agents/agent-instructions.md` and Setup Manifest;
  - the Roundfix skill's `SKILL.md` pair and its canonical `baseline`, `modes` and `setup` references;
  - the `setup-context-driven` skill's `SKILL.md` pair.
- These governed paths are not granted, and no Task may change them:
  - the Source Baseline corpus and manifests under `internal/baseline/assets/source-baselines/`;
  - every Setup Snapshot under `internal/baseline/assets/setups/`;
  - every profile other than `standard-typescript-monorepo.json`;
  - the formatter golden `AGENTS.md`;
  - `docs/agents/specific-repository.md`;
  - `internal/baseline/plan_test.go`, `internal/cli/cli_test.go` and `internal/docscontract/publicdocs_test.go`;
  - the three controls: the `Makefile`, `.github/workflows/ci.yml` and `go.mod`.

Spec 0405 regenerates some of the same derived files before this Spec runs. If, on the tree the Run starts from, the sanctioned regeneration rewrites a governed path outside `paths`, the operator amends this record before the Run, outside the change it authorizes.

## What is not governed

- Every Go source, test and fixture file under `internal/baseline` and `internal/cli` that the Tasks declare, including `internal/baseline/module-versions.json` and the catalog records and plan goldens under `internal/baseline/testdata/`.
- This repository's `AGENTS.md` and `CLAUDE.md`, which this Spec does not change.
- The mirrors under `skills/roundfix/references/`, and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/baseline.md`, `docs/user-guide/commands/doctor.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, ADR-0100, ADR-0411, ADR-0432 and ADR-0433.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the derived records, the skill mirrors and this repository's managed guide. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/core.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/agent-instructions.md
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
  - .agents/skills/setup-context-driven/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/setup-context-driven/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest, golden or record is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Vendored Skill under `.agents/skills`;
  - any Source Baseline, Setup Snapshot, or profile other than `standard-typescript-monorepo.json`;
  - `docs/agents/specific-repository.md`;
  - `docs/references/coverage-record.json`;
  - any Archive Record;
  - the `### QA settlement` section of any skill.
- No change to any file of another repository. The QA gate may read the adopted checkouts the inventory names, and may run `roundfix baseline update` there without confirmation, which writes nothing. It confirms each `git status --porcelain` unchanged afterwards. Never with `--yes`, `--confirm-plan` or `--promote`.
- No test or Task Verification:
  - reads another repository, except a temporary one the test itself created;
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
