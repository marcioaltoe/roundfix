---
task: task_04
spec: 0272-the-history-holds-only-archive-records
status: completed
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

## Result

Implementation evidence for the daemon handoff:

- `HistoryCitations` now skips a tracked Markdown path whose own entry is a symbolic link, while files reached through a symbolic-linked directory remain refused. The `superseded` citation scan applies the same rule and ignores tracked paths already removed by an earlier refresh. `AGENTS.md` remains the rewrite surface in the fixture; `CLAUDE.md` remains a link to it and is never read, reported or rewritten.
- Added `TestHistoryCitationsSkipATrackedSymbolicLink` and `TestHistorySanitizeAppliesBesideATrackedSymbolicLink`; both focused tests passed. The five new clause tests also passed: `TestTheRetirementRuleIsStatedOnce`, `TestTheDocsLayoutRegionStaysLean`, `TestTheRetirementClauseCarriesTheSentinel`, `TestTheRetirementClauseReplacesTheClausesItAbsorbs` and `TestTheCaptureTriggerAsksForAConditionalLesson`.
- Rewrote the Baseline retirement guidance into one mandatory clause, recorded both absorbed clause IDs in `replaces`, extended conditional Secondbrain lesson capture, regenerated the managed guides and Setup Manifest, and updated the guide, catalog tests, force characterization, archive layout characterization, and adopter wording.
- Changed the named existing tests: `TestTheRetirementClausesCarryTheirText`, `TestTheRetirementClausesRenderInTheDocsLayoutGuide`, `TestAFindingInTheGuideReducedFormNeedsNoSanitize`, `TestAnAdopterRetainsTheRetirementClauses`, `TestHistorySanitizeClauseIsAppended`, `TestArchiveRecordClausesAreAppended`, `TestLifecycleClausesCarryTheScopedWording`, `TestBaselineClauseForceIsCharacterized`, `TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily`, `TestDocsLayoutGuideNamesOnlyTheArchiveRecordFamily`, and `copyTrackedRepository` now account for deletion and the Archive Record only history root.
- Recorded module versions; `make baseline-digests` regenerated the declared derived files and the second run reported `changed=false`. Re-recorded skill coverage and the Behavior Surface Record for the five unchanged surfaces. The coverage inventory command wrote its intended record; its package guard reported the intentional record mutation, and the read-only `TestCoverageEquivalence` rerun passed.
- Moved the promoted Finding's routing to `0159-archive-override-and-authoring-rules`, retained the former `0129` owner as a Git Locator addendum, fixed both broken reference links, and ran the Managed Refresh with history application. The final refresh removed 139 retired files plus the remaining `0129` record; the following refresh exited 0 with `History: current` and no refused units. The tree has zero non-Spec history files, zero `0120`–`0129` records, and a 9,665-byte rendered `guide.docs-layout` region.
- Focused checks passed for the changed Spec, CLI, Baseline and coverage paths, including `TestCoverageEquivalence`, `TestPlanSupersededRecordsSkipsCitedAndGrantedRecords`, `TestTheRetirementClausesCarryTheirText`, `TestHistorySanitizeWithNothingPending`, and compile checks for all three changed Go packages.
- Verification feedback identified two stale contract assertions: `TestLifecycleClausesCarryTheScopedWording` still expected the removed terminal-state enumeration, and `TestHistorySanitizeClauseIsAppended` still searched the removed `docs-one-job-per-directory` clause. They now assert the deletion wording and the consolidated `findings-09-archive` clause. Focused reruns of both tests and `TestTheRetirementClauseCarriesTheSentinel` passed.

The Daemon owns the Task status and must run the declared Verification, including the final `roundfix spec check` and contract checks. This handoff does not claim a terminal Task verdict.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/history/backlog/2026-08-06-atomic-inbox-capture-helper.md`
- `docs/history/backlog/2026-08-08-a-failed-proof-appends-a-cleanup-error-the-maintainer-cannot-act-on.md`
- `docs/history/backlog/2026-08-08-a-runtime-that-advertises-a-catalog-is-not-unusable.md`
- `docs/history/backlog/2026-08-08-record-usage-per-agent-session.md`
- `docs/history/backlog/2026-08-09-a-verification-command-passes-only-by-exiting-zero.md`
- `docs/history/backlog/2026-08-09-the-canonical-method-lives-in-rendered-guides-instead-of-skills.md`
- `docs/history/backlog/2026-08-09-the-qa-gate-commits-scratch-repositories-and-binaries.md`
- `docs/history/backlog/2026-08-09-the-roundfix-skill-carries-nine-commands-in-one-file.md`
- `docs/history/backlog/2026-08-10-an-excluded-artifact-outlives-the-decision-that-excluded-it.md`
- `docs/history/backlog/2026-08-10-one-reader-in-cli-still-couples-verify-to-the-docs-tree.md`
- `docs/history/backlog/2026-08-10-the-gate-accepts-a-manifest-that-names-a-catalog-that-is-gone.md`
- `docs/history/backlog/2026-08-10-the-preflight-reads-the-catalog-fixtures-as-repository-carriers.md`
- `docs/history/backlog/2026-08-12-a-hook-failure-kills-a-run-that-already-verified-its-work.md`
- `docs/history/backlog/2026-08-12-a-redirected-verification-hands-the-agent-an-empty-diagnostic.md`
- `docs/history/backlog/2026-08-12-archive-refuses-a-graph-that-declined-the-qa-gate.md`
- `docs/history/backlog/2026-08-12-release-plan-requires-a-v-prefixed-tag.md`
- `docs/history/backlog/2026-08-12-speak-acp-directly-instead-of-through-acpx.md`
- `docs/history/backlog/2026-08-12-the-secondbrain-guide-permits-capture-without-triggering-it.md`
- `docs/history/backlog/2026-08-13-a-refused-gate-writes-a-report-its-own-contract-rejects.md`
- `docs/history/backlog/2026-08-13-the-changed-path-audit-does-not-know-sanctioned-fallout.md`
- `docs/history/backlog/2026-08-14-a-spec-archives-before-continuous-integration-answers.md`
- `docs/history/backlog/2026-08-14-a-spec-that-coins-a-term-cannot-pass-its-own-gate.md`
- `docs/history/backlog/2026-08-14-staging-a-task-commit-fails-when-the-task-deleted-a-file.md`
- `docs/history/backlog/2026-08-14-the-archive-destination-still-depends-on-which-spec-root-is-configured.md`
- `docs/history/backlog/2026-08-14-the-blocked-row-parser-requires-a-literal-the-message-never-names.md`
- `docs/history/backlog/2026-08-14-the-prd-template-and-the-checker-disagree-on-the-tooling-row.md`
- `docs/history/backlog/2026-08-14-the-underscored-orphan-review-root-is-never-migrated.md`
- `docs/history/backlog/2026-08-14-the-vacuous-verification-event-lists-the-commands-that-passed.md`
- `docs/history/backlog/2026-09-25-jev-gateway-baseline-before-routing.md`
- `docs/history/backlog/2026-10-03-a-grok-fallback-for-docs-and-chore.md`
- `docs/history/backlog/2026-10-05-the-judge-log-under-counts-router-spend.md`
- `docs/history/backlog/2026-10-12-the-worktree-admin-lock-test-is-timing-dependent.md`
- `docs/history/findings/2026-07-16-vortex-pr87-detached-watch-notification.md`
- `docs/history/findings/2026-07-17-global-run-storage-sanitation-and-compaction.md`
- `docs/history/findings/2026-07-23-setup-context-driven-adoption-process-improvements.md`
- `docs/history/findings/2026-07-24-greenfield-agent-guidance-acceptance-target.md`
- `docs/history/findings/2026-07-25-npm-trusted-publishing-and-release-preflight.md`
- `docs/history/findings/2026-07-25-spec-owned-reference-lifecycle.md`
- `docs/history/findings/2026-07-26-baseline-profile-refresh-retention-gap.md`
- `docs/history/findings/2026-07-26-claude-adapter-configoptions-migration.md`
- `docs/history/findings/2026-07-26-vortex-baseline-capability-remediation.md`
- `docs/history/findings/2026-07-27-bare-go-build-writes-an-untracked-root-binary.md`
- `docs/history/findings/2026-07-27-claude-adapter-standardization.md`
- `docs/history/findings/2026-07-27-derived-skill-digest-pins-have-no-regeneration-path.md`
- `docs/history/findings/2026-07-27-owner-identity-forks-ps-and-fails-closed-under-load.md`
- `docs/history/findings/2026-07-27-sandboxed-agents-cannot-reach-the-default-go-cache.md`
- `docs/history/findings/2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs.md`
- `docs/history/findings/2026-07-28-profiles-configure-replaces-the-whole-profiles-map.md`
- `docs/history/findings/2026-07-28-qa-gate-cannot-reach-pull-request-journeys.md`
- `docs/history/findings/2026-07-28-same-day-qa-reruns-are-ignored-by-the-verdict-selector.md`
- `docs/history/findings/2026-07-28-tooling-tasks-need-a-green-repo-and-an-undocumented-commit-choreography.md`
- `docs/history/findings/2026-07-29-doctor-requires-roundfix-own-development-skills.md`
- `docs/history/findings/2026-07-29-qa-cycle-cost-is-cold-environments-and-agent-turns.md`
- `docs/history/findings/2026-07-29-qa-cycle-latency-and-detector-placement.md`
- `docs/history/findings/2026-07-29-qa-gate-round-economics.md`
- `docs/history/findings/2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md`
- `docs/history/findings/2026-07-30-failed-qa-runs-accumulate-unreleasable-run-branches.md`
- `docs/history/findings/2026-07-30-run-termination-does-not-reach-the-acpx-child.md`
- `docs/history/findings/2026-07-30-the-autonomous-loop-orders-qa-before-its-own-preconditions.md`
- `docs/history/findings/2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md`
- `docs/history/findings/2026-08-01-characterization-corpus-is-outside-the-regeneration-command.md`
- `docs/history/findings/2026-08-02-a-spec-cycle-leaves-branches-and-worktrees-nobody-audits.md`
- `docs/history/findings/2026-08-02-a-verification-naming-a-missing-test-passes-vacuously.md`
- `docs/history/findings/2026-08-03-a-200ms-attach-budget-fails-under-ci-load.md`
- `docs/history/findings/2026-08-03-a-recording-order-defect-blocked-every-public-qa-row.md`
- `docs/history/findings/2026-08-03-gate-and-review-rounds-need-a-convergence-rule.md`
- `docs/history/findings/2026-08-04-a-doctor-next-action-that-does-not-reach-green.md`
- `docs/history/findings/2026-08-04-a-finding-does-not-prevent-the-recurrence-it-describes.md`
- `docs/history/findings/2026-08-04-a-spec-archives-with-pass-while-a-user-story-was-never-exercised.md`
- `docs/history/findings/2026-08-04-a-static-gate-row-reported-one-instance-per-cycle.md`
- `docs/history/findings/2026-08-04-an-accepted-gap-has-no-terminal-state-so-the-loop-cannot-close.md`
- `docs/history/findings/2026-08-04-branch-integrity-preflight-prescribes-a-remedy-that-reintroduces-superseded-work.md`
- `docs/history/findings/2026-08-04-fail-fast-verification-spends-the-single-repair-turn-on-the-first-of-n-defects.md`
- `docs/history/findings/2026-08-04-pre-contract-spec-graphs-run-with-no-qa-gate-and-say-nothing.md`
- `docs/history/findings/2026-08-04-review-runs-halt-autonomous-delivery-on-unrelated-dirty-files.md`
- `docs/history/findings/2026-08-04-the-documented-sandbox-escape-does-not-exist-on-the-codex-adapter.md`
- `docs/history/findings/2026-08-04-watch-derives-a-review-head-it-never-checks-is-reachable.md`
- `docs/history/findings/2026-08-04-what-still-needs-a-supervisor-between-a-prd-and-a-merge.md`
- `docs/history/findings/2026-08-05-agent-full-access-passes-config-validation-and-fails-every-task.md`
- `docs/history/findings/2026-08-05-an-absence-grep-rejects-the-work-it-was-written-to-protect.md`
- `docs/history/findings/2026-08-05-archive-refuses-a-spec-whose-graph-declined-the-qa-gate.md`
- `docs/history/findings/2026-08-05-authored-verification-gates-are-untested-code.md`
- `docs/history/findings/2026-08-05-authoring-has-no-procedure-for-a-disproven-premise.md`
- `docs/history/findings/2026-08-05-contract-seams-between-daemon-gate-and-archive.md`
- `docs/history/findings/2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md`
- `docs/history/findings/2026-08-05-o-loop-devolve-controle-por-motivos-mecanicos.md`
- `docs/history/findings/2026-08-05-preflight-prescribes-integrating-a-superseded-run-branch.md`
- `docs/history/findings/2026-08-05-review-issues-have-no-identity-across-rounds.md`
- `docs/history/findings/2026-08-05-what-a-six-spec-autonomous-session-asks-roundfix-to-change.md`
- `docs/history/findings/2026-08-05-what-roundfix-should-do-differently-measured-over-one-queue-night.md`
- `docs/history/findings/2026-08-05-what-this-repository-should-change-after-a-full-queue-night.md`
- `docs/history/findings/2026-08-06-a-promoted-backlog-entry-has-nowhere-valid-to-go.md`
- `docs/history/findings/2026-08-06-a-test-mutates-the-repository-another-test-is-reading.md`
- `docs/history/findings/2026-08-06-authoring-rules-a-release-night-made-checkable.md`
- `docs/history/findings/2026-08-06-every-run-that-failed-tonight-failed-on-a-contract.md`
- `docs/history/findings/2026-08-06-manual-thread-resolution-is-load-bearing-and-undocumented.md`
- `docs/history/findings/2026-08-06-minting-an-adr-opens-gaps-no-one-can-ever-close.md`
- `docs/history/findings/2026-08-06-rollup-agent-selection-and-execution-environments.md`
- `docs/history/findings/2026-08-06-rollup-baseline-and-derived-tooling.md`
- `docs/history/findings/2026-08-06-rollup-qa-gates-and-verification-evidence.md`
- `docs/history/findings/2026-08-06-rollup-review-and-delivery-convergence.md`
- `docs/history/findings/2026-08-06-rollup-run-lifecycle-and-branch-integrity.md`
- `docs/history/findings/2026-08-06-rollup-spec-authoring-and-contract-enforcement.md`
- `docs/history/findings/2026-08-06-six-parallel-runs-on-one-machine-show-the-seams.md`
- `docs/history/findings/2026-08-06-the-detach-tests-leak-the-process-they-prove-survives.md`
- `docs/history/findings/2026-08-06-the-gate-checks-that-adrs-were-cited-not-that-they-were-obeyed.md`
- `docs/history/findings/2026-08-06-the-loop-cannot-fix-comments-about-its-own-artifacts.md`
- `docs/history/findings/2026-08-06-the-qa-gate-and-the-pull-request-cannot-both-be-current.md`
- `docs/history/findings/2026-08-06-three-gigabytes-of-event-journal-inside-the-retention-window.md`
- `docs/history/findings/2026-08-07-a-live-contract-lives-inside-an-archived-spec.md`
- `docs/history/findings/2026-08-07-a-sixty-four-value-bound-locks-out-the-opencode-runtime.md`
- `docs/history/findings/2026-08-07-changing-the-http-contract-discards-its-exceptions.md`
- `docs/history/findings/2026-08-07-claude-agent-selections-are-never-proven.md`
- `docs/history/findings/2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md`
- `docs/history/findings/2026-08-07-the-only-gate-reports-green-on-a-red-suite.md`
- `docs/history/findings/2026-08-07-the-setup-refresh-interviews-a-repository-that-already-answered.md`
- `docs/history/findings/2026-08-07-two-http-contract-defaults-and-only-one-is-read.md`
- `docs/history/findings/2026-08-10-a-fake-adapter-goes-silent-under-a-dense-start.md`
- `docs/history/findings/2026-08-10-a-head-the-loop-did-not-push-is-a-head-nobody-reviews.md`
- `docs/history/findings/2026-08-14-a-review-retires-on-whatever-the-object-store-happens-to-hold.md`
- `docs/history/findings/2026-08-14-preflight-starves-when-the-machine-is-busy.md`
- `docs/history/findings/2026-08-26-a-qa-tasks-authored-verification-is-never-executed.md`
- `docs/history/findings/2026-08-26-the-daemon-has-no-postcondition-so-a-task-settles-clean-breaking-the-tree.md`
- `docs/history/findings/2026-09-08-run-branches-ignore-the-selected-prefix.md`
- `docs/history/findings/2026-09-16-a-qa-commit-can-carry-verifier-output.md`
- `docs/history/findings/2026-09-25-an-opt-in-jev-judgment-layer-is-proposed.md`
- `docs/history/findings/2026-09-29-the-first-full-queue-trial-needed-manual-recovery.md`
- `docs/history/findings/2026-09-30-jev-judgments-measured-against-the-spec-archive.md`
- `docs/history/findings/2026-09-30-the-v0-22-0-queue-needed-twelve-manual-interventions.md`
- `docs/history/findings/2026-10-01-the-skills-catalog-reshaped-the-baseline-inputs.md`
- `docs/history/specs/0037-terminal-outcome-integrity.md`
- `docs/history/specs/0039-review-source-evidence-and-detached-outcomes.md`
- `docs/history/specs/0040-mandatory-context-driven-adrs.md`
- `docs/history/specs/0046-public-context-driven-baseline-command.md`
- `docs/history/specs/0047-context-driven-guidance-composition.md`
- `docs/history/specs/0048-context-driven-project-decisions-and-spec-constraints.md`
- `docs/history/specs/0052-claude-adapter-standardization.md`
- `docs/history/specs/0053-qa-gate-reachability-and-verdict-semantics.md`
- `docs/history/specs/0054-tooling-task-and-verification-hygiene.md`
- `docs/history/specs/0055-owner-identity-without-fork.md`
- `docs/history/specs/0056-profiles-configure-merge-semantics.md`
- `docs/history/specs/0057-baseline-capability-evidence-and-retention.md`
- `docs/history/specs/0058-npm-trusted-publishing-and-release-preflight.md`
- `docs/history/specs/0059-run-storage-compaction-and-global-sanitation.md`
- `docs/history/specs/0060-spec-owned-reference-lifecycle.md`
- `docs/history/specs/0061-repository-derived-skill-requirements.md`
- `docs/history/specs/0062-baseline-digest-regeneration-bootstrap.md`
- `docs/history/specs/0079-one-door-for-fleet-knowledge.md`
- `docs/history/specs/0120-knowledge-lifecycle-and-durable-capture.md`
- `docs/history/specs/0121-baseline-decisions-and-complete-regeneration.md`
- `docs/history/specs/0122-verified-content-and-terminal-settlement.md`
- `docs/history/specs/0123-runtime-readiness-and-model-capabilities.md`
- `docs/history/specs/0124-verification-capacity-and-measured-economics.md`
- `docs/history/specs/0125-repository-identity-and-run-branch-policy.md`
- `docs/history/specs/0126-agent-review-before-pull-request.md`
- `docs/history/specs/0127-durable-unattended-spec-workflow.md`
- `docs/history/specs/0128-release-planning-with-bare-stable-tags.md`
- `docs/history/specs/0129-spec-authoring-and-gate-recovery.md`
- `docs/history/specs/0193-baseline-guides-that-describe-the-product-as-it-is.md`
- `docs/history/specs/0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own.md`
- `docs/history/specs/0213-a-test-suite-that-does-not-flake.md`
- `docs/history/specs/0222-a-review-through-claude-that-keeps-its-verdict.md`
- `docs/history/specs/0228-queue-items-that-stay-current-with-main.md`
- `docs/history/specs/0231-checks-that-hold-in-delivery.md`
- `docs/history/specs/0232-a-queue-that-sees-a-pull-request-merged-by-hand.md`
- `docs/history/specs/0233-a-light-tier-on-open-models.md`
- `docs/history/specs/0235-one-qa-partial-policy.md`
- `docs/history/specs/0237-review-selection-failures-that-say-why.md`
- `docs/history/specs/0239-a-glossary-every-spec-keeps-current.md`
- `docs/history/specs/0240-a-lost-rollout-is-infrastructure.md`
- `docs/history/specs/0241-retire-the-fool-autoresearch-and-council.md`
- `docs/history/specs/0252-one-gate-per-run-and-a-smaller-baseline.md`
- `docs/history/specs/0262-a-run-worktree-under-a-symlinked-root-is-cleaned.md`

## Carry-forward provenance

- Source Run: `run_20261017T065250Z_4a0bb2ff0617aa6a`
- Source commit: `f61afa85aeefd195eb8d88f7729c91f346dce202`
