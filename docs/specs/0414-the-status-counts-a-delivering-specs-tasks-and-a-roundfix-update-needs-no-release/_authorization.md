---
status: approved
granted: 2026-11-19
action: make the status view of the Follow-up Command count a Delivery Queue item's Tasks on its Item Branch, and make the Release Plan count the files Roundfix manages in an adopter and every CONTEXT.md as maintenance-only, with the guides, the release reference of the Roundfix skill, the glossary and ADR-0447
consuming: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/release.md
  - skills/roundfix/SKILL.md
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

# Approved authority for Spec 0414

The maintainer's answers are recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md). The answer this record rests on is, verbatim: "Spec fix 0414 logo após a 0411 (Recommended)" — uma Spec com as duas correções (como a 0407), TechSpec no Opus, sem Ideia, "Mesmas do BE", release patch, entregue antes da 0409. The boundary of the release correction is the maintainer's answer "Arquivos geridos pelo Roundfix + CONTEXT.md (Recommended)", quoted in the PRD `## Decisions`.

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

For this Spec the answer grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, with up to five corrections. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: fix` and its Release Ceiling is `patch`: it adds no command, flag, field or exit code. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

## How the governed set was measured

The author measured the governed set on 2026-11-19 on the planning clone of the branch `docs/plan-0414`, based on `145d6c76`, with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository. The probe covered every file the Tasks declare, this Spec's planning files, ADR-0447 and six controls: the `Makefile`, `.github/workflows/ci.yml`, `go.mod`, `.roundfixrc.yml`, `skills-lock.json` and the `qa-gate` skill's `SKILL.md`.

- The three paths above are the governed ones this Spec changes: the canonical `SKILL.md` and `references/release.md` of the Roundfix skill, and the mirror `SKILL.md`.
- These governed paths are not granted, and no Task may change them:
  - every owned skill other than the Roundfix skill, canonical or mirror, and every file of the Roundfix skill other than the three above and the mirror `skills/roundfix/references/release.md`;
  - every Vendored Skill under `.agents/skills`;
  - every Baseline module, template, profile, Source Baseline, Setup Snapshot and golden under `internal/baseline/assets/`;
  - the guides under `docs/agents/`;
  - the six controls.

## What is not governed

- `internal/spec/spec.go`, the new `internal/spec/load_at_test.go`, `internal/cli/follow_up_queue.go`, `internal/cli/follow_up_status.go`, `internal/cli/follow_up_queue_test.go` and `internal/cli/follow_up_status_test.go`.
- `internal/releaseplan/classifier.go`, `internal/releaseplan/build.go`, `internal/releaseplan/model.go`, the new `internal/releaseplan/managed.go` and `internal/releaseplan/managed_test.go`, `internal/cli/cli.go`, `internal/cli/releaseplan_command_test.go`, the new `internal/cli/releaseplan_managed_paths.go` and `internal/cli/releaseplan_managed_paths_test.go`.
- `skills/roundfix/references/release.md`, the new `skills/release_adopter_maintenance_test.go` and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/follow-up.md`, `docs/user-guide/release-runbook.md`, `docs/user-guide/usage.md`, `CONTEXT.md`, `docs/references/skill-coverage.json`, `docs/references/behavior-surfaces.json` and ADR-0447.
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
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
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
  - `.roundfixrc.yml` or `skills-lock.json` of this repository;
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
