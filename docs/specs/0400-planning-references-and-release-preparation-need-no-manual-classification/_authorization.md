---
status: approved
granted: 2026-11-16
action: count docs/references as planning evidence and recognise a release-preparation commit from the release_paths of the Specs a release plan range delivers, report it in the JSON change evidence, and describe both rules and both publication triggers in the Roundfix skill's release reference, the release plan usage text, the release runbook, the usage page and the glossary
consuming: 0400-planning-references-and-release-preparation-need-no-manual-classification
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
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0400

The Spec adopts one Backlog Entry whole and the release-guide half of another, recorded by the operator from the Oraculum's and the Quaestor's messages of 2026-11-16. The maintainer's answers are recorded whole in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md). The answers this record rests on are, verbatim (AskUserQuestion, 2026-11-16 ~20:55 BRT):

- Achado Oraculum 01a1e695 → "Spec fix depois da 0399 (Recommended)": Spec de release plan que adota o achado do Oraculum + parte 2 do backlog do Quaestor (release.md supõe release por tag); parte 1 do Quaestor (fila atualizar branch atrasada) fica no backlog.
- Autorização dessa Spec → "Mesmas do BE (Recommended)".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it. These include the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas as skills se necessário") and the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`. The Spec is a `type: fix`. Its Release Ceiling is `minor`, not `patch`, because the design adds a public output: the `releasePreparation` field of each change in the `roundfix.release-plan/0.0.1` JSON result. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge, up to three corrections included.

## How the governed set was measured

The author measured the governed set on 2026-11-16 in a clone on the branch `docs/plan-0400`, based on `c2d1f77c`, with `speccheck.GovernedPath`.

- The measurement used a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository. The probe file never existed in the tree.
- It covered 39 paths: every source, test, guide, skill, record and document file the TechSpec names or expects the Tasks to create, this Spec's planning files, ADR-0412, the revised Backlog Entry, and three controls.
- The three paths above are the governed ones this Spec's design changes: the Roundfix skill's canonical release reference and its `SKILL.md` pair.
- These governed paths are not granted, and no Task may change them:
  - `internal/cli/cli_test.go` and `internal/docscontract/publicdocs_test.go`; new tests go in new files;
  - the three controls: the `Makefile`, `.github/workflows/ci-verify.yml` and `go.mod`.

`write-tasks` measures the exact derived set on a scratch copy before the Run. If a sanctioned regeneration rewrites a governed path outside `paths`, the operator amends this record before the Run, outside the change it authorizes.

## What is not governed

- Every Go source and test file under `internal/releaseplan` and `internal/cli` that the Tasks declare, including `internal/cli/cli.go`.
- The mirror `skills/roundfix/references/release.md`, `skills/testdata/owned-skill-versions.json` and the tests under `skills/`.
- `docs/user-guide/release-runbook.md`, `docs/user-guide/usage.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json` and ADR-0412.
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

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot, and any file under `docs/agents/`;
  - any other skill, Vendored Skill or Archive Record;
  - `.github/workflows/release.yml` or the publication steps of `docs/user-guide/release-runbook.md`;
  - the `### QA settlement` section of any skill.
- No change to any file of another repository, and no read of one.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - creates, moves or deletes a tag, a ref or a remote of this repository;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- The QA gate may run `roundfix release plan` read-only over this repository's existing tags in a clean worktree it creates and removes.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
