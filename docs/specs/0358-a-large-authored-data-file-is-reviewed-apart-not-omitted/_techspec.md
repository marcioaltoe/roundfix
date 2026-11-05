---
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
prd: _prd.md
created: 2026-11-05
---

# A large authored data file is reviewed apart, not omitted — Technical Spec

## Executive Summary

The pre-PR review can leave a file out (`delivery.derived_paths` as
`generated`, `pre_pr_review.omit_paths` as `omit-path`) or split a large
candidate into Review Segments by first-parent commits. A large authored data
file fits neither: omitting it hides its values from the reviewer, and one
commit that adds it stays above the bound (`_prd.md` → Sources, point 5).

The design adds `pre_pr_review.data_paths` next to `omit_paths`, a reason
`data-path` in `reviewScopedDiff` that takes the files out of every code diff,
and a plan, `planReviewDataSegments` in a new
`internal/cli/review_data_segments.go`, that cuts their diff into at most four
Data Review Segments by file, hunk and line within the provider's bound.
`runConfiguredReviewSession` appends those prompts after the code prompts of
the round, in the same prepared session, and the record, the findings and the
answer file keep them apart (ADR-0368).

The trade-off accepted: a review with data files sends up to four more prompts
in its session, and a cut inside a hunk recomputes hunk headers, which Git
never wrote, so the tests prove every added line keeps the line number the
head's file gives it.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or schema version
  is minted; one Project Config key and two optional record fields are added.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0366 (a conflict-free merge of the default branch is a review delta, a sibling Spec in the same queue) touches the same review code but carries no obligation for this Spec; ADR-0368 (new) records the decision;
  it amends ADR-0282, which rejected "Splitting the diff by file or by bytes"
  for Review Segments, for declared data files only, and refines ADR-0352,
  whose `omit_paths` keeps its meaning; task_04 writes the amendment notice.
  Accounted for and unchanged, each one citation edge from a listed ADR on a part this Spec leaves alone: ADR-0092, ADR-0120, ADR-0165, ADR-0169, ADR-0192, ADR-0196, ADR-0197, ADR-0227, ADR-0233, ADR-0247, ADR-0248, ADR-0250, ADR-0251, ADR-0254, ADR-0256, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0269, ADR-0272, ADR-0277, ADR-0278, ADR-0279, ADR-0281, ADR-0286, ADR-0290, ADR-0295, ADR-0299, ADR-0306, ADR-0308, ADR-0317, ADR-0345, ADR-0348, ADR-0349. They keep the Claude budget this Spec applies to each Data Review Segment (ADR-0227), the queue's conflict regeneration (ADR-0192), and the history, archive, record, module and binary decisions the design never touches. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 2 answers of 2026-11-05 quoted in the authorization record; the
  Governed Paths are the Roundfix skill files task_04 changes; no Makefile,
  workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0358-a-large-authored-data-file-is-reviewed-apart-not-omitted/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-05 at `ad77a797`; `_prd.md` → Validation records the run.

1. `prePRReviewOverlay` (`internal/config/config.go`) decodes `provider` and
   `omit_paths` with `KnownFields(true)`, so any other key under
   `pre_pr_review` is refused; `validateReviewOmitPaths` and
   `PrePRReview.Omits` (`internal/config/pre_pr_review.go`) validate and match
   `omit_paths`, and the config parser checks the YAML node is a list of
   strings.
2. `reviewScopedDiff` (`internal/cli/review.go`) reads the Project Config at
   the declaration base with `roundconfig.DeliveryConfigAtCommit` and gives
   each changed path the first reason of `archived-spec-source`,
   `qa-evidence`, `upstream-skill`, `retired-history`, `generated`,
   `omit-path`; it diffs with `:(exclude,literal)<path>` for each.
3. `appendReviewOmissions` writes `Omitted from this diff (not reviewed): ...`
   with one count per reason and one `<path> (<reason>)` line each.
4. `runConfiguredReviewSession` builds the round's code prompt (round 1, round
   2 or Correction Review) and, above the bound, Review Segments; it keeps the
   start of the measured diff in `measuredBase`. With segment prompts it sends
   them in order through `RunPrepared`, joins the answers under `--- Review
   Segment <i> of <k> ---`, and `classifyReviewSegmentAnswers`
   (`internal/cli/review_segments.go`) numbers findings `F1...` with a
   one-based `segment`.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/config/config.go`, `internal/config/pre_pr_review.go`: the key,
  its validation and `PrePRReview.ReviewsApart`.
- `internal/cli/review.go`: the `data-path` reason, the omission line, the
  data step in `runConfiguredReviewSession`, the answer headers, the record
  fields.
- `internal/cli/review_data_segments.go` (new): the plan and its prompts.
- `internal/cli/review_segments.go`: the `dataSegment` label in
  `classifyReviewSegmentAnswers`.
- Tests: `internal/config/pre_pr_review_data_paths_test.go`,
  `internal/cli/review_data_segments_test.go` and
  `internal/cli/review_data_paths_test.go` (new).
- Documentation: the configuration and `review` guides, the Roundfix `review`
  reference, `CONTEXT.md` and ADR-0282's amendment notice.

## Implementation Design

### Interfaces

```go
// internal/config/pre_pr_review.go
func validateReviewDataPaths(paths []string) error
// ReviewsApart reports whether name is a Review Data Path (ADR-0368).
func (review PrePRReview) ReviewsApart(name string) bool

// internal/cli/review_data_segments.go
type reviewDataSegment struct {
	Paths                 []string `json:"paths"`
	DiffBytes             int      `json:"diffBytes"`
	EstimatedPromptTokens int      `json:"estimatedPromptTokens,omitempty"`
}
// reviewDataFit reports whether a prompt of promptBytes carrying diffBytes of diff fits.
type reviewDataFit func(promptBytes, diffBytes int) bool
func planReviewDataSegments(diff, from, head string, fit reviewDataFit) ([]reviewDataSegment, []string, []string, error)
```

### Data Models

- `PrePRReview` gains `DataPaths []string`; `prePRReviewOverlay` gains
  `DataPaths *[]string` with YAML name `data_paths`.
- `reviewRecord` gains `DataSegments []reviewDataSegment` with JSON name
  `dataSegments,omitempty`; `reviewFinding` gains `DataSegment int` with JSON
  name `dataSegment,omitempty`.

A record or config without the new fields reads as before.

### API Contracts

1. API Contract: the key. `pre_pr_review.data_paths` takes the path forms of
   `pre_pr_review.omit_paths` (clean repository-relative files, directories
   ending in `/`, whole-path Go `path.Match` patterns) and is refused with
   `pre_pr_review.data_paths must be a list of paths`,
   `pre_pr_review.data_paths has non-string path "<value>"`,
   `pre_pr_review.data_paths has unsafe path "<entry>"` or
   `pre_pr_review.data_paths has invalid pattern "<entry>": <error>`, each
   after the `parse config "<path>": ` prefix where the parser adds it. A
   Project Config list replaces the User Config list and `[]` clears it.
   `ReviewsApart` matches as `Omits` does.
2. API Contract: the reason. `reviewScopedDiff` gives a changed path the reason
   `data-path` when no earlier reason applies and the declarations at the
   declaration base answer `ReviewsApart`; the path leaves the diff like any
   omission. `appendReviewOmissions` leaves `data-path` paths out of the
   `Omitted from this diff (not reviewed):` counts and lines, and when at least
   one exists writes after them:

   ```text
   Reviewed apart, after this diff, in Data Review Segments: <n> path(s) of pre_pr_review.data_paths.
   <path> (data-path)
   ```

   with one `<path> (data-path)` line per path in sorted order. The record's
   `omittedPaths` never lists a `data-path` path. Without one, every prompt
   keeps its bytes.
3. API Contract: the plan. `planReviewDataSegments(diff, from, head, fit)`
   splits `diff` into file sections at each `diff --git ` line, keeps a
   section whole when it fits, else cuts it at its `@@ ` hunks, and a hunk
   that does not fit alone at line boundaries, each piece holding the
   section's header lines and a recomputed hunk header `@@ -<o>,<oc> +<n>,<nc>
   @@`, where `<oc>` and `<nc>` count the piece's lines and `<o>` and `<n>`
   are the old and new line numbers of its first line, or one less when the
   count is `0`. It packs pieces greedily in order into segments that fit,
   measuring sizes in bytes without rebuilding a string per line, and records
   each segment's paths in order. A line that fits no segment returns
   `review data too large: line <n> of <path> does not fit one data review segment`;
   more than four segments return `review data too large: <bytes> bytes of
   pre_pr_review.data_paths needs <d> data review segments; a review allows at
   most 4`. Each prompt is these two lines followed by
   `buildReviewPrompt(from, head, <segment diff>)`:

   ```text
   Data Review Segment <i> of <d>: pre_pr_review.data_paths changes <from>..<head> of the candidate; answer for this segment only.
   These files are data the project writes by hand and declares in pre_pr_review.data_paths, such as rule catalogs. Check each changed value against the Spec contexts, the code that reads it and the rest of the file; report a wrong, missing, duplicated or contradictory value.
   ```

4. API Contract: the review. In every round, after the code prompt or Review
   Segments are planned and before the readiness proof, when the round's code
   diff left out `data-path` paths, `runConfiguredReviewSession` runs `git diff
   --no-ext-diff --no-textconv --no-color <measuredBase> <head> --
   :(literal)<path>...`, plans it with the `codex` fit `diffBytes <=
   917504` or the `claude` fit `estimateReviewPromptTokens` of the prompt `<=`
   the budget, and on an error blocks before any provider call with that
   reason and the size suffix of `reviewSizeRefusalForDiff` over the data diff.
   Otherwise it appends the data prompts after the code prompts, the single
   code prompt becoming the first of the list, and sends them in order in the
   same prepared session; the first blocked answer stops the sequence, as for
   Review Segments.
5. API Contract: the record and answer. The record sets `dataSegments`, with
   `estimatedPromptTokens` only for `claude`, and keeps `segments` for code
   Review Segments only. A finding from a data prompt carries `dataSegment`
   (one-based) and no `segment`, and its anchor is checked against that
   segment's diff; a finding from a single code prompt carries no `segment`;
   finding ids stay `F1...` in send order. The answer file writes a single
   code answer under `--- Review ---`, code Review Segments under their
   current headers with `<k>` the code segment count, and each data answer
   under `--- Data Review Segment <i> of <d> ---`. Verdict reuse,
   dispositions, validation and lineage treat the data findings as findings
   of the round.
6. API Contract: the rule in writing, word for word, line wrapping free.
   - `docs/user-guide/configuration.md` (task_01): after the
     `pre_pr_review.omit_paths` row of the key table, this row:

     ```markdown
     | `pre_pr_review.data_paths` | `[]` | Reviews matching authored data files apart from the code diff, in Data Review Segments, recorded as `data-path`. Accepts the same paths as `pre_pr_review.omit_paths`; Project Config replaces the User Config list; `[]` clears it. The review reads only Project Config at its base commit, and an omission reason wins over it. |
     ```

     and after the sentence "This list affects review scope and requires no
     regeneration command." this sentence:

     ```markdown
     An authored data file too large for the review, such as a rules catalog, belongs in `pre_pr_review.data_paths`, not in `delivery.derived_paths` or `pre_pr_review.omit_paths`: the review then reads it apart instead of skipping it.
     ```

   - In `docs/user-guide/commands/review.md` and both copies of the Roundfix
     `review` reference (task_04), after the paragraph ending "round two's
     delta and the Correction Review delta.":

     ```markdown
     A changed file matched by `pre_pr_review.data_paths` in the Project Config at the candidate's base commit, and by no omission above, is a Review Data Path with reason `data-path`. It leaves every code diff, in every round, and is reviewed apart: after the omission list the code prompt adds `Reviewed apart, after this diff, in Data Review Segments: <n> path(s) of pre_pr_review.data_paths.` and one `<path> (data-path)` line each. The record's `omittedPaths` never lists it.

     The diff of the Review Data Paths, from the start of the round's code diff to the head, is sent after the code prompts, in the same session, as at most four Data Review Segments. Whole files are packed in path order while a segment fits the provider's bound, the 917504-byte diff bound for `codex` or the token budget for `claude`; a file that does not fit alone is cut at its hunks, and a hunk at its lines, each piece keeping the file header and a recomputed hunk header. Each prompt starts with `Data Review Segment <i> of <d>: pre_pr_review.data_paths changes <from>..<head> of the candidate; answer for this segment only.` The record lists `dataSegments` with their `paths`, `diffBytes` and, for `claude`, `estimatedPromptTokens`; a finding from one carries `dataSegment` and is anchored against that segment's diff; the answer file joins the answers under `--- Data Review Segment <i> of <d> ---`, after the code answer under `--- Review ---` or its Review Segment headers. A fifth segment blocks with `review data too large: <bytes> bytes of pre_pr_review.data_paths needs <d> data review segments; a review allows at most 4`, and a line that fits no segment with `review data too large: line <n> of <path> does not fit one data review segment`, each with the size suffix.
     ```

   - `CONTEXT.md` (task_04), after the **Review Omit Path** entry, each
     closing `(ADR-0368)` being the Markdown link
     `docs/adr/0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`
     as the neighbouring entries write theirs:

     ```markdown
     **Review Data Path**:
     A Project Config path in `pre_pr_review.data_paths` at the review base whose matching changed files, authored data such as rule catalogs, the pre-PR review takes out of the code diff with reason `data-path` and reviews apart in Data Review Segments, distinct from a Review Omit Path and a Derived Path Declaration (ADR-0368).
     _Avoid_: Generated path, omitted data

     **Data Review Segment**:
     One of at most four prompts, sent after a review round's code prompts in the same session, that carries only the diff of the Review Data Paths, packed by file and cut at hunks and lines, each within the provider's bound; its findings join the round's Pre-PR Review Record (ADR-0368).
     _Avoid_: Review Segment, data summary
     ```

   - ADR-0282 (task_04): a new last paragraph "**Amended by ADR-0368
     (2026-11-05).** Files declared in `pre_pr_review.data_paths` leave the
     Review Segments and are reviewed apart in Data Review Segments, which
     split their diff by file, hunk and line; Review Segments stay commit
     ranges.", where `ADR-0368` links
     `0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`;
     at the end of the Rejected bullet "**Splitting the diff by file or by
     bytes.**", the sentence "ADR-0368 takes it for declared data files
     only."; and `updated_at: 2026-11-05T12:00:00Z`.

### Surface Transcripts

None. `roundfix review` writes its record as JSON to stdout and a provider
answer is needed for any run past the plan; the tests of Testing Approach
drive it through `runCLI` with a fake reviewer, and the only new refusals are
asserted there. The config refusal prints through the existing `parse config`
path.

### Invariants

```text
1. A candidate whose base commit declares no pre_pr_review.data_paths is reviewed with the same prompts, record and answer file bytes as before.
2. A path omitted for any existing reason stays omitted with that reason, whatever pre_pr_review.data_paths says.
3. No code prompt of a candidate with Review Data Paths holds a line of those files, in any round.
4. Every added or context line of a Data Review Segment sits, by its hunk header, at the line where the head's file holds that text.
5. delivery.derived_paths and pre_pr_review.omit_paths keep their meaning, and conflict regeneration reads only delivery.derived_paths.
```

### Declared breaks

- A Project Config with `pre_pr_review.data_paths` is accepted; before, it was
  refused as an unknown field.
- A review of a candidate with Review Data Paths sends more prompts and uses
  the prepared-session path even when its code diff is one prompt; its answer
  file then writes the code answer under `--- Review ---`.
- `classifyReviewSegmentAnswers` gains the parameters it needs to label data
  findings; its callers change with it. No existing test changes, and none is
  renamed or removed, so the coverage record is not re-recorded.

## Coverage Map

- Goal 1, Core Feature 1, User Story 1 → API Contract 1,
  `TestPrePRReviewDataPathsLoadsAndValidates`, Invariant 5.
- Goal 2, Core Features 2 and 4, User Stories 1 and 2, Success Metric 1 → API
  Contracts 2, 4 and 5, `TestReviewSendsADataPathApart`,
  `TestReviewReadsADataPathApartInRoundTwo`, Invariants 1 and 3.
- Goal 3, Core Feature 3, Success Metrics 2 and 3 → API Contract 3,
  `TestReviewDataSegmentsPackCutAndBound`,
  `TestReviewSplitsALargeDataFileAtLineBoundaries`,
  `TestReviewBlocksDataNeedingMoreThanFourDataSegments`, Invariant 4.
- Goal 4, Success Metric 4 → API Contract 2,
  `TestReviewOmissionWinsOverADataPath`, the existing review, scope and
  segment suites, Invariants 1, 2 and 5.
- Goal 5, Core Feature 5, Success Metric 5 → API Contract 6.

Deviations from the PRD: None.

## Absence Evidence

1. `pre_pr_review.data_paths`, `DataPaths` and `ReviewsApart` — command:
   `grep -rn 'data_paths\|DataPaths\|ReviewsApart' internal`; result: no
   match.
2. The `data-path` reason and Data Review Segments — command:
   `grep -rn 'data-path\|dataSegment\|Data Review Segment' internal docs/user-guide .agents/skills/roundfix CONTEXT.md`;
   result: no match.
3. The new files — command: `ls internal/cli/review_data_segments.go internal/cli/review_data_segments_test.go internal/cli/review_data_paths_test.go internal/config/pre_pr_review_data_paths_test.go`;
   result: no such files.

## ADR Alternatives

- ADR-0366: it rejects a tree-equality test of the merge, `--remerge-diff`, a larger bound, history rewriting and a delta across a merge of another branch; this design takes none of them.
- ADR-0368: the decision this Spec records; it rejects a structured summary,
  a larger bound, reusing Review Segments and `.gitattributes`, and the design
  takes none of them.
- ADR-0282: amended; it rejected splitting the diff by file or by bytes for
  Review Segments, and this design takes that option for declared data files
  only, cutting by file, hunk and line, never by bytes inside a line; Review
  Segments stay commit ranges and the bound stays.
- ADR-0352: refined; `omit_paths` keeps its meaning and its precedence over
  the new reason, and the design does not take its rejected automatic
  omission or `linguist-generated`.
- ADR-0227: unchanged; each Claude Data Review Segment fits the same budget,
  and none of its rejected options is taken.
- ADR-0192: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0278: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0281: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0290: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0299: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0092: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0120: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0233: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0247: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0248: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0250: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0254: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0256: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0258: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0260: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0261: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0277: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0295: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0308: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0345: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0165: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0251: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0262: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0269: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0272: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0348: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0349: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0169: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0279: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0286: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0306: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0196: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0197: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0317: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.

## Integration Points

None external. The `git` binary is a Base Test Tool; the tests use it to make
the repositories and as the oracle of Invariant 4 (`git show
<head>:<path>`).

## Testing Approach

All tests are hermetic: `t.TempDir()` repositories through
`newReviewCommandFixture`, `commitScopeFile`, `gittest.InitRepo`,
`gittest.Run`, the fake reviewer `segmentReviewRunner`, and `runCLI`. Every new
top-level test calls `t.Parallel()` first, and each new file opens with the
`Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines. A
review's base commit holds the Project Config with the key, committed and
made `main` before the candidate commits, because the review reads it there.

- `internal/config/pre_pr_review_data_paths_test.go` (task_01):
  `TestPrePRReviewDataPathsLoadsAndValidates`, through
  `ResolveConfigProposal` only: rows `absent`, `inherits user`, `project
  replaces user`, `empty clears user`; each refusal of API Contract 1 (`a
  mapping`, `a number entry`, `an absolute path`, `a parent segment`, `a bad
  pattern`); and `ReviewsApart` on an exact file, a directory and a pattern.
- `internal/cli/review_data_segments_test.go` (task_02):
  `TestReviewDataSegmentsPackCutAndBound` on diffs Git writes, with a byte
  fit of a few kilobytes: `small files share one segment` (one segment, both
  paths, the first header line), `a new file is cut at lines` (400 numbered
  lines; every piece within the fit, every added line at its line in `git
  show`), `a modified file is cut at hunks and lines` (600 lines with one
  change and 200 inserted lines; every added and context line at its line),
  `a fifth segment is refused` (the refusal text) and `a line that fits no
  segment is refused`.
- `internal/cli/review_data_paths_test.go` (task_03):
  - `TestReviewSendsADataPathApart`, subtests `codex` and `claude`: base
    declares `data_paths: ['rules/*.json']`; the candidate adds `code.go` and
    a 9,000-entry `rules/icms.json`; two prompts, the code prompt without the
    catalog and with the API Contract 2 lines, the data prompt starting with
    the API Contract 3 header; `dataSegments` of one with that path, no
    `omittedPaths`, no `segments`; the data finding with `dataSegment: 1`, no
    `segment`, standing; the answer file headers of API Contract 5.
  - `TestReviewSplitsALargeDataFileAtLineBoundaries`: a 24,000-entry file
    under `data_paths: ['rules/']` with `codex`; two data segments within
    917,504 bytes; a finding on line 20,000 stands with `dataSegment: 2`.
  - `TestReviewBlocksDataNeedingMoreThanFourDataSegments`: five 14,000-entry
    files; exit `2`, outcome `blocked`, the refusal of API Contract 3 with the
    size suffix, and no session prepared.
  - `TestReviewOmissionWinsOverADataPath`: a file in both `omit_paths` and
    `data_paths` is omitted as `omit-path`, with no data segment.
  - `TestReviewReadsADataPathApartInRoundTwo`: round 1 findings on `code.go`,
    then a commit changing one entry of the catalog; round 2 sends the round-2
    delta without the catalog and one Data Review Segment holding only the
    changed entry.
- Existing suites: `./internal/cli` and `./internal/config` whole, `./skills`,
  `-tags docscontract ./internal/docscontract`.
- Sabotages, each run by the author on the prototype, with the case whose
  input makes it observable:
  - `ReviewsApart` always false: `TestReviewSendsADataPathApart` fails on
    both subtests (one prompt, catalog in the code prompt);
  - the `data-path` reason before `omit-path`: `TestReviewOmissionWinsOverADataPath` fails;
  - the line cursor not advanced between pieces: `a new file is cut at lines`
    fails on the second piece's first line;
  - the segment limit unchecked: `a fifth segment is refused` and
    `TestReviewBlocksDataNeedingMoreThanFourDataSegments` fail;
  - `appendReviewOmissions` listing `data-path` as omitted:
    `TestReviewSendsADataPathApart` fails on the code prompt lines.

## Build Order

1. The key, its validation and `ReviewsApart`, with
   `pre_pr_review_data_paths_test.go`, and the configuration guide row and
   sentence with the Behavior Surface Record re-recorded last (task_01).
2. The `data-path` reason, the omission line and `planReviewDataSegments`,
   with `review_data_segments_test.go` (task_02) (depends on: 1).
3. The data step in `runConfiguredReviewSession`, the record fields, the
   finding label and the answer headers, with `review_data_paths_test.go`
   (task_03) (depends on: 2).
4. The `review` guide and reference, `CONTEXT.md`, ADR-0282's amendment, the
   skill version and the Behavior Surface Record recorded last (task_04)
   (depends on: 3).
5. QA gate (task_05) (depends on: 4).

task_02 reads `ReviewsApart` from task_01 and task_01's guide gives the CLI
files of tasks 02 and 03 a guide in their dependency closure; task_03 sends
what task_02 plans. No Verification depends on another Spec. Spec 0356, which
the queue delivers first, changes `planReviewSegments`, the lineage and the
round-2 delta of the same files; this Spec's edits sit in other functions or
after them, and its data diff starts from whatever `measuredBase` holds.

## Risks & Considerations

- A data file with very long lines can still exceed the bound with a single
  line; the refusal names the line and the path.
- A recomputed hunk header is Roundfix's, not Git's; Invariant 4 and its test
  hold it to the head's file.
- Up to four more prompts lengthen a review session; the runtime compacts
  earlier turns as it does for Review Segments.
- An adopter that keeps a catalog in `delivery.derived_paths` keeps today's
  omission; the guide tells it to move the entry.

## Glossary

- adds: **Review Data Path**
- adds: **Data Review Segment**

## Decisions

- The key under `pre_pr_review`, with the `omit_paths` grammar and
  base-commit reading; no new package.
- Data Review Segments in every round, from the round's `measuredBase`.
- Cut by file, then hunk, then line, with recomputed headers; never inside a
  line.
- The four-segment limit and the size suffix reused from Review Segments.
- ADR-0368 records the decision, amends ADR-0282 and refines ADR-0352.
- Secondbrain and Exa as `_prd.md` → Validation records.
- Four implementation Tasks in series and a QA gate.
