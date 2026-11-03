---
spec: 0342-the-claude-review-segments-and-omits-declared-paths
prd: _prd.md
created: 2026-11-03
---

# The Claude review segments a large candidate and omits declared paths — Technical Spec

## Executive Summary

Three small changes to the pre-PR review and one to the Delivery Queue answer
the Vortex report (`_prd.md` → Sources). Project Config gains
`pre_pr_review.omit_paths`, read from the base commit like the derived-path
declarations and applied in the one scoped-diff function every round already
uses. The Review Segment planner, today bounded only by the Codex byte bound,
takes a fit rule, so a Claude round-1 prompt above its budget is planned with
the token estimate of each segment's prompt as sent. Every size refusal keeps
its text and gains a suffix naming the three largest files of the diff it
measured and the key; the queue parks a blocked review with the record's
reason line. The trade-off accepted is that Claude segments share one session,
as Codex segments do, relying on the runtime's compaction instead of opening a
session per segment, which would break the one-session Reviewer Lineage.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; one optional
  config list, one omission reason value and one optional record field.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git reads and the existing
  Agent Session only. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0227 sets the Claude token budget,
  ADR-0227: "The Claude review prompt is bounded in tokens, not bytes"; it now
  bounds each Claude Review Segment. ADR-0282 owns Review Segments and the
  base-commit rule for review omissions, ADR-0282: "the candidate's own
  Project Config is not used, because a candidate must not choose what its
  review skips"; it is amended so Claude segments too. ADR-0352 is new. The
  accepted ADRs one citation edge from those are listed as in the PRD: the
  review lineage and record decisions ADR-0165, ADR-0169, ADR-0196, ADR-0197, ADR-0233, ADR-0272, ADR-0279, ADR-0286 and ADR-0317; the history and
  archive decisions ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0251, ADR-0254, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0269, ADR-0277, ADR-0278, ADR-0295, ADR-0299, ADR-0306, ADR-0308 and ADR-0348; and ADR-0250, ADR-0256, ADR-0281, ADR-0290, ADR-0345 and ADR-0349. ADR
  Alternatives relates each. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AW answers of 2026-11-03 quoted in the authorization record; the Governed
  Paths are the Roundfix skill's `SKILL.md` pair and its canonical review and
  deliver references; no Makefile, workflow, module or lint configuration
  changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0342-the-claude-review-segments-and-omits-declared-paths/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-03 at `4e61bf6d`; `_prd.md` → Validation records the runs.

- `reviewScopedDiff` (`internal/cli/review.go`) reads Project Config at the
  declaration base with `roundconfig.DeliveryConfigAtCommit` (user path
  empty), assigns each changed path the first matching reason among
  `archived-spec-source`, `qa-evidence`, `upstream-skill`, `retired-history`
  and `generated` (a non-directory `delivery.derived_paths` entry), and runs
  `git diff --no-ext-diff --no-textconv --no-color <base> <head> --` with one
  `:(exclude,literal)` pathspec per omitted path. Rename detection follows
  Git's default.
- `appendReviewOmissions` always writes `Omitted from this diff (not
  reviewed): <q> path(s) of QA evidence, <u> of upstream-managed skills, <a>
  of archived Spec sources, <r> of retired history, <g> of generated files.`
  followed by one `<path> (<reason>)` line per omission.
- `runConfiguredReviewSession` enters segmentation only when the provider is
  `codex` or the lineage is a Correction Review, only in round 1, and only
  when the scoped diff exceeds `reviewDiffBound` (917,504 bytes); otherwise it
  applies `checkReviewDiffBound` in that branch. Spec contexts are appended
  afterwards, to the first segment prompt only. For `claude` it then refuses a
  prompt whose `estimateReviewPromptTokens` (two bytes per token) exceeds
  `claudeReviewPromptBudget` (500,000) with `review prompt too large: <n>
  estimated tokens after omitting <k> path(s) exceeds the claude review budget
  of 500000 tokens, half of its 1000000-token context window`.
- `planReviewSegments` (`internal/cli/review_segments.go`) walks
  `rev-list --first-parent --reverse`, refuses with the unsegmented byte
  message when one commit's own scoped diff exceeds the bound, grows each
  segment greedily, refuses a fifth segment with `review diff too large: <n>
  bytes after omitting <k> path(s) needs <s> review segments; a review allows
  at most 4`, and prefixes each prompt with `Review Segment <i> of <s>:
  commits <b>..<h> of the candidate <base>..<head>; answer for this segment
  only.`. `reviewSegment` holds `base`, `head` and `diffBytes`.
- `Engine.reviewCandidate` (`internal/delivery/engine.go`) parks every
  `ReviewOutcomeBlocked` as the bare `review-blocked`, although
  `deliveryReviewResult` (`internal/cli/deliver_workflow.go`) already passes
  the record's reason in `ReviewResult.Reason`. `ClassifyPark` reads the text
  before the first colon.

## System Architecture

No new package. `internal/config` gains the key, its validation and a
matcher; `internal/cli` changes the scoped diff, the segment planner, the
Claude admission and every size refusal through one helper file;
`internal/delivery` changes one park call. The guides, the Roundfix skill
references, `CONTEXT.md` and the ADRs describe them.

## Implementation Design

### Interfaces

```go
// internal/config/config.go
type PrePRReview struct {
	Provider  string
	Source    string
	OmitPaths []string // Project Config replaces User Config, as derived_paths
}

// internal/config/pre_pr_review.go
func validateReviewOmitPaths(paths []string) error // "pre_pr_review.omit_paths has unsafe path %q" | "... has invalid pattern %q: %w"
func (review PrePRReview) Omits(name string) bool  // exact, directory prefix ending "/", whole-path path.Match

// internal/cli/review_segments.go
type reviewSegment struct {
	Base                  string `json:"base"`
	Head                  string `json:"head"`
	DiffBytes             int    `json:"diffBytes"`
	EstimatedPromptTokens int    `json:"estimatedPromptTokens,omitempty"` // claude only
}
type reviewSegmentFit func(index int, prompt string, diff string) bool

// internal/cli/review_size_refusal.go
type reviewDiffFile struct{ Path string; Bytes int }
func reviewLargestDiffFiles(ctx context.Context, git preflight.GitRunner, root, base, head string, omitted []reviewOmittedPath, diff string) ([]reviewDiffFile, error)
func reviewSizeRefusal(reason string, files []reviewDiffFile) error // reviewPrePromptError
```

`planReviewSegments` keeps its walk and gains the fit rule: `codex` passes
`len(diff) <= reviewDiffBound`, byte for byte today's tests; `claude` passes
`estimateReviewPromptTokens(prompt) <= claudeReviewPromptBudget`, where
`prompt` is the segment prompt as it will be sent: its header, written for
segment `<i> of 4` (both numbers are one digit in any admitted plan, so its
length equals the sent header's), `buildReviewPrompt` of the segment, and for
the first segment the rendered Spec contexts.

### Data Models

- Project Config: `pre_pr_review.omit_paths: [<entry>, …]`; absent means
  `[]`. An explicit `[]` in Project Config clears a User Config list.
- Pre-PR Review Record: `omittedPaths[].reason` may be `omit-path`; a Claude
  segment carries `estimatedPromptTokens`; a Codex segment never does, so its
  record bytes are unchanged. A segmented Claude record's
  `estimatedPromptTokens` is the sum of its segments'.

### API Contracts

1. API Contract: `pre_pr_review.omit_paths` accepts a YAML list of strings;
   each entry follows `validateDerivedPathPatterns`. A non-list value is
   refused with `pre_pr_review.omit_paths must be a list of paths`; an unsafe
   entry with `pre_pr_review.omit_paths has unsafe path "<entry>"`; a bad
   pattern with `pre_pr_review.omit_paths has invalid pattern "<entry>":
   <error>`. Each error is wrapped by the existing `parse config "<label>"`.
2. API Contract: the review reads the key only from the Project Config at the
   declaration base (`DeliveryConfigAtCommit` with an empty user path, the
   same read as the derived declarations, keeping the error text `review
   scope: read derived paths: <error>`). A changed path no earlier reason
   claims and `Omits` matches gets reason `omit-path`, in round 1, round 2,
   the Correction Review and every segment.
3. API Contract: when at least one path has reason `omit-path`, the summary
   sentence ends `…, <g> of generated files, <o> of pre_pr_review.omit_paths.`
   instead of `…, <g> of generated files.`; with none, it is byte-identical
   to today's.
4. API Contract: a `claude` round-1 prompt (with Spec contexts) whose
   estimate exceeds the budget is planned into Review Segments with the
   Claude fit rule and run as Codex segments are run. One commit whose own
   segment prompt does not fit refuses with today's Claude message, the
   estimate being the whole unsegmented prompt's. More than four segments
   refuses with `review prompt too large: <n> estimated tokens after omitting
   <k> path(s) needs <s> review segments; a review allows at most 4`.
5. API Contract: round 2 and the Correction Review are never segmented; a
   `claude` round-2 prompt above the budget refuses with today's Claude
   message, and a Correction Review keeps the byte bound for either provider.
6. API Contract: every size refusal of API Contracts 4–5 and of the Codex
   byte bound and segment limit keeps its text and appends
   `; largest files in the diff: <path> (<bytes> bytes)[, <path> (<bytes>
   bytes)[, <path> (<bytes> bytes)]]; pre_pr_review.omit_paths in the base
   commit's Project Config leaves generated files out of the review`. The
   files are those of the diff the refusal measured (round 1's whole scoped
   candidate, or the round-2 or Correction Review delta), ranked by the bytes
   of each file's section, descending, ties by path ascending. A section is
   the text from one line starting `diff --git ` to the next; its path is the
   entry at the same position of `git diff --name-only -z` run with the
   measured diff's own arguments, so a quoted or renamed path is written as
   Git lists it, never parsed from the header.
7. API Contract: `Engine.reviewCandidate` parks `ReviewOutcomeBlocked` as
   `review-blocked: <first line of ReviewResult.Reason, trimmed>` when that
   line is not empty, else `review-blocked`. `ClassifyPark` still answers
   `review`.

### Surface Transcripts

None. No command, flag or exit code is added; the changed refusal lines are
given exactly in API Contracts 4 and 6, and the Task tests assert each whole
standard-error line `roundfix: review blocked: <reason>` and exit `2`.

### Invariants

1. A review whose base Project Config has no `pre_pr_review.omit_paths` sends
   the same prompt bytes and writes the same record as before.
2. A Codex candidate's segments, segment prompts and record bytes are
   unchanged.
3. Every Claude segment prompt as sent is estimated within 500,000 tokens.
4. The key is never read from the candidate's head or a segment head.
5. A path an existing rule omits keeps its existing reason.
6. The text of every size refusal before its new suffix is unchanged.

### Declared breaks

- Every size refusal gains the suffix of API Contract 6, Codex included. The
  exact-reason pins change: `TestReviewCommandRefusesAClaudePromptOverHalfItsWindow`
  and `TestReviewCommandKeepsTheCodexDiffByteBound`
  (`internal/cli/review_prompt_bound_test.go`),
  `TestReviewBlocksACandidateNeedingMoreThanFourSegments`
  (`internal/cli/review_segment_test.go`) and
  `TestReviewBlocksAboveTheBoundWithoutAProviderCall`
  (`internal/cli/review_scope_test.go`) compare with the old text plus the
  suffix.
- A `claude` round-1 prompt above the budget whose commits fit separately is
  reviewed instead of refused.
- A blocked review parks with its reason after the colon; the blocker
  expectation of `TestDeliveryEngineParksDeclaredBlockersAndContinues` in
  `internal/delivery/engine_test.go` for the item `review-blocked`, whose
  reason is `review unavailable`, becomes `review-blocked: review
  unavailable`.

## Coverage Map

- Goal 1, Core Feature 1, Success Metrics 1–2 → `PrePRReview.OmitPaths`,
  `validateReviewOmitPaths`, `PrePRReview.Omits`, `reviewScopedDiff`,
  `appendReviewOmissions` (Build Order 1).
- Goal 2, Core Feature 2, Success Metrics 3–4 → `planReviewSegments` fit rule,
  `runConfiguredReviewSession` Claude admission, `reviewSegment` (Build
  Order 2).
- Goal 3, Core Features 3–4, Success Metrics 5–6 → `reviewLargestDiffFiles`,
  `reviewSizeRefusal`, `Engine.reviewCandidate` (Build Order 3).
- Goal 4 → the zero-count branch of `appendReviewOmissions` keeps the prompt
  bytes (Invariant 1), the Codex fit rule of `planReviewSegments` keeps the
  Codex segments (Invariant 2), `reviewSizeRefusal` keeps every refusal's
  text before its suffix (Invariant 6), and the round-2 and Correction Review
  branches of `runConfiguredReviewSession` stay unsegmented (Build Orders
  1–3).
- Core Feature 5, Success Metric 7 → the configuration guide (Build Order 1);
  the review and deliver guides, references, `CONTEXT.md`, ADR-0352 and the
  ADR-0227/ADR-0282 amendments (Build Order 4).

Deviations from the PRD: None.

## Absence Evidence

1. The omit key — command: `grep -rn 'omit_paths\|OmitPaths\|omit-path' internal`; result: no match.
2. Claude segmentation — command: `grep -rn 'planReviewSegments(' internal/cli`; result: one non-test caller, inside the branch for `codex` or a Correction Review.
3. Largest-file refusal — command: `grep -rn 'largest files\|reviewLargestFiles' internal`; result: no match.
4. Per-segment token estimate — command: `grep -n 'EstimatedPromptTokens' internal/cli/review_segments.go`; result: no match.
5. Park with the review reason — command: `grep -n 'BlockerReviewBlocked' internal/delivery/engine.go`; result: the constant and one bare park call.

## ADR Alternatives

- ADR-0227: followed; the budget and its two-bytes-per-token estimate stay, now applied to each segment prompt; its rejected up-front tool denial stays unreachable and untouched.
- ADR-0282: amended; its rejected alternatives stay rejected — raising the bound (the bound and budget are unchanged), splitting by file or by bytes (segments stay commit ranges, Claude bounded by tokens) and reading declarations from the candidate (the key is read at the base).
- ADR-0352: written by this Spec.
- ADR-0165: followed; a blocking review after the archive still asks for a corrective Spec; no alternative it rejected is taken.
- ADR-0169: followed; the diff still starts at the merge base, which is why the key is read there.
- ADR-0196: followed; findings still park only after validation, per segment as today.
- ADR-0197: followed; the lineage keeps two rounds and round 2 stays one delta prompt; its rejected review cache stays rejected.
- ADR-0233: unaffected; skill version regeneration and review-only corrections are untouched.
- ADR-0272: followed; the record keeps its per-checkout key and round history; only an omission reason value and a segment field are added.
- ADR-0279: unaffected; dismissals that survive a fix are untouched.
- ADR-0286: followed; the Correction Review stays unsegmented with the byte bound for either provider.
- ADR-0317: unaffected; the Result Status Note dismissal is untouched.
- ADR-0092: unaffected; the typed backlog keeps its contract, and this Spec's source is a Backlog Entry adopted at birth.
- ADR-0120: unaffected; the history root is untouched.
- ADR-0247: unaffected; the archive leaves an Archive Record as before.
- ADR-0248: unaffected; history sanitization is untouched.
- ADR-0251: unaffected; legacy archive reading is untouched.
- ADR-0254: unaffected; Baseline update history sanitization is untouched.
- ADR-0258: unaffected; QA rerun classes and reduced history are untouched.
- ADR-0260: unaffected; archive relinks are untouched.
- ADR-0261: unaffected; the record's pull request and head fetch are untouched.
- ADR-0262: unaffected; cancellation is untouched.
- ADR-0269: unaffected; archive and cancel refusals are untouched.
- ADR-0277: unaffected; Legacy Archive Folder conversion is untouched.
- ADR-0278: followed; deleted history files keep the `retired-history` omission and its precedence.
- ADR-0295: unaffected; reachability is untouched.
- ADR-0299: unaffected; the history conversion is untouched.
- ADR-0306: unaffected; the repository-path refusal is untouched.
- ADR-0308: unaffected; unproven-link rebasing is untouched.
- ADR-0348: unaffected; the absent Spec Root rule is untouched.
- ADR-0250: unaffected; no Baseline module changes, so no module version is recorded.
- ADR-0256: followed; Task 04 updates the skills that describe the changed review and deliver surfaces, so the release does not lag.
- ADR-0281: unaffected; no test rebuilds an executable.
- ADR-0290: unaffected; item binary and worktree reclamation are untouched.
- ADR-0345: followed; the Behavior Surface and skill version records are rewritten only by their declared commands.
- ADR-0349: followed; the authorization record carries a `minor` Release Ceiling and no Task releases.

## Integration Points

The Agent runtime receives up to four Claude prompts in one prepared session;
no adapter contract changes. Git is read through `preflight.GitRunner`.

## Testing Approach

Every test uses the existing `newReviewCommandFixture` disposable repositories
and fake runners (`reviewCommandRunner`, `segmentReviewRunner`); no session,
adapter or network is reached. New top-level tests call `t.Parallel()` first.

- `internal/config/pre_pr_review_omit_paths_test.go`:
  `TestPrePRReviewOmitPathsLoadsAndValidates` (glob, file, directory entry;
  Project replaces User; each refusal of API Contract 1). Neither `roundfix
  init` template carries the key: the existing
  `TestDefaultConfigSelectsNoPrePRReviewProvider` guard keeps its bytes.
- `internal/cli/review_omit_paths_test.go`:
  `TestReviewOmitsConfiguredPathsFromTheBaseCommit` (base glob omits
  `gen/a_snapshot.json` as `omit-path`, keeps `src/a.go`; a path both a
  derived entry and the key match stays `generated`; the key only in the head
  omits nothing; round 2 delta omits it too),
  `TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths` (the summary line
  equals today's literal sentence).
- `internal/cli/review_claude_segment_test.go`:
  `TestReviewSegmentsAClaudeCandidateOverItsBudget` (three commits of 475,000
  bytes: exactly two segments, `<base>..<c2>` and `<c2>..<c3>`, one prepared
  session, Spec context only in the first prompt, each segment estimate within
  the budget and equal to its sent prompt's, record estimate their sum),
  `TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments` (five commits
  of 700,000 bytes),
  `TestReviewKeepsClaudeRoundTwoAndOversizedCommitsUnsegmented`,
  `TestReviewKeepsCodexSegmentRecordsWithoutTokenEstimates`.
- `internal/cli/review_size_refusal_test.go`:
  `TestReviewSizeRefusalNamesTheLargestFiles`, one subtest per refusal of API
  Contract 6, with files of distinct sizes committed out of alphabetical
  order, two of equal size, a path `dados gerados/ação.json`, an omitted
  larger file never listed, and a two-file case; expected bytes come from
  `git diff <base> <head> -- <path>` of each path.
- `internal/delivery/review_blocked_park_test.go`:
  `TestReviewBlockedParkCarriesTheReviewReason` (a two-line reason parks with
  its first line; an empty reason parks bare; both classify `review`).

Sabotages the QA gate runs, each on its own temporary copy, with the test and
the input that makes the mutation observable:

1. Drop the `omit-path` match → `TestReviewOmitsConfiguredPathsFromTheBaseCommit` (the snapshot appears in the prompt).
2. Read the key from the head → same test, head-only case.
3. Always print the `omit_paths` count → `TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths`.
4. Restrict segmentation to `codex` again → `TestReviewSegmentsAClaudeCandidateOverItsBudget`.
5. Bound Claude segments by `reviewDiffBound` → same test (three segments instead of two, because 950,000 bytes exceed the byte bound but fit the token budget).
6. Drop the four-segment limit for Claude → `TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments`.
7. Rank by path or ascending bytes → `TestReviewSizeRefusalNamesTheLargestFiles` (out-of-order sizes, the tie).
8. Parse the path from the `diff --git` header → same test (`dados gerados/ação.json` is quoted with octal escapes there).
9. Park bare `review-blocked` → `TestReviewBlockedParkCarriesTheReviewReason`.

## Build Order

1. Review Omit Paths: config key, validation, matcher (no init template line), scoped
   diff reason and summary clause, with their tests, and the configuration
   guide's key row and Behavior Surface Record.
2. Claude Review Segments: fit rule, Claude admission, per-segment estimate
   (depends on: 1, both edit `internal/cli/review.go`).
3. Actionable size refusal and park reason (depends on: 2, both edit
   `internal/cli/review.go` and `internal/cli/review_segments.go`).
4. Review and deliver guides, Roundfix references, `CONTEXT.md`, ADR
   amendments, skill version and Behavior Surface Record (depends on: 3, it
   documents 1–3).
5. QA gate (depends on: 4).

Shared with the parallel Specs 0343 and 0344, which deliver after this one:
the Roundfix skill's `SKILL.md` pair and `skills/testdata/owned-skill-versions.json`
(version fields), `docs/references/behavior-surfaces.json` and `CONTEXT.md`.
This Spec edits only the `review` and `deliver` references of the Roundfix
skill, never `modes.md`, `follow-up.md` or `implement.md`, and adds glossary
entries next to **Review Segment** only.

## Risks & Considerations

- A later Claude segment can still meet `Prompt is too long` if the runtime
  does not compact; ADR-0227's quoted reason blocks it and the item parks with
  that reason.
- The key is read at the merge base, so an adopter's item branched before the
  key landed needs the default branch merged into it; the guides and the
  adopter note say so.
- Large fixtures (up to 3.5 MB per test repository) mirror today's segment
  tests; each test removes its repository through `t.TempDir`.

## Glossary

- adds: **Review Omit Path** — declared in `_prd.md`
- changes: **Review Segment** — declared in `_prd.md`

## Decisions

- The key reuses the derived-path entry grammar and validator. See ADR-0352.
- Claude segments are bounded by the estimate of each prompt as sent, in one
  session. See ADR-0352.
- Refusals rank files by section bytes and attribute sections by Git's name
  list. See ADR-0352.
- The park carries the reason's first line for every blocked review. See
  ADR-0352.
