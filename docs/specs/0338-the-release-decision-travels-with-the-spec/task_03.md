---
task: task_03
spec: 0338-the-release-decision-travels-with-the-spec
status: completed
type: infra
complexity: medium
---

# Task 03: The Baseline's release and authorization clauses name the Release Ceiling

## Overview

The Baseline's release clause says a generic release request authorizes only
a conclusive patch plan, and its authorization clause lists the record's keys
without a ceiling, so an adopter's agent cannot learn that a recorded Release
Ceiling answers the plan. Both clauses gain one sentence or phrase; the
modules' next version is recorded, the derived Baseline files are
regenerated, and this repository's managed guides are refreshed. This answers
the adopted Backlog Entry of 2026-11-01, "A Spec's authorization carries the
release decision"
([references](references/2026-11-01-a-spec-authorization-carries-the-release-decision.md)),
whose Expected section asks the Baseline wording to change together.

## Requirements

1. MUST append to the `guidance` of `clause.core.plan-the-release-first`
   (`internal/baseline/assets/modules/core.json`) the sentence of
   `_techspec.md` → API Contract 11, keeping every existing sentence byte for
   byte, and raise the `version` of `rule.core.git-delivery` by one over the
   value the tree holds.
2. MUST insert into the `guidance` of
   `clause.spec.project-constraints-02-tooling-authorization`
   (`internal/baseline/assets/modules/spec-workflow.json`), right after the
   operations list, the words API Contract 11 gives, keeping the sentence
   `internal/baseline/authorization_clauses_test.go` pins, and raise the
   `version` of `rule.spec.project-constraints` by one. Both clauses keep
   their `id`, `enforcement` and `replaces`; nothing else in either module is
   edited by hand.
3. MUST add `TestTheReleaseClauseLetsACeilingAnswerTheApproval` to
   `internal/baseline/release_clause_test.go`, calling `t.Parallel()` first,
   that reads the embedded catalog as the existing test does and asserts the
   new core sentence in the core module and in the formatter golden
   `docs/agents/agent-instructions.md`, and the words
   "(`release_ceiling`: `patch`, `minor` or `major`)" in the spec-workflow
   module and in the formatter golden `docs/agents/spec-routing.md`. The
   existing test keeps its bytes.
4. MUST then run, in this order, and hand-edit none of their outputs:
   - `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
     which writes the next version of both modules and records them in
     `internal/baseline/module-versions.json`;
   - `make baseline-digests`, which rewrites the derived files this Task
     declares (the two Standard TypeScript formatter goldens, the Standard
     TypeScript profile, the three catalog files and the four plan
     characterization goldens);
   - `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
     the Managed Refresh of this repository, which rewrites
     `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md` and
     `docs/agents/setup-context.json`, then the same command again, which
     must report no file change. Run both with `HOME` set to a temporary
     directory that is removed afterwards, so the refresh reads and writes no
     real Roundfix Home.
   A file the commands change that this Task does not declare is a
   stop-and-report condition.
5. MUST NOT change the Source Baseline corpus, manifest, accounting or index,
   `internal/baseline/assets/retention/**`, any other module, any skill, any
   user guide page, `.roundfixrc.yml`, the Makefile, the CI workflows or any
   Go production file.

## Subtasks

- [ ] Append the core sentence and insert the spec-workflow words; raise both rule versions.
- [ ] Add the clause test.
- [ ] Record the module versions, run `make baseline-digests` and the Managed Refresh twice.

## Acceptance Criteria

- [ ] The embedded clauses, the Standard TypeScript goldens and this repository's `agent-instructions.md` and `spec-routing.md` carry the Release Ceiling wording, and the existing release and authorization sentences keep their bytes.
- [ ] A second Managed Refresh changes nothing, and the Setup Manifest test passes.
- [ ] The Baseline package's catalog, characterization, clause and module-version tests pass.

## Context

- instruction: `docs/adr/0349-a-spec-release-ceiling-answers-the-release-plan-approval.md`
- instruction: `internal/baseline/module_versions_test.go`
- instruction: `internal/baseline/catalog_test.go`
- instruction: `internal/baseline/plan_characterization_test.go`
- instruction: `internal/baseline/clause_characterization_test.go`
- instruction: `internal/baseline/authorization_clauses_test.go`
- instruction: `internal/docscontract/publicdocs_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/release_clause_test.go`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`

## Consumers

- The two clause texts and the catalog digests — consumers: `./internal/baseline`, `./internal/docscontract`

## Verification

- `for file in docs/agents/agent-instructions.md internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'answers a minor or major plan at or below that ceiling, never a breaking or manually classified one' || { printf '%s lacks the ceiling sentence\n' "$file" >&2; exit 1; }; done; for file in docs/agents/spec-routing.md internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'with the version-only paths that release may change' || { printf '%s lacks the ceiling keys\n' "$file" >&2; exit 1; }; done; go test -timeout 20m -count=1 -run '^(TestTheReleaseClauseLetsACeilingAnswerTheApproval|TestTheReleaseClauseNamesTheSkillsAndGuidesCheck|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestCatalogDiagnosticCharacterization|TestBaselinePlanCharacterization|TestFormatterComposition|TestBaselineCompatibilityCorpus|TestNoTwoBaselineClausesShareText|TestBaselineClauseForceIsCharacterized)$' ./internal/baseline && go test -timeout 20m -count=1 ./internal/baseline` — expected: exit 0. Before this Task the guides and goldens lack the ceiling wording, so the first check fails; afterwards the four rendered copies carry it and the clause, module record, catalog, characterization and whole Baseline suites pass.
- `out="$(go test -timeout 20m -count=1 -v -run '^TestTheReleaseClauseLetsACeilingAnswerTheApproval$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheReleaseClauseLetsACeilingAnswerTheApproval ' || { echo 'missing pass: TestTheReleaseClauseLetsACeilingAnswerTheApproval' >&2; exit 1; }; go test -timeout 20m -count=1 -tags docscontract -run '^(TestTheSetupManifestIsCurrent|TestReleasePlanDocumentation.*)$' ./internal/docscontract` — expected: exit 0. Before this Task the clause test does not exist, so the pass line is missing; afterwards it passes and this repository's Setup Manifest matches the regenerated guides.

## References

- `_prd.md` → Sources 2; Goals; Core Feature 5; Success Metric 7
- `_techspec.md` → Current behavior 7; API Contract 11; Testing Approach; Build Order 3
- ADR-0349

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, tech-writer, setup-context-driven, roundfix, evidence-gate.

The two Baseline clauses now carry API Contract 11's Release Ceiling wording.
The core guidance appends the full sentence; the authorization guidance inserts
both keys immediately after the operations list. The owning rule versions rose
once (git-delivery 8 → 9; project-constraints 7 → 8). The sanctioned recorder
raised the module versions (core 25 → 26; spec-workflow 22 → 23) and appended
their digests. No generated output was edited by hand.

Acceptance evidence:

- Clause wording and preservation: added
  `TestTheReleaseClauseLetsACeilingAnswerTheApproval`, with `t.Parallel()` as
  its first statement and named core/spec-workflow subtests reading the
  embedded catalog. It asserts the full new sentence or phrase in the exact
  clause and its corresponding formatter golden. Before the module edits,
  `go test ./internal/baseline -run '^TestTheReleaseClauseLetsACeilingAnswerTheApproval/(core|spec-workflow)$' -count=1`
  exited 1 with all four expected missing-wording assertions. After
  regeneration, the focused clause/catalog command below exited 0.
  A structural comparison against HEAD confirmed that only each module's
  guidance, owning rule version and recorded module version changed; removing
  the inserted words recovers the previous authorization guidance exactly,
  and the core guidance preserves its previous bytes as a prefix. Clause
  identities, enforcement and replaces are unchanged. A byte-prefix comparison
  confirmed the existing release test is untouched. The generated local guides
  were inspected and carry the same wording.
- Refresh idempotence and Setup Manifest: ran
  `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
  twice, after the generators, with HOME set to a temporary directory. The
  existing module cache was reused with network access disabled. Both runs
  exited 0. The first reported three file changes, exclusively the two
  declared guides and Setup Manifest. The second reported `File changes: 0`,
  `Idempotence: verified`, and approved postimages verified. SHA-256 snapshots
  of every declared Task path were identical across the second refresh. The
  temporary HOME was removed. The focused Setup Manifest test below passed.
- Catalog, characterization, clause and module records: ran the required
  sanctioned regeneration in order:
  1. `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
     — exit 0; only the two modules and their version record changed.
  2. `make baseline-digests` — exit 0, `ok: true`; regenerated exactly the
     ten declared derived files: the two formatter goldens, profile, three
     catalog files and four plan characterization goldens. Its regeneration
     and strict catalog validation steps all passed.
  3. The two Managed Refresh invocations described above — exit 0 each.
  The full Baseline suite and authored Verification remain for the Daemon.

Additional focused checks:

- `go test -count=1 ./internal/baseline -run '^(TestTheReleaseClause.*|TestAuthorizationClauses.*|TestEmbeddedCatalog|TestCatalogDigest)$'`
  — exit 0, including the new ceiling test and existing release/authorization
  sentence checks.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$' -count=1`
  — exit 0.
- `gofmt -l internal/baseline/release_clause_test.go` — exit 0, no output.
- Changed-path postflight against this Task's exact interface list — all
  18 changed paths declared, including this Task file; no other Task or
  Task Graph mutation.

The Managed Refresh reported the existing nested-carrier warnings for the
formatter fixture and Source Baseline corpus AGENTS.md; both carriers remained
unchanged. No follow-up implementation or blocker was identified. Task status
and the Subtasks/Acceptance Criteria checkboxes remain Daemon-owned. No declared
Verification command, repository Verification, full suite, commit, push or pull
request was performed in this turn.
