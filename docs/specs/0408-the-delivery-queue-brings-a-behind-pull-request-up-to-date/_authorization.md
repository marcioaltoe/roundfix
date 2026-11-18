---
status: approved
granted: 2026-11-18
action: make the Delivery Queue bring a delivery Pull Request that GitHub reports BEHIND up to date with a Base Update instead of parking it, bound the Base Updates of an item, and describe the behavior in the Roundfix skill's deliver reference, the user guide and the glossary
consuming: 0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
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

# Approved authority for Spec 0408

The maintainer answered on 2026-11-18 about 19:55 BRT, through `/roundfix --follow-up` (AskUserQuestion), recorded by the operator in the Cycle BG brief. Verbatim:

- Backlog 2026-11-16 Delivery Queue com branch estrita (Quaestor): "Spec fix após a corretiva (Recommended)" — "Mesmas do BE".
- Ordem: "Corretiva → regras → Indexus (Recommended)"; o fix da fila entra logo após a corretiva.

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

For this Spec the answer grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, a release after the Spec, with up to five corrections. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: fix` with a `patch` Release Ceiling. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

## How the governed set was measured

The author measured the governed set on 2026-11-18 on the branch `docs/plan-0408`, based on `1a8e08b4`, with `speccheck.GovernedPath`.

- The measurement used a `go test -overlay` probe in `internal/speccheck` whose test file lived outside the repository; it wrote nothing to the tree.
- It covered 19 paths: every source, test, guide, skill, record and document file the TechSpec names or expects the Tasks to change, ADR-0435, this Spec's PRD, and three controls.
- The 3 paths above are the governed ones this Spec's design changes: the canonical `deliver` reference of the Roundfix skill and the two `SKILL.md` files whose version fields rise.
- These governed paths are not granted, and no Task may change them: the three controls, the `Makefile`, `.github/workflows/ci.yml` and `go.mod`.

Spec 0407 delivers first and may raise the same Roundfix skill version fields. If, on the tree the Run starts from, a sanctioned command below rewrites a governed path outside `paths`, the operator amends this record before the Run, outside the change it authorizes.

## What is not governed

- `internal/delivery/engine.go`, `internal/delivery/park_class.go`, `internal/delivery/base_update_test.go`, `internal/delivery/park_class_test.go` and `internal/cli/deliver_base_update_test.go`.
- The mirror `skills/roundfix/references/deliver.md`, and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/deliver.md`, `internal/docscontract/command_documentation_test.go`, `CONTEXT.md`, `docs/references/behavior-surfaces.json` and ADR-0435.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the skill mirror, the version record and the Behavior Surface Record. This declaration records the regeneration that follows the approved edits and adds no source path.

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
  - any Vendored Skill under `.agents/skills`, or any owned skill other than `roundfix`;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot;
  - `docs/references/coverage-record.json`;
  - any Archive Record or any ADR other than ADR-0435;
  - the `### QA settlement` section of any skill.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads another repository, except a temporary one the test itself created;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
