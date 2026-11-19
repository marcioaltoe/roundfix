---
status: approved
granted: 2026-11-18
action: state in the qa-gate and write-tasks skills that a fresh home keeps the host toolchains, and make the status view of the Follow-up Command show a renamed Spec once, with the user guide, the glossary and the decision records
consuming: 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/write-tasks/SKILL.md
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

# Approved authority for Spec 0407

The maintainer's answers are recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The answer this record rests on is, verbatim, from the adopted Backlog Entry's `## Decisions`: "2026-11-18: \"Spec após a release\"; authorization \"Mesmas do BE (Recommended)\" (to be recorded in the Spec's `_authorization.md`)." The order is the maintainer's answer of 2026-11-18 about 19:55 BRT: "Ordem: \"Corretiva → regras → Indexus (Recommended)\"; o fix da fila entra logo após a corretiva." This Spec is the corrective one and delivers first.

The author measured that the two defects share no component, and on 2026-11-18 about 22:35 BRT the maintainer kept them together: "Uma Spec só (Recommended)" — one Spec, one patch release.

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

For this Spec the answer grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, with up to five corrections. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: fix` and its Release Ceiling is `patch`: it adds no command, flag, field or output line. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

## How the governed set was measured

The author measured the governed set on 2026-11-18 on the planning clone of the branch `docs/plan-0407`, based on `1a8e08b4`, and measured it again after rebasing on `d1c43f9e`, with the same result, with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository. The probe covered every file the Tasks declare, this Spec's planning files, ADR-0410, ADR-0434, the adopted Backlog Entry and four controls: the `Makefile`, `.github/workflows/ci.yml`, `go.mod` and the Roundfix skill's `SKILL.md`.

- The four paths above are the governed ones this Spec changes: the `SKILL.md` of the `qa-gate` and `write-tasks` skills and their mirrors.
- These governed paths are not granted, and no Task may change them:
  - the Roundfix skill and every other owned skill, canonical or mirror;
  - every Vendored Skill under `.agents/skills`;
  - every Baseline module, template, profile, Source Baseline, Setup Snapshot and golden under `internal/baseline/assets/`;
  - the guides under `docs/agents/`;
  - the four controls.

## What is not governed

- `internal/cli/follow_up_queue.go`, `internal/cli/follow_up_queue_test.go` and the new `skills/qa_fresh_home_test.go`.
- `docs/user-guide/commands/follow-up.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, `skills/testdata/owned-skill-versions.json`, ADR-0410 and ADR-0434.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the skill mirrors and the derived records. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/write-tasks/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field or record is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Baseline source or derived Baseline file, and any guide under `docs/agents/`;
  - `docs/references/coverage-record.json`;
  - any Archive Record;
  - the `### QA settlement` section of any skill.
- No test or Task Verification:
  - reads another repository, except a temporary one the test itself created;
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
