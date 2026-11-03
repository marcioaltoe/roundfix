---
task: task_01
spec: 0342-the-claude-review-segments-and-omits-declared-paths
status: completed
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
3. MUST NOT add the key to either `roundfix init` template
   (`defaultConfigYAML`): `TestDefaultConfigSelectsNoPrePRReviewProvider`
   (`internal/config/review_source_scope_test.go`, #274) keeps every
   generated config free of `pre_pr_review`, and that guard keeps its bytes.
   The configuration guide is where the key is described. Both templates
   keep their bytes.
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
   `TestPrePRReviewOmitPathsLoadsAndValidates`, and
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

- [ ] Add the key, its validation and its matcher; leave both init templates unchanged.
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

- `grep -q 'func TestPrePRReviewOmitPathsLoadsAndValidates(' internal/config/pre_pr_review_omit_paths_test.go && grep -q 'func TestReviewOmitsConfiguredPathsFromTheBaseCommit(' internal/cli/review_omit_paths_test.go && grep -q 'func TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths(' internal/cli/review_omit_paths_test.go && go test -count=1 -timeout 20m ./internal/config -run '^(TestPrePRReviewOmitPathsLoadsAndValidates|TestDefaultConfigSelectsNoPrePRReviewProvider)$' && go test -count=1 -timeout 20m ./internal/cli -run '^(TestReviewOmitsConfiguredPathsFromTheBaseCommit|TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths|TestReviewOmitsRetiredHistoryAndGeneratedFiles|TestReviewReadsDerivedPathsFromTheBaseCommit|TestReviewKeepsDeletedArchiveRecordsAndOmissionPrecedence|TestReviewScopesRoundTwoDelta|TestReviewOmitsQAEvidenceAndUpstreamSkills)$' && tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'but nobody regenerates belongs in' && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task the new test files do not exist, so the first check fails; afterwards the new tests pass, the existing scope tests still pass unchanged, the configuration guide carries the key and the Behavior Surface Record is current.

## References

- `_prd.md` → Goals; Core Feature 1; Success Metric 1; Success Metric 2; Decisions
- `_techspec.md` → API Contracts 1-3; Invariants 1, 4 and 5; Testing Approach; Build Order 1
- ADR-0352; ADR-0282

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, agentic-cli-design, systematic-debugging.

Task status remains Daemon-owned. Implementation and focused evidence only;
declared Verification and repository Verification remain for the Daemon.

Added `PrePRReview.OmitPaths` with list replacement and explicit-empty clearing,
shape validation before typed decoding, labeled path validation through the
existing derived-path validator, and exact/directory/whole-path matching.
Provider selection and its Source are unchanged. Both init templates retain
their bytes. The scoped diff uses its existing immutable Project Config read
and assigns `omit-path` only after every existing omission reason. The summary
adds the new count only when it is nonzero.

The configuration guide documents the key and its base-commit adoption rule,
and distinguishes review omissions from conflict regeneration. After the last
source and guide edit, the sanctioned recorder updated only the `config keys`
and configuration-guide fingerprints in the Behavior Surface Record.

Acceptance evidence:

- Base glob, head-only key, and precedence: the five parallel cases of
  `TestReviewOmitsConfiguredPathsFromTheBaseCommit` passed. They assert the
  provider's diff body, omission reason, exact summary, diff byte count, and
  persisted record. Round 2 still uses the original base even after the head
  clears its list; a head-only list still omits nothing in round 2. A path
  matched by both keys retains `generated`.
- Invalid configuration: `TestPrePRReviewOmitPathsLoadsAndValidates` passed
  proposal parsing and disk loading, Project/User precedence, empty clearing,
  provider/Source preservation, exact/glob/directory matching, unsafe paths,
  malformed patterns, non-list values, and non-string list entries. Refusals
  retain the config label and key; unsafe/pattern refusals name the entry.
  Direct `Validate` checks reject the same unsafe and malformed entries.
- Byte compatibility: the new summary test compares the required literal
  sentence and the command's whole prompt. The fixed-input whole-prompt
  SHA-256 was identical before and after production edits:
  `825ac238f6467063b26ed28c9f2fcfeaea0e2dbde74cbc3f99991e674184f2f4`.
  The unchanged existing scope and init-template guards passed. A source-byte
  comparison also confirmed `defaultConfigYAML` is unchanged for both scopes.
- Guide and record: inspection confirmed the required table row and exact
  distinction sentence. The sanctioned record command below exited 0 after
  all implementation and guide edits, checking the record it wrote.

Focused checks (each exited 0 on the final implementation):

```sh
go test ./internal/config -run '^(TestPrePRReviewOmitPathsLoadsAndValidates|TestDefaultConfigSelectsNoPrePRReviewProvider|TestDefaultConfigScopesReviewSourceToPullRequestFeedback|TestProjectConfigReadsDerivedPathDeclarations|TestDerivedPathDeclarationsRefuseUnsafeEntries)$' -count=1 -timeout 3m -v
go test ./internal/cli -run '^TestReview(OmitsConfiguredPathsFromTheBaseCommit|OmissionSummaryKeepsItsBytesWithoutOmitPaths|OmitsRetiredHistoryAndGeneratedFiles|ReadsDerivedPathsFromTheBaseCommit|KeepsDeletedArchiveRecordsAndOmissionPrecedence|ScopesRoundTwoDelta|OmitsQAEvidenceAndUpstreamSkills|BlocksOnUnreadableBaseDerivedDeclarations|KeepsTheQAReportInTheDiff)$' -count=1 -timeout 3m -v
go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
git diff --check
```

The first pre-implementation config run failed on the absent field and matcher;
the first review run rejected the unknown key, while its byte-compatibility
case passed. During test development, the new round-2 fixture initially named
an anchor already in its base; validation correctly dismissed that finding.
The fixture now anchors `src/a.go:1`, which its candidate changes; both affected
round-2 cases and the final focused review set passed.

No existing test, other Task, Task Graph, skill, ADR, derived-path grammar or
conflict-regeneration code was changed. No commit, push, or Pull Request was
created. Follow-up scope (Claude segmentation, size refusals, and the remaining
review documentation) stays with the later Tasks.
