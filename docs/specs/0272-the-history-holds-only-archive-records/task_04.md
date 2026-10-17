---
task: task_04
spec: 0272-the-history-holds-only-archive-records
status: pending
type: docs
complexity: high
---

# Task 04: The Baseline states retirement once and this repository's history leaves

## Overview

The `context-workflow` module still tells an Agent to move retired Findings
and Backlog Entries into history families as reduced entries, keeps retired
ADRs whole, calls handoffs the capture door, and states the retirement rule
in several clauses; the rendered docs-layout region is 14,494 bytes
(`_techspec.md` → Current behavior). This Task rewrites the module to one
retirement clause, adds the conditional lesson capture to the `secondbrain`
module, regenerates the derived files, and runs this repository's Managed
Refresh, which turns task_03's check on and, in the same plan, removes the 98
reduced Findings, the 32 reduced Backlog Entries and the ten `superseded`
records and gives every citation and broken record link its Git Locator. It
is verifiable on its own through the clause tests, the tree and a second
refresh that reports the repository current.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-16,
   [the Baseline keeps retired ADRs whole and names two capture doors](references/2026-10-16-the-baseline-keeps-retired-adrs-and-names-two-capture-doors.md).
2. MUST implement API Contract 15 in
   `internal/baseline/assets/modules/context-workflow.json`, keeping each
   clause's identity and force: `clause.context.docs-one-job-per-directory`
   states the directory jobs, with the history root for Archive Records,
   `docs/_inbox/` only for a repository without a Secondbrain, handoffs for
   session continuity and no "capture door", and keeps the authorization
   record locations and the Spec-owned adopted reference rule;
   `clause.context.findings-09-archive` opens with
   `spec.RetirementRuleSentinel`, states that a Finding, Rollup or Backlog
   Entry with a terminal status (`deferred` included), a handoff on the user's
   explicit confirmation, a finished Review Artifact, a rejected or deprecated
   ADR and a superseded ADR once an accepted successor names it with its date
   are deleted in the change that records the disposition, that Git keeps the
   text, that a Spec absorbing a source adopts it and any other retirement
   names its reason in the commit, that a citation of a deleted file becomes
   the Git Locator `roundfix spec check` prints, and that the runtime's
   `roundfix history sanitize` removes what an earlier layout left; it lists
   `clause.context.findings-10-live-work-health` and
   `clause.context.findings-11-rollup-closure` in `replaces`, both mandatory
   like it, and those two clauses leave the module. `findings-05-done`,
   `findings-08-rollup` and the last paragraph of
   `backlog-01-operational-contract` stop naming `docs/history/` families and
   point at the retirement clause. A preference stays an obligation about its
   exception (ADR-0202), and no clause gains force for what the product does
   not enforce (ADR-0186).
3. MUST extend `clause.secondbrain.capture-trigger` in
   `internal/baseline/assets/modules/secondbrain.json`: the change that closes
   a Finding states whether it teaches a lesson its ADR, Archive Record or
   reference does not already state, and only then captures it as one pending
   Inbox Entry with its source pinned as `<repository>@<commit>:<path>`.
4. MUST record the next module versions with
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   then run `make baseline-digests` twice (the second rewrites nothing), and
   change only the derived files this Task declares; a derived file outside
   them stops the Task.
5. MUST, before the refresh, move the routing of the promoted Finding
   `docs/references/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md`
   from the 0129 record to the successor its record names in `superseded_by`,
   in a dated addendum that names the prior owner by its Git Locator, and
   replace its broken `../_prd.md` link with the Git Locator the check prints;
   and fix the broken link in
   `docs/references/jev-router-measurement-2026-10-04.md` the same way. No
   other file is edited by hand.
6. MUST run the Managed Refresh itself, as this Task's own implementation
   work:
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`,
   which refreshes the managed guides and the Setup Manifest and applies the
   history section; then run it again without `--yes` and require exit 0 with
   no pending history. The run with `--yes` is not a Verification command and
   the Daemon never runs it: the Verification below only runs the refresh
   without `--yes`, to prove the tree current after this Task's own refresh.
   A diff that still holds `docs/history/findings` or `docs/history/backlog`
   fails settlement on its `deletes:` entries, as the Run of 2026-10-17 06:28
   UTC did when its Agent left the refresh to the Daemon. After it, the History Root holds only Archive Records
   (243 files), none of the ten `superseded` records 0120 to 0129 remains,
   the five user guide and reference links to reduced Findings and the 40
   broken links in Archive Records carry Git Locators, and
   `roundfix spec check` reports no `SC-HISTORY-RETAINED` finding. The ten
   records are not listed under `deletes:` because a Spec file naming them
   would cite them and keep them; the Daemon records them, the 130 reduced
   entries and the rewritten records under `## Recorded paths`.
7. MUST record the Coverage Review of the five surfaces the PRD's `## Skills`
   section lists as unchanged (`command reconcile`, `command stop` and the
   guides `docs/user-guide/commands/detached-runs.md`, `reconcile.md` and
   `stop.md`) as a dated `review` entry in `docs/references/skill-coverage.json`
   naming this Task and the rewritten link, and re-record the Behavior Surface
   Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
8. MUST describe deletion in `docs/user-guide/context-driven-development.md`,
   replacing the paragraph that moves terminal Findings and Backlog Entries
   into history families, with the phrase `Git Locator` and without
   `Write each as a reduced entry`.
9. MUST change these existing tests to the new contract and name each change
   in the Result: in `internal/baseline/reduced_retirement_clauses_test.go`,
   `TestTheRetirementClausesCarryTheirText`,
   `TestTheRetirementClausesRenderInTheDocsLayoutGuide`,
   `TestAFindingInTheGuideReducedFormNeedsNoSanitize` and
   `TestAnAdopterRetainsTheRetirementClauses`;
   `TestHistorySanitizeClauseIsAppended`,
   `TestArchiveRecordClausesAreAppended`,
   `TestLifecycleClausesCarryTheScopedWording` and
   `TestBaselineClauseForceIsCharacterized`; and
   `TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily`, whose
   Findings family is no longer present in this repository.
10. MUST add, in `internal/baseline/retirement_rule_clauses_test.go`, each
   calling `t.Parallel()` first: `TestTheRetirementRuleIsStatedOnce` (no other
   clause of the module names a retired history family, a reduced entry,
   "stay whole" or a handoff capture door), `TestTheDocsLayoutRegionStaysLean`
   (the formatter golden's `guide.docs-layout` region is at most 11,000
   bytes), `TestTheRetirementClauseCarriesTheSentinel` (it reads
   `spec.RetirementRuleSentinel`), `TestTheRetirementClauseReplacesTheClausesItAbsorbs`
   (the Standard TypeScript Source Baseline records both removed clauses as
   replaced) and `TestTheCaptureTriggerAsksForAConditionalLesson`.
11. MUST, when this Task removes or renames a test, regenerate the test
    inventory with `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
    after its last test change, so `TestCoverageEquivalence` reports no
    regression. The sanctioned regeneration in `_authorization.md` already
    declares that command, so this Task MUST NOT edit
    `internal/spec/coverage_test.go`.
12. MUST run both Managed Refresh runs on the tree as it is, with the tracked
    `CLAUDE.md` link left a link: Requirement 14 makes the citation scan skip
    it, so no file is materialized, replaced or restored around a refresh.
13. MUST change the two Repository Contract Tests that read the removed
    families, and name each change in the Result:
    `TestDocsLayoutGuideNamesEveryHistoryFamily` in
    `internal/docscontract/publicdocs_test.go` becomes
    `TestDocsLayoutGuideNamesOnlyTheArchiveRecordFamily`, which requires the
    guide to name `spec.ArchiveDir(spec.ArchiveKindSpec)` and to name no other
    kind's `spec.ArchiveDir`; and `copyTrackedRepository` in
    `skills/baseline_skill_contract_test.go` skips a tracked path the working
    tree no longer holds (`os.ErrNotExist` from `os.Lstat`), because the
    change removes it, so `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`
    copies the tree the commit will hold. Every other Repository Contract
    Test passes unchanged; `make verify-contracts` measured on 2026-10-17
    against this Task's first delivery found only these two.
14. MUST, before the first Managed Refresh, make `HistoryCitations` in
    `internal/spec/history_entries.go` skip a tracked Markdown path whose own
    entry is a symbolic link: it is not read, not reported and never
    rewritten, and its target is scanned only under the target's own tracked
    path. The citation scan of the `superseded` unit (API Contract 8) skips
    such a link the same way, so a link to a file outside the repository is
    never followed. A regular file under a symbolically linked directory is
    still refused, the apply path keeps refusing a link it is asked to
    rewrite, and every other refusal and the plan output keep their text.
    Today `roundfix baseline update` on this tree prints `History result:
    citing file CLAUDE.md is not a regular repository file`, and adopters keep
    the same `CLAUDE.md` link.
15. MUST add, each calling `t.Parallel()` first:
    `TestHistoryCitationsSkipATrackedSymbolicLink` in
    `internal/spec/history_symlink_test.go` (a fixture whose tracked
    `AGENTS.md` cites a removed path and whose tracked `CLAUDE.md` links to
    it: one citation, in `AGENTS.md`; after apply `CLAUDE.md` is still a link
    to `AGENTS.md` and `AGENTS.md` carries the Git Locator), and
    `TestHistorySanitizeAppliesBesideATrackedSymbolicLink` in
    `internal/cli/history_symlink_test.go` (`roundfix history sanitize
    --apply` on the same kind of fixture exits 0 and leaves the link a link).

## Subtasks

- [ ] Make the history citation scan skip a tracked symbolic link, with its two tests.
- [ ] Rewrite the `context-workflow` clauses and extend the `secondbrain` capture trigger.
- [ ] Record the module versions and regenerate the derived files.
- [ ] Move the 0129 routing and fix the two broken reference links by hand.
- [ ] Run the Managed Refresh with its history section, then again to prove it current.
- [ ] Update the context-driven development guide and the changed tests, and add the clause tests.

## Acceptance Criteria

- [ ] The rendered docs-layout region is at most 11,000 bytes and states retirement in one clause.
- [ ] The History Root holds only Archive Records and none of the ten superseded records.
- [ ] A second Managed Refresh exits 0 with no pending history, and the check reports no retained history.
- [ ] Every derived file changed is one this Task declares.
- [ ] The tracked `CLAUDE.md` link stays a link and is never read through, reported or rewritten.

## Consumers

- History Citation planning for both history commands — consumers: `./internal/spec`, `./internal/cli`
- Baseline modules, the rendered guides and the Setup Manifest — consumers: `./internal/baseline`
- This repository's History Root and its citing files — consumers: `./internal/spec`

## Context

- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `docs/adr/0202-a-stack-rule-carries-force-and-a-preference-names-what-licenses-the-exception.md`
- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- instruction: `docs/agents/specific-repository.md`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/modules/secondbrain.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/secondbrain.md`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/reduced_retirement_clauses_test.go`
- interface: `internal/baseline/history_sanitize_clause_test.go`
- interface: `internal/baseline/archive_record_clause_test.go`
- interface: `internal/baseline/adopter_neutral_clauses_test.go`
- interface: `internal/baseline/clause_characterization_test.go`
- creates: `internal/baseline/retirement_rule_clauses_test.go`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md`
- interface: `docs/references/jev-router-measurement-2026-10-04.md`
- interface: `docs/references/model-selection.md`
- interface: `docs/user-guide/commands/detached-runs.md`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `docs/user-guide/commands/stop.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- deletes: `docs/history/findings`
- deletes: `docs/history/backlog`
- interface: `docs/references/coverage-record.json`
- interface: `internal/spec/history_entries.go`
- interface: `internal/spec/history_entries_test.go`
- creates: `internal/spec/history_symlink_test.go`
- creates: `internal/cli/history_symlink_test.go`
- interface: `internal/cli/history_test.go`
- interface: `internal/docscontract/publicdocs_test.go`
- interface: `skills/baseline_skill_contract_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestHistoryCitationsSkipATrackedSymbolicLink$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestHistoryCitationsSkipATrackedSymbolicLink ' || { printf 'missing passing test TestHistoryCitationsSkipATrackedSymbolicLink\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^TestHistorySanitizeAppliesBesideATrackedSymbolicLink$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestHistorySanitizeAppliesBesideATrackedSymbolicLink ' || { printf 'missing passing test TestHistorySanitizeAppliesBesideATrackedSymbolicLink\n' >&2; exit 1; }; go test -count=1 -run 'History' ./internal/spec ./internal/cli || exit 1` — expected: exit 0. Before this Task the two symbolic link tests do not exist, so the PASS checks fail.
- `out="$(go test -count=1 -v -run '^(TestTheRetirementRuleIsStatedOnce|TestTheDocsLayoutRegionStaysLean|TestTheRetirementClauseCarriesTheSentinel|TestTheRetirementClauseReplacesTheClausesItAbsorbs|TestTheCaptureTriggerAsksForAConditionalLesson)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheRetirementRuleIsStatedOnce TestTheDocsLayoutRegionStaysLean TestTheRetirementClauseCarriesTheSentinel TestTheRetirementClauseReplacesTheClausesItAbsorbs TestTheCaptureTriggerAsksForAConditionalLesson; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Clause|Retire|Catalog|Formatter|PlanCharacterization|Compatibility|ModuleVersion|SourceBaseline' ./internal/baseline || exit 1; go test -count=1 -run '^TestArchiveLayoutCharacterization' ./internal/spec || exit 1` — expected: exit 0. Before this Task the five clause tests do not exist, so the PASS check fails.
- `matches="$(find docs/history -type f ! -path 'docs/history/specs/*' -print)" || exit 1; test -z "$matches" || { printf 'retained history:\n%s\n' "$matches" >&2; exit 1; }; superseded="$(find docs/history/specs -name '012[0-9]-*.md' -print)" || exit 1; test -z "$superseded" || { printf 'superseded records remain:\n%s\n' "$superseded" >&2; exit 1; }; if tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'Write each as a reduced entry'; then printf 'guide still writes reduced entries\n' >&2; exit 1; fi; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'Git Locator' || { printf 'guide does not name the Git Locator\n' >&2; exit 1; }; dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$dir/roundfix" ./cmd/roundfix || exit 1; HOME="$dir" "$dir/roundfix" baseline update --repo . --no-skills --format text || { printf 'the Managed Refresh is not current\n' >&2; exit 1; }; HOME="$dir" "$dir/roundfix" spec check 0272-the-history-holds-only-archive-records || exit 1; go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestDocsLayoutGuideNamesOnlyTheArchiveRecordFamily)$' ./internal/docscontract || exit 1; go test -count=1 -tags repocontract -run '^TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical$' ./skills || exit 1` — expected: exit 0. Before this Task the History Root holds reduced entries and the refresh exits 3, so the command fails.

## References

- `_prd.md` → Goals; User Stories 1, 4 and 5; Core Feature 1; Core Feature 5; Core Feature 6; Core Feature 7; Success Metric 1; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contract 8; API Contract 10; API Contract 15; Invariant 6; Coverage Map deviations; Build Order 4
- ADR-0278; ADR-0186; ADR-0202; ADR-0257; ADR-0095
