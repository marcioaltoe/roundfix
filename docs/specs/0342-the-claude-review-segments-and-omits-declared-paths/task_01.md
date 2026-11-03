---
task: task_01
spec: 0342-the-claude-review-segments-and-omits-declared-paths
status: pending
type: backend
complexity: medium
---

# Task 01: A Review Omit Path leaves a declared file out of the pre-PR review

## Overview

Project Config gains `pre_pr_review.omit_paths`. The pre-PR review reads it
from the Project Config at the candidate's base commit and leaves every
changed file it matches out of the diff of every round, recorded in
`omittedPaths` with reason `omit-path`, without a derived-path declaration or
a regeneration command. A review whose base has no key sends the same prompt
bytes as before. This answers the adopted Backlog Entry of 2026-11-03
([references](references/2026-11-03-the-claude-pre-pr-review-blocks-on-size-instead-of-segmenting.md)),
reported by vortex, whose second suggestion asks for this key.

## Requirements

1. MUST add `OmitPaths []string` to `PrePRReview` in
   `internal/config/config.go`, decoded from `pre_pr_review.omit_paths`;
   Project Config replaces the User Config list, and an explicit `[]` clears
   it, as `delivery.derived_paths` does. The provider and its `Source` keep
   their behavior.
2. MUST create `internal/config/pre_pr_review.go` with
   `validateReviewOmitPaths` and `PrePRReview.Omits` exactly as
   `_techspec.md` → Interfaces and API Contract 1 state: entries are checked
   by the existing `validateDerivedPathPatterns` with the key
   `pre_pr_review.omit_paths`; `Omits` matches an exact path, a directory
   entry ending in `/` as a prefix, or a whole-path `path.Match` pattern.
   `Validate` calls `validateReviewOmitPaths`. A non-list value is refused
   with `pre_pr_review.omit_paths must be a list of paths` before typed
   decoding, wrapped by the existing `parse config` label.
3. MUST add to the Project scope of the `roundfix init` template
   (`defaultConfigYAML`), and only there, the block
   `pre_pr_review:` / `# Paths the pre-PR review leaves out, read from the Project Config at the review base.` / `omit_paths: []`.
   The User template keeps its bytes.
4. MUST, in `reviewScopedDiff` (`internal/cli/review.go`), give a changed
   path that no earlier reason claims and that
   `declarations.PrePRReview.Omits` matches the reason `omit-path`, using the
   same `DeliveryConfigAtCommit` read at the declaration base, so a key only
   in the head or a segment head omits nothing (API Contract 2). The read's
   error text stays `review scope: read derived paths: <error>`.
5. MUST, in `appendReviewOmissions`, end the summary sentence with
   `, <o> of pre_pr_review.omit_paths.` after `<g> of generated files` only
   when `<o>` is not zero; with no `omit-path` omission the sentence and the
   whole prompt keep their bytes (API Contract 3, Invariants 1 and 5).
6. MUST create `internal/config/pre_pr_review_omit_paths_test.go` with
   `TestPrePRReviewOmitPathsLoadsAndValidates` and
   `TestProjectConfigTemplateDescribesOmitPaths`, and
   `internal/cli/review_omit_paths_test.go` with
   `TestReviewOmitsConfiguredPathsFromTheBaseCommit` and
   `TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths`, covering the
   cases `_techspec.md` → Testing Approach lists. The summary test compares
   the line with the literal sentence `Omitted from this diff (not reviewed):
   0 path(s) of QA evidence, 0 of upstream-managed skills, 0 of archived Spec
   sources, 0 of retired history, 0 of generated files.` Every new top-level
   test and subtest calls `t.Parallel()` first.
7. MUST, in `docs/user-guide/configuration.md`, add a
   `pre_pr_review.omit_paths` row to the Review and Run settings table
   (default `[]`; Project Config replaces User Config; read only from the
   Project Config at the review's base commit, so an item sees a new value
   once its merge base includes it) and, after the Delivery conflict recovery
   section's derived-path text, one paragraph containing the sentence "A file
   the review should skip but nobody regenerates belongs in
   pre_pr_review.omit_paths, not in delivery.derived_paths", with the key and
   the derived-path key written as inline code.
8. MUST re-record the Behavior Surface Record as the LAST step, after the
   last guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   and commit `docs/references/behavior-surfaces.json`.
9. MUST NOT change conflict regeneration, the `delivery.derived_paths`
   grammar or matching, any other guide, any skill or ADR, or any existing
   test.

## Subtasks

- [ ] Add the key, its validation, its matcher and the Project template block.
- [ ] Omit matching paths at the base commit with reason `omit-path` and extend the summary only when needed.
- [ ] Add the config and review tests.
- [ ] Document the key in the configuration guide and re-record the fingerprints last.

## Acceptance Criteria

- [ ] A base-commit glob omits a matching file with reason `omit-path` in round 1 and round 2, a head-only key omits nothing, and a file a derived entry also matches stays `generated`.
- [ ] Unsafe, malformed and non-list values are refused naming the key and the entry.
- [ ] Without the key, the omission summary and the prompt keep their bytes, and the existing scope tests pass unchanged.
- [ ] The configuration guide describes the key, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- instruction: `docs/adr/0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md`
- instruction: `internal/config/delivery.go`
- instruction: `internal/cli/review_scope_test.go`
- interface: `internal/config/config.go`
- creates: `internal/config/pre_pr_review.go`
- creates: `internal/config/pre_pr_review_omit_paths_test.go`
- interface: `internal/cli/review.go`
- creates: `internal/cli/review_omit_paths_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `PrePRReview` config value and the review's `omittedPaths` reasons — consumers: `./internal/config`, `./internal/cli`
- The configuration guide's Behavior Surface fingerprint — consumers: `./internal/docscontract`

## Verification

- `grep -q 'func TestPrePRReviewOmitPathsLoadsAndValidates(' internal/config/pre_pr_review_omit_paths_test.go && grep -q 'func TestReviewOmitsConfiguredPathsFromTheBaseCommit(' internal/cli/review_omit_paths_test.go && grep -q 'func TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths(' internal/cli/review_omit_paths_test.go && go test -count=1 -timeout 20m ./internal/config -run '^(TestPrePRReviewOmitPathsLoadsAndValidates|TestProjectConfigTemplateDescribesOmitPaths)$' && go test -count=1 -timeout 20m ./internal/cli -run '^(TestReviewOmitsConfiguredPathsFromTheBaseCommit|TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths|TestReviewOmitsRetiredHistoryAndGeneratedFiles|TestReviewReadsDerivedPathsFromTheBaseCommit|TestReviewKeepsDeletedArchiveRecordsAndOmissionPrecedence|TestReviewScopesRoundTwoDelta|TestReviewOmitsQAEvidenceAndUpstreamSkills)$' && tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'but nobody regenerates belongs in' && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task the new test files do not exist, so the first check fails; afterwards the new tests pass, the existing scope tests still pass unchanged, the configuration guide carries the key and the Behavior Surface Record is current.

## References

- `_prd.md` → Goals; Core Feature 1; Success Metric 1; Success Metric 2; Decisions
- `_techspec.md` → API Contracts 1-3; Invariants 1, 4 and 5; Testing Approach; Build Order 1
- ADR-0352; ADR-0282
