---
status: approved
granted: 2026-11-18
action: make profiles configure and profiles check --apply prove under the configured profile proof setup deadline, and give the Research Subagent of the Recommended Authoring Declaration the built-in web tools with a report for a research entry that declares none, with the user guides, the Roundfix skill, the glossary and the decision records
consuming: 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/profiles.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/modes.md
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

# Approved authority for Spec 0411

The maintainer's answers are recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The answers this record rests on are, verbatim:

- to "Defeito confirmado ... Como tratamos?": "Spec fix 0411 após a 0408 (Recommended)" — TechSpec no Opus, sem Ideia, "Mesmas do BE", release patch; entra na fila logo depois da 0408 e antes da 0409;
- to "... Onde entra?" about the report of conexus: "Juntar à Spec 0411 (Recommended)" — dois defeitos que os adotantes acharam ao adotar as recomendações da 0.148.0; uma Spec só, mesmo release patch, logo após a 0408.

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

For this Spec the answer grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, with up to five corrections. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: fix` and its Release Ceiling is `patch`, by the maintainer's answer. It adds one optional JSON member and one report line that name what the fix recommends, and no command, flag or exit code. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

## How the governed set was measured

The author measured the governed set on 2026-11-18 on the planning clone of the branch `docs/plan-0411`, based on `a6a4a15b`, with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository. The probe covered every file the Tasks declare, ADR-0415, ADR-0441 and four controls: the `Makefile`, `.github/workflows/ci.yml`, `go.mod` and `.roundfixrc.yml`.

- The five paths above are the governed ones this Spec changes: the Roundfix skill's `SKILL.md` in both copies, whose version fields the record command raises, and the canonical `profiles.md`, `setup.md` and `modes.md` references of the Roundfix skill.
- These governed paths are not granted, and no Task may change them:
  - every other file of the Roundfix skill and every other owned skill, canonical or mirror, except the mirrors that `make skills-sync` writes from the granted references;
  - every Vendored Skill under `.agents/skills`;
  - every Baseline module, template, profile, Source Baseline, Setup Snapshot and golden under `internal/baseline/assets/`;
  - the guides under `docs/agents/`;
  - the four controls.

## What is not governed

- `internal/config/proof_setup_deadline.go`, `internal/config/proof_setup_deadline_test.go`, `internal/config/authoring_recommendation.go`, `internal/config/authoring_recommendation_test.go`.
- `internal/cli/profiles_configure.go`, `internal/cli/profile_proof_setup_deadline_test.go`, `internal/cli/authoring_reports.go`, `internal/cli/authoring_research_tools_test.go`, `internal/cli/authoring_recommendation_test.go`.
- `internal/docscontract/authoring_recommendation_test.go`, `skills/authoring_test.go`.
- The mirrors `skills/roundfix/references/profiles.md`, `skills/roundfix/references/setup.md` and `skills/roundfix/references/modes.md`.
- `docs/user-guide/configuration.md`, `docs/user-guide/commands/profiles.md`, `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/update.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `skills/testdata/owned-skill-versions.json`, ADR-0415 and ADR-0441.
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
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Baseline source or derived Baseline file, and any guide under `docs/agents/`;
  - `docs/references/coverage-record.json` and `docs/references/skill-coverage.json`;
  - `skills/authoring.go`, `internal/config/authoring.go`, `internal/agent/fallback.go` and `internal/baselineacp/`;
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
