---
status: approved
granted: 2026-10-06
action: make a normal archive drop a Spec's raw QA evidence and keep its QA Reports with an Evidence Manifest, let the Delivery Queue accept that cut as an exact move, add a request-only cut of one archived Spec to the Archive Command, stop the coverage record and the repository-copy helper from depending on archived evidence, and describe the cut in the archive-spec, qa-gate and Roundfix skills and in the Spec workflow Baseline guides
consuming: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - docs/references/coverage-record.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/spec/archive.go
  - internal/spec/coverage_test.go
  - skills/archive-spec/SKILL.md
  - skills/baseline_skill_contract_test.go
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0238

On 2026-10-06 the maintainer answered through AskUserQuestion. The question
was "Uma Spec de retenção no repositório (o archive descarta a evidência
bruta de QA e guarda o relatório)?", and the answer was "Sim, como Spec". The
option chosen said that the Spec fixes the three tests that read archived
evidence, records an ADR allowing the cut, applies to new Specs and
optionally to already-archived ones, and ships in its own release. This is
the explicit approval of a removal that ADR-0215 requires. In the same
session the maintainer answered "Não agora" to rewriting the Secondbrain's
Git history, which stays out of scope, and the broad mirror exclusion was
done separately (`.secondbrain-export`, Pull Request 419).

On 2026-09-30 the maintainer had granted the skills ("considere autorizado a
ajustar todas as skills se necessário") and the Baseline source and guides
("Autorizar os dois").

The governed set was measured with `GovernedPath` on the authoring branch at
`08660fcc`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. The derived Baseline files
were measured by applying task_04's clause edits and three skill edits in a
disposable clone and running `make skills-sync`, the skill version record,
`make baseline-digests` twice and the Managed Refresh twice there; the
second runs changed nothing.

## Why each governed path is unavoidable

- `docs/references/coverage-record.json` and `internal/spec/coverage_test.go`
  — the record lists three Go packages inside archived evidence, so removing
  that evidence fails `make verify`; the collector stops counting packages
  under `docs/` and the record is re-recorded.
- `skills/baseline_skill_contract_test.go` — it holds `copyTrackedRepository`,
  which fails `make verify-docs` on a tracked file missing from the working
  tree.
- `internal/spec/archive.go` — the archive itself calls the cut and accepts
  an evidence link that now reaches the manifest.
- `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md` and the three `SKILL.md`
  mirrors — they describe what archive keeps and what the gate must write in
  its report; each owned skill's version is raised with its content.
- `internal/baseline/assets/modules/spec-workflow.json` — two clauses say
  that `qa/` evidence stays under the Spec folder and that archived Specs
  stay byte-identical; each gains one sentence about the cut.
- `docs/agents/docs-layout.md`, `docs/agents/spec-routing.md`,
  `docs/agents/setup-context.json`, the two formatter goldens and
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json` — the
  deterministic fallout of that module edit, written by `make
  baseline-digests` and the Managed Refresh, never by hand.

## What is not governed

The new and changed Go sources and tests under `internal/spec`,
`internal/cli` and `internal/speccheck` other than those listed above,
`internal/cli/deliver_workflow.go`, `internal/cli/archive.go`, the new
`skills/copy_tracked_repository_test.go` and
`internal/baseline/evidence_cut_clause_test.go`, the Baseline test data under
`internal/baseline/testdata/`, the skill reference mirror
`skills/roundfix/references/archive.md`,
`skills/testdata/owned-skill-versions.json`, the archive user guide and
ADR-0243 are ordinary.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions,
the derived Baseline artifacts and this repository's Setup Manifest. This
declaration records the regeneration that follows the approved edits and
adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
records the raised versions, and
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
renders this repository's guides and Setup Manifest; a second refresh reports
no file change. `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record`
re-records the coverage record. No digest pin, golden, generated guide or
record is hand-edited.

## Limits

- No file under `docs/history/` is deleted, moved or rewritten by this Spec;
  the cut of the existing History Root is a later change.
- No rewrite of Git history, here or in the Secondbrain.
- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to the `### QA settlement` section of any skill, to QA
  eligibility, to the review scope or to `CONTEXT.md` and `CHANGELOG.md`.
- No test, Verification command or QA row reaches the network or reads or
  writes the real `~/.roundfix`.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
